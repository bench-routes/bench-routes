package parser

import (
	"github.com/prometheus/common/log"
	"strconv"
)

// ValidateIntervalProp validates the `test_interval` property
// in the configuration file.
func ValidateIntervalProp(intervals []Interval) {
	if len(intervals) == 0 {
		log.Errorln("`test_interval` property is missing.")
	} else {
		for i, interval := range intervals {
			if interval.Test == "" {
				log.Errorln("`test` property for interval #" + strconv.Itoa(i+1) + " is missing.")
			}
			if interval.Type == "" {
				log.Errorln("`type` property for interval #" + strconv.Itoa(i+1) + " is missing.")
			}
			if interval.Duration == nil {
				log.Errorln("`duration` property for interval #" + strconv.Itoa(i+1) + " is missing.")
			}
		}
	}
}

// ValidateUtilsConf validates the `utils` property
// in the configuration file.
func ValidateUtilsConf(config *UConfig) {
	if config.ServicesSignal.FloodPing == "" {
		log.Errorln("`flood-ping` property under `Services-state` is missing.")
	}
	if config.ServicesSignal.Jitter == "" {
		log.Errorln("`jitter` property under `Services-state` is missing.")
	}
	if config.ServicesSignal.Ping == "" {
		log.Errorln("`ping` property under `Services-state` is missing.")
	}
	if config.ServicesSignal.ReqResDelayMonitoring == "" {
		log.Errorln("`req-res-delay-or-monitoring` property under `Services-state` is missing.")
	}
}

// ValidatePasswordProp validates the `password` property in
// the configuration file.
func ValidatePasswordProp(password string) {
	if password == "" {
		log.Errorln("`password` property is missing.")
	}
}

// Validate validates the local configuration file.
func (inst *Config) Validate() bool {
	config := *inst.Config
	ValidatePasswordProp(config.Password)
	ValidateIntervalProp(config.Interval)
	ValidateUtilsConf(&config.UtilsConf)
	return true
}
