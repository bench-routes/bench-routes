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
		err:       "Should have an error of missing `Name` field",
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
		err:       "Should have an error of missing `Every` field",
	},
	{
		name:      "`domain_or_ip` field is MISSING",
		file:      "./testdata/config_domain_missing.bad.yml",
		shouldErr: true,
		err:       "Should have an error of missing `domain_or_ip` field",
	},
	{
		name:      "`Route` field is MISSING",
		file:      "./testdata/config_route_missing.bad.yml",
		shouldErr: true,
		err:       "Should have an error of missing `Route` field",
	},
	{
		name:      "`Method` field is MISSING",
		file:      "./testdata/config_method_missing.bad.yml",
		shouldErr: true,
		err:       "Should have an error of missing `Method` field",
	},
	{
		name:      "`Method` field is INVALID",
		file:      "./testdata/config_method_invalid.bad.yml",
		shouldErr: true,
		err:       "Should have an error of INVALID `Method` field",
	},
	{
		name:      "`Domain` field is INVALID",
		file:      "./testdata/config_domain_invalid.bad.yml",
		shouldErr: true,
		err:       "Should have an error of INVALID `Domain` field",
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
		err:       "Should have error loading due to invalid file structure",
	},
}

var addAPITests = []test{
	{
		name:      "Adding valid API to config file",
		file:      "./testdata/config_API_valid.good.yml",
		shouldErr: false,
		api: API{
			Name:   "API_name_1",
			Method: "Get",
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
			c, err := c.Reload()

			if err != nil {
				t.Fatal("Error in reloading: %w", err)
			}

			if err = c.Validate(); err != nil {
				if !s.shouldErr {
					t.Fatalf("%s: error was not expected", err)
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
			c, err = c.AddAPI(s.api)
			if err != nil {
				if !s.shouldErr {
					t.Fatalf("%s: error was not expected", err)
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
