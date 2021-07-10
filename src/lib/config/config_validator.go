package config

import (
	"fmt"
	"regexp"
	"strings"
)

// https://regexr.com/3au3g
const validDomainRegex = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`

func validateAPI(index int, api API) error {
	if api.Name == "" {
		return fmt.Errorf("`Name` field of # %d API cannot be empty", index+1)
	}
	if api.Protocol == "" {
		return fmt.Errorf("`Protocol` field of %d API cannot be empty", index+1)
	}
	if api.Every.String() == "0s" {
		return fmt.Errorf("`Every` field value of # %d API is not supported", index)
	}
	if api.Domain == "" {
		return fmt.Errorf("`Domain_or_Ip` field of # %d API cannot be empty", index)
	}
	if api.Route == "" {
		return fmt.Errorf("`Route` field of # %d API cannot be empty", index)
	}
	method := strings.ToLower(api.Method)
	if method == "" {
		return fmt.Errorf("`Method` field of # %d API cannot be empty", index)
	}
	if method != "get" && method != "post" && method != "put" && method != "delete" && method != "patch" {
		return fmt.Errorf("`Method` field of # %d API is not supported", index)
	}
	regExp := regexp.MustCompile(validDomainRegex)
	if !regExp.MatchString(api.Domain) {
		return fmt.Errorf("`Domain` field of # %d API does not match the valid regex", index)
	}
	return nil
}

// Validate validates APIs data which is parsed from the config file.
func (c *Config) Validate() error {
	apis := c.APIs
	for i := range apis {
		if err := validateAPI(i, apis[i]); err != nil {
			return fmt.Errorf("validation error : %w", err)
		}
	}
	return nil
}
