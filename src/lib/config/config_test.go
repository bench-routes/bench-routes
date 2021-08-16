package config

import (
	"fmt"
	"testing"
	"time"
)

type test struct {
	name      string
	file      string
	err       string
	api       API
	shouldErr bool
}

var validateTests = []test{
	{
		name:      "`Name` field is MISSING",
		file:      "./testdata/config_name_missing.bad.yml",
		shouldErr: true,
		err:       "validation error : `Name` field of # 1 API cannot be empty",
	},
	{
		name:      "VALID config file",
		file:      "./testdata/config_valid.good.yml",
		shouldErr: false,
	},
	{
		name:      "`Every` field is MISSING",
		file:      "./testdata/config_every_missing.bad.yml",
		shouldErr: true,
		err:       "validation error : `Every` field value of # 0 API is not supported",
	},
	{
		name:      "`domain_or_ip` field is MISSING",
		file:      "./testdata/config_domain_missing.bad.yml",
		shouldErr: true,
		err:       "validation error : `Domain_or_Ip` field of # 0 API cannot be empty",
	},
	{
		name:      "`Route` field is MISSING",
		file:      "./testdata/config_route_missing.bad.yml",
		shouldErr: true,
		err:       "validation error : `Route` field of # 0 API cannot be empty",
	},
	{
		name:      "`Method` field is MISSING",
		file:      "./testdata/config_method_missing.bad.yml",
		shouldErr: true,
		err:       "validation error : `Method` field of # 0 API cannot be empty",
	},
	{
		name:      "`Method` field is INVALID",
		file:      "./testdata/config_method_invalid.bad.yml",
		shouldErr: true,
		err:       "validation error : `Method` field of # 0 API is not supported",
	},
	{
		name:      "`Domain` field is INVALID",
		file:      "./testdata/config_domain_invalid.bad.yml",
		shouldErr: true,
		err:       "validation error : `Domain` field of # 0 API does not match the valid regex",
	},
}

var loadTests = []test{
	{
		name:      "Loading config file",
		file:      "./testdata/normal_load.good.yml",
		shouldErr: false,
	},
	{
		name:      "Loading INVALID config file",
		file:      "./testdata/invalid_load.bad.yml",
		shouldErr: true,
		err:       "unmarshalling file content: yaml: line 3: mapping values are not allowed in this context",
	},
}

var addAPITests = []test{
	{
		name:      "Adding valid API to config file",
		file:      "./testdata/config_API_valid.good.yml",
		shouldErr: false,
		api: API{
			Name:   "API_name_1",
			Method: "get",
			Every:  time.Second * 5,
		},
	},
	{
		name:      "Adding valid API to config file",
		file:      "./testdata/config_API_valid.good.yml",
		shouldErr: false,
		api: API{
			Name:   "API_name_1",
			Method: "GET",
			Every:  time.Minute * 5,
		},
	},
}

func TestLoad(t *testing.T) {
	for _, s := range loadTests {
		t.Run(s.name, func(t *testing.T) {
			c := &Config{
				path: s.file,
			}
			_, err := c.Reload()
			if err != nil {
				if !s.shouldErr {
					t.Fatalf("%s: error was not expected", err)
				}
				if err.Error() != s.err {
					t.Fatalf("%s: error does not match", err)
				}
			} else {
				if s.shouldErr {
					t.Fatalf("%s: error was expected", s.err)
				}
			}

		})
	}
}

func TestValidate(t *testing.T) {
	for _, s := range validateTests {
		t.Run(s.name, func(t *testing.T) {
			c := &Config{
				path: s.file,
			}
			_, err := c.Reload()
			// if err != nil {
			// 	t.Fatal("Error in reloading: %w", err)
			// }
			if err != nil {
				if !s.shouldErr {
					t.Fatalf("%s: error was not expected", err)
				}
				if err.Error() != s.err {
					t.Fatalf("%s: error does not match", err)
				}
			} else {
				if s.shouldErr {
					t.Fatalf("%s: error was expected", s.err)
				}
			}
		})
	}
}

func TestAddAPI(t *testing.T) {
	for _, s := range addAPITests {
		t.Run(s.name, func(t *testing.T) {
			c := &Config{
				path: s.file,
			}
			c, err := c.Reload()
			if err != nil {
				t.Fatal("Error in reloading: %w", err)
			}
			c, err = c.Add(s.api)
			if err != nil {
				if !s.shouldErr {
					t.Fatalf("%s: error was not expected", err)
				}
				if err.Error() != s.err {
					t.Fatalf("%s: error does not match", err)
				}
			} else {
				if s.shouldErr {
					t.Fatalf("%s: error was expected", s.err)
				}
			}
			api := c.APIs[len(c.APIs)-1]
			if api.Name != s.api.Name {
				t.Error(fmt.Errorf("ERROR occured in adding API"))
			}
		})
	}
}
