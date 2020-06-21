package parser

import (
	"strconv"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
)

// ValidateRoutesProp validates the `routes` property
// in the configuration file.
func ValidateRoutesProp(routes []Route) {
	if len(routes) == 0 {
		logger.Terminal("`routes` property is missing.", "f")
	} else {
		for i, route := range routes {
			if route.URL == "" {
				msg := "`url` property for route #" + strconv.Itoa(i+1) + " is missing."
				logger.Terminal(msg, "f")
			}
		}
	}
}

// ValidateIntervalProp validates the `test_interval` property
// in the configuration file.
func ValidateIntervalProp(intervals []Interval) {
	if len(intervals) == 0 {
		logger.Terminal("`test_interval` property is missing.", "f")
	} else {
		for i, interval := range intervals {
			if interval.Test == "" {
				msg := "`test` property for interval #" + strconv.Itoa(i+1) + " is missing."
				logger.Terminal(msg, "f")
			}
			if interval.Type == "" {
				msg := "`type` property for interval #" + strconv.Itoa(i+1) + " is missing."
				logger.Terminal(msg, "f")
			}
			if interval.Duration == nil {
				msg := "`duration` property for interval #" + strconv.Itoa(i+1) + " is missing."
				logger.Terminal(msg, "f")
			}
		}
	}
}

// ValidateUtilsConf validates the `utils` property
// in the configuration file.
func ValidateUtilsConf(config *UConfig) {
	if config.ServicesSignal.FloodPing == "" {
		logger.Terminal("`flood-ping` property under `Services-state` is missing.", "f")
	}
	if config.ServicesSignal.Jitter == "" {
		logger.Terminal("`jitter` property under `Services-state` is missing.", "f")
	}
	if config.ServicesSignal.Ping == "" {
		logger.Terminal("`ping` property under `Services-state` is missing.", "f")
	}
	if config.ServicesSignal.ReqResDelayMonitoring == "" {
		logger.Terminal("`req-res-delay-or-monitoring` property under `Services-state` is missing.", "f")
	}
}

// ValidatePasswordProp validates the `password` property in
// the configuration file.
func ValidatePasswordProp(password string) {
	if password == "" {
		logger.Terminal("`password` property is missing.", "f")
	}
}

// Validate validates the local configuration file.
func (inst *Config) Validate() bool {
	config := *inst.Config
	ValidatePasswordProp(config.Password)
	ValidateRoutesProp(config.Routes)
	ValidateIntervalProp(config.Interval)
	ValidateUtilsConf(&config.UtilsConf)
	return true
}
