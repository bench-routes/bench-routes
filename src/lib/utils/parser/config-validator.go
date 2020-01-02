package parser

import (
	"github.com/zairza-cetb/bench-routes/src/lib/utils/logger"
)

// ValidateRoutesProp validates the `routes` property
// in the configuration file.
func ValidateRoutesProp(routes []Routes) {
	if len(routes) == 0 {
		logger.TerminalandFileLogger.Fatalf("`routes` property is missing.")
	} else {
		for i, route := range routes {
			if route.Method == "" {
				logger.TerminalandFileLogger.Fatalf("`method` property for route #%d is missing.", i+1)
			}
			if route.Route == "" {
				logger.TerminalandFileLogger.Fatalf("`route` property for route #%d is missing.", i+1)
			}
			if route.URL == "" {
				logger.TerminalandFileLogger.Fatalf("`url` property for route #%d is missing.", i+1)
			}
			if len(route.Header) != 0 {
				for j, header := range route.Header {
					if header.OfType == "" {
						logger.TerminalandFileLogger.Fatalf("`type` property for header #%d of route `%s` is missing.", j+1, route.Route)
					}
					if header.Value == "" {
						logger.TerminalandFileLogger.Fatalf("`value` property for header #%d of route `%s` is missing.", j+1, route.Route)
					}
				}
			}
			if len(route.Params) != 0 {
				for j, param := range route.Params {
					if param.Name == "" {
						logger.TerminalandFileLogger.Fatalf("`name` property for param #%d of route `%s` is missing.", j+1, route.Route)
					}
					if param.Value == "" {
						logger.TerminalandFileLogger.Fatalf("`value` property for param #%d of route `%s` is missing.", j+1, route.Route)
					}
				}
			}
		}
	}
}

// ValidateIntervalProp validates the `test_interval` property
// in the configuration file.
func ValidateIntervalProp(intervals []Interval) {
	if len(intervals) == 0 {
		logger.TerminalandFileLogger.Fatalf("`test_interval` property is missing.")
	} else {
		for i, interval := range intervals {
			if interval.Test == "" {
				logger.TerminalandFileLogger.Fatalf("`test` property for interval #%d is missing.", i+1)
			}
			if interval.Type == "" {
				logger.TerminalandFileLogger.Fatalf("`type` property for interval #%d is missing.", i+1)
			}
			if interval.Duration == nil {
				logger.TerminalandFileLogger.Fatalf("`duration` property for interval #%d is missing.", i+1)
			}
		}
	}
}

// ValidateUtilsConf validates the `utils` property
// in the configuration file.
func ValidateUtilsConf(config *UConfig) {
	if config.RespChanges.Mean == nil {
		logger.TerminalandFileLogger.Fatalf("`mean` property under `response-length` is missing.")
	}
	if config.RespChanges.Mode == nil {
		logger.TerminalandFileLogger.Fatalf("`mode` property under `response-length` is missing.")
	}
	if config.ServicesSignal.FloodPing == "" {
		logger.TerminalandFileLogger.Fatalf("`flood-ping` property under `services-state` is missing.")
	}
	if config.ServicesSignal.Jitter == "" {
		logger.TerminalandFileLogger.Fatalf("`jitter` property under `services-state` is missing.")
	}
	if config.ServicesSignal.Ping == "" {
		logger.TerminalandFileLogger.Fatalf("`ping` property under `services-state` is missing.")
	}
	if config.ServicesSignal.ReqResDelayMonitoring == "" {
		logger.TerminalandFileLogger.Fatalf("`req-res-delay-or-monitoring` property under `services-state` is missing.")
	}
}

// ValidatePasswordProp validates the `password` property in
// the configuration file.
func ValidatePasswordProp(password string) {
	if password == "" {
		logger.TerminalandFileLogger.Fatalf("`password` property is missing.")
	}
}

// Validate validates the local configuration file.
func (inst YAMLBenchRoutesType) Validate() bool {
	var config = *inst.Config
	ValidatePasswordProp(config.Password)
	ValidateRoutesProp(config.Routes)
	ValidateIntervalProp(config.Interval)
	ValidateUtilsConf(&config.UtilsConf)
	return true
}
