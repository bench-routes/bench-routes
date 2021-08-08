package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/bench-routes/bench-routes/src/lib/api"
	brConfig "github.com/bench-routes/bench-routes/src/lib/config"
	"github.com/bench-routes/bench-routes/src/lib/log"
	"github.com/bench-routes/bench-routes/src/lib/modules/module"
	"github.com/bench-routes/bench-routes/tsdb/file"
	"github.com/rs/cors"
)

type config struct {
	path      string // Path of bench-routes configuration file.
	logLevel  string
	logFormat string
	port      string
}

func main() {
	memChainsRef := file.Chains
	conf := new(config)
	parsesFlags(conf, os.Args[1:])
	err := validateFlags(conf)
	if err != nil {
		printErrAndExit("validate-flags", err)
		return
	}

	if err = log.Init(log.Config{Format: conf.logFormat, Level: conf.logLevel}); err != nil {
		printErrAndExit("initializing logger", err)
		return
	}

	brConf, err := brConfig.New(conf.path)
	if err != nil {
		printErrAndExit("initializing bench-routes configuration", err)
		return
	}

	log.Info("msg", "configuration", "out", fmt.Sprintf("%v", brConf.APIs))

	machineErrCh := make(chan error)
	machine, err := module.New(module.MachineType, memChainsRef, machineErrCh)
	if err != nil {
		printErrAndExit("creating machine jobs", err)
		return
	}
	log.Info("msg", "launching machine routine")
	go machine.Run()
	if err = machine.Reload(brConf); err != nil {
		printErrAndExit("reloading machine with configuration", err)
		return
	}

	monitorErrCh := make(chan error)
	monitor, err := module.New(module.MonitorType, memChainsRef, monitorErrCh)
	if err != nil {
		printErrAndExit("creating monitor jobs", err)
		return
	}

	log.Info("msg", "launching monitor routine")
	go monitor.Run()
	if err = monitor.Reload(brConf); err != nil {
		printErrAndExit("reloading monitor with configuration", err)
		return
	}

	go func() {
		// Listen to error channels and exit if any error occurs.
		// Todo: This can be optimized to accept crashable and non-crashable errors.
		// Non-crashable errors can be errors from HTTP or Ping request.
		// Rest all errors (like format error, validation error, etc) are crashable
		// and hence, the application should stop in those cases.
		// Right now, we are treating all errors as crashable.
		for i := 1; i <= 2; i++ {
			select {
			case err := <-machineErrCh:
				printErrAndExit("running machine", err)
			case err := <-monitorErrCh:
				printErrAndExit("running monitor", err)
			}
		}
	}()

	reloadSig := make(chan struct{})
	apiModule := api.New(reloadSig, brConf)
	reloader := func(sig <-chan struct{}) {
		// Reloader should not stop the already running application in case of
		// any error in the new changes. If error is found, log and continue using
		// the earlier perfect configuration.
		for range sig {
			tmpBrConf, err := brConfig.New(conf.path)
			if err != nil {
				log.Error("msg", fmt.Sprintf("error reloading configuration: %s", err.Error()))
				continue
			}
			if err = machine.Reload(tmpBrConf); err != nil {
				log.Error("msg", fmt.Sprintf("error reloading machine jobs: %s", err.Error()))
				continue
			}
			brConf = tmpBrConf // We need to update the
			if err = monitor.Reload(tmpBrConf); err != nil {
				log.Error("msg", fmt.Sprintf("error reloading monitor jobs: %s", err.Error()))
				continue
			}
			apiModule.UpdateConf(brConf)
		}
	}
	go reloader(reloadSig)

	log.Info("msg", "Listening at "+conf.port)
	log.Error(http.ListenAndServe(conf.port, cors.Default().Handler(apiModule.Router())).Error())
}

func parsesFlags(cfg *config, args []string) {
	flag.StringVar(&cfg.path, "config.path", "config.yml", "Address of configuration file.")
	flag.StringVar(&cfg.logLevel, "logger.level", "debug", "Level of logging. Valid options include ['debug', 'info', 'warn', 'error'].")
	flag.StringVar(&cfg.logFormat, "logger.format", "logfmt", "Format of logger. Valid options include ['logfmt', json].")
	flag.StringVar(&cfg.port, "app.port", ":9990", "Port at which bench-routes service will listen.")

	_ = flag.CommandLine.Parse(args)
}

func validateFlags(cfg *config) error {
	var err error
	if err = verifyLoggerLevel(cfg.logLevel); err != nil {
		return fmt.Errorf("verify logger level: %w", err)
	}
	if err = verifyLoggerFormat(cfg.logFormat); err != nil {
		return fmt.Errorf("verify logger format: %w", err)
	}
	return nil
}

func verifyLoggerFormat(format string) error {
	switch format {
	case "logfmt", "json":
		return nil
	default:
		return fmt.Errorf("invalid log-format. Valid options include ['logfmt', 'json']")
	}
}

func verifyLoggerLevel(level string) error {
	switch level {
	case "debug", "info", "warn", "error":
		return nil
	default:
		return fmt.Errorf("invalid log-level. Valid options include ['debug', 'info', 'warn', 'error']")
	}
}

// printErrAndExit prints the wrapped the err with the message and halts the execution of program.
func printErrAndExit(message string, err error) {
	log.Error("msg", message, "error", fmt.Errorf("%s: %w", message, err).Error())
	os.Exit(1)
}
