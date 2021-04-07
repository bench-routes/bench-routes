package parser

import (
	"errors"
	"fmt"
	"strconv"

	valid "github.com/asaskevich/govalidator"
	"github.com/prometheus/common/log"
)

// ValidateIntervalProp validates the `test_interval` property
// in the configuration file.
func ValidateIntervalProp(intervals []Interval) []error {
	services := map[string]bool{
		"ping":       true,
		"jitter":     true,
		"monitoring": true,
	}
	durationType := map[string]bool{
		"sec": true,
		"min": true,
		"hr":  true,
	}
	intervalErr := []error{}
	if len(intervals) == 0 {
		intervalErr = append(intervalErr, errors.New("`test_interval` property is missing"))
	} else {
		if !services["ping"] {
			intervalErr = append(intervalErr, fmt.Errorf("ping service missing in the rest_interval"))
		}
		if !services["jitter"] {
			intervalErr = append(intervalErr, fmt.Errorf("jitter service missing in the rest_interval"))
		}
		if !services["monitoring"] {
			intervalErr = append(intervalErr, fmt.Errorf("monitoring service missing in the rest_interval"))
		}
		for i, interval := range intervals {
			if interval.Test == "" {
				intervalErr = append(intervalErr, errors.New("`test` property for interval #"+strconv.Itoa(i+1)+" is missing"))
			}
			if !durationType[interval.Type] {
				intervalErr = append(intervalErr, errors.New("`type` property for interval #"+strconv.Itoa(i+1)+" is not valid, It must be sec, min or hr for second, minute and hour respectively"))
			}
			if interval.Duration == nil {
				intervalErr = append(intervalErr, errors.New("`duration` property for interval #"+strconv.Itoa(i+1)+" is missing"))
			}
		}
	}
	return intervalErr
}

// ValidateUtilsConf validates the `utils` property
// in the configuration file.
func ValidateUtilsConf(config *UConfig) []error {
	utilsErr := []error{}
	if config.ServicesSignal.FloodPing == "" {
		utilsErr = append(utilsErr, errors.New("`flood-ping` property under `Services-state` is missing"))
	}
	if config.ServicesSignal.Jitter == "" {
		utilsErr = append(utilsErr, errors.New("`jitter` property under `Services-state` is missing"))
	}
	if config.ServicesSignal.Ping == "" {
		utilsErr = append(utilsErr, errors.New("`ping` property under `Services-state` is missing"))
	}
	if config.ServicesSignal.ReqResDelayMonitoring == "" {
		utilsErr = append(utilsErr, errors.New("`req-res-delay-or-monitoring` property under `Services-state` is missing"))
	}
	return utilsErr
}

// ValidatePasswordProp validates the `password` property in
// the configuration file.
func ValidatePasswordProp(password string) error {
	if password == "" {
		return errors.New("password should not be empty")
	}
	return nil
}

func ValidateRoutes(routes *[]Route) []error {
	Method := map[string]bool{
		"POST":    true,
		"GET":     true,
		"DELETE":  true,
		"PUT":     true,
		"PATCH":   true,
		"COPY":    true,
		"HEAD":    true,
		"OPTIONS": true,
	}
	routesErr := []error{}
	for _, route := range *routes {
		if !valid.IsURL(route.URL) {
			routesErr = append(routesErr, fmt.Errorf("\"%s\" is not a valid URL, Update the URL", route.URL))
		}
		if !Method[route.Method] {
			routesErr = append(routesErr, fmt.Errorf("for \"%s\" \"%s\" is not a valid http method, Update the Method", route.URL, route.Method))
		}
	}
	return routesErr
}

// Validate validates the local configuration file.
func (inst *Config) Validate() bool {
	er := []error{}
	config := *inst.Config
	err := ValidatePasswordProp(config.Password)
	if err != nil {
		er = append(er, err)
	}
	intervalErr := ValidateIntervalProp(config.Interval)
	if len(intervalErr) != 0 {
		er = append(er, intervalErr...)
	}
	utilsErr := ValidateUtilsConf(&config.UtilsConf)
	if len(utilsErr) != 0 {
		er = append(er, utilsErr...)
	}
	routesErr := ValidateRoutes(&config.Routes)
	if len(routesErr) != 0 {
		er = append(er, routesErr...)
	}
	if len(er) != 0 {
		for i, err := range er {
			log.Errorln(i+1, err)
		}
		panic(er)
	}
	return true
}
