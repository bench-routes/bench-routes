package parser

import (
	"strconv"

	"github.com/zairza-cetb/bench-routes/src/lib/logger"
)

// ValidateRoutesProp validates the `routes` property
// in the configuration file.
func ValidateRoutesProp(routes []Routes) {
	if len(routes) == 0 {
		logger.Terminal("`routes` property is missing.", "f")
	} else {
		for i, route := range routes {
			if route.Method == "" {
				msg := "`method` property for route #" + strconv.Itoa(i+1) + " is missing."
				logger.Terminal(msg, "f")
			}
			if route.Route == "" {
				msg := "`route` property for route #" + strconv.Itoa(i+1) + " is missing."
				logger.Terminal(msg, "f")
			}
			if route.URL == "" {
				msg := "`url` property for route #" + strconv.Itoa(i+1) + " is missing."
				logger.Terminal(msg, "f")
			}
			if len(route.Header) != 0 {
				for j, header := range route.Header {
					if header.OfType == "" {
						msg := "`type` property for header #" + strconv.Itoa(j+1) + " of route " + route.Route + "is missing"
						logger.Terminal(msg, "f")
					}
					if header.Value == "" {
						msg := "`value` property for header #" + strconv.Itoa(j+1) + " of route " + route.Route + "is missing"
						logger.Terminal(msg, "f")
					}
				}
			}
			if len(route.Params) != 0 {
				for j, param := range route.Params {
					if param.Name == "" {
						msg := "`name` property for param #" + strconv.Itoa(j+1) + " of route " + route.Route + "is missing"
						logger.Terminal(msg, "f")
					}
					if param.Value == "" {
						msg := "`value` property for param #" + strconv.Itoa(j+1) + " of route " + route.Route + "is missing"
						logger.Terminal(msg, "f")
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
	if config.RespChanges.Mean == nil {
		logger.Terminal("`mean` property under `response-length` is missing.", "f")
	}
	if config.RespChanges.Mode == nil {
		logger.Terminal("`mode` property under `response-length` is missing.", "f")
	}
	if config.ServicesSignal.FloodPing == "" {
		logger.Terminal("`flood-ping` property under `services-state` is missing.", "f")
	}
	if config.ServicesSignal.Jitter == "" {
		logger.Terminal("`jitter` property under `services-state` is missing.", "f")
	}
	if config.ServicesSignal.Ping == "" {
		logger.Terminal("`ping` property under `services-state` is missing.", "f")
	}
	if config.ServicesSignal.ReqResDelayMonitoring == "" {
		logger.Terminal("`req-res-delay-or-monitoring` property under `services-state` is missing.", "f")
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
func (inst YAMLBenchRoutesType) Validate() bool {
	var config = *inst.Config
	ValidatePasswordProp(config.Password)
	ValidateRoutesProp(config.Routes)
	ValidateIntervalProp(config.Interval)
	ValidateUtilsConf(&config.UtilsConf)
	return true
}
