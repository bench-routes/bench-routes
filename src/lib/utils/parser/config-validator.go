package parser

import (
	"log"
)

// ValidateRoutesProp validates the `routes` property
// in the configuration file.
func ValidateRoutesProp(routes []Routes) {
	if len(routes) == 0 {
		log.Fatalf("`routes` property is missing.")
	} else {
		for i, route := range routes {
			if route.Method == "" {
				log.Fatalf("`method` property for route #%d is missing.", i+1)
			}
			if route.Route == "" {
				log.Fatalf("`route` property for route #%d is missing.", i+1)
			}
			if route.URL == "" {
				log.Fatalf("`url` property for route #%d is missing.", i+1)
			}
			if len(route.Header) != 0 {
				for j, header := range route.Header {
					if header.OfType == "" {
						log.Fatalf("`type` property for header #%d of route `%s` is missing.", j+1, route.Route)
					}
					if header.Value == "" {
						log.Fatalf("`value` property for header #%d of route `%s` is missing.", j+1, route.Route)
					}
				}
			}
			if len(route.Params) != 0 {
				for j, param := range route.Params {
					if param.Name == "" {
						log.Fatalf("`name` property for param #%d of route `%s` is missing.", j+1, route.Route)
					}
					if param.Value == "" {
						log.Fatalf("`value` property for param #%d of route `%s` is missing.", j+1, route.Route)
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
		log.Fatalf("`test_interval` property is missing.")
	} else {
		for i, interval := range intervals {
			if interval.Test == "" {
				log.Fatalf("`test` property for interval #%d is missing.", i+1)
			}
			if interval.Type == "" {
				log.Fatalf("`type` property for interval #%d is missing.", i+1)
			}
			if interval.Duration == nil {
				log.Fatalf("`duration` property for interval #%d is missing.", i+1)
			}
		}
	}
}

// ValidateUtilsConf validates the `utils` property
// in the configuration file.
func ValidateUtilsConf(config *UConfig) {
	if config.RespChanges.Mean == nil {
		log.Fatalf("`mean` property under `response-length` is missing.")
	}
	if config.RespChanges.Mode == nil {
		log.Fatalf("`mode` property under `response-length` is missing.")
	}
	if config.ServicesSignal.FloodPing == "" {
		log.Fatalf("`flood-ping` property under `services-state` is missing.")
	}
	if config.ServicesSignal.Jitter == "" {
		log.Fatalf("`jitter` property under `services-state` is missing.")
	}
	if config.ServicesSignal.Ping == "" {
		log.Fatalf("`ping` property under `services-state` is missing.")
	}
	if config.ServicesSignal.ReqResDelayMonitoring == "" {
		log.Fatalf("`req-res-delay-or-monitoring` property under `services-state` is missing.")
	}
}

// ValidatePasswordProp validates the `password` property in
// the configuration file.
func ValidatePasswordProp(password string) {
	if password == "" {
		log.Fatalf("`password` property is missing.")
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
