package configparser

import (
	"errors"
	"testing"
)

const (
	LoadType     uint8 = 1
	Validatetype uint8 = 2
	AddAPIType   uint8 = 3
)

type Test struct {
	name      string
	file      string
	Type      uint8
	err       string
	api       API
	shouldErr bool
}

var tests = []Test{
	{
		name:      "Loading config file",
		file:      "./testdata/normal_load.yml",
		Type:      LoadType,
		shouldErr: false,
	},
	{
		name:      "Loading INVALID config file",
		file:      "./testdata/invalid_load.yml",
		Type:      LoadType,
		shouldErr: true,
		err:       "Should have error loading due to invalid file structure",
	},
	{
		name:      "`Name` field is MISSING",
		file:      "./testdata/config_name_missing.yml",
		Type:      Validatetype,
		shouldErr: true,
		err:       "Should have an error of missing `Name` field",
	},
	{
		name:      "`Every` field is MISSING",
		file:      "./testdata/config_every_missing.yml",
		Type:      Validatetype,
		shouldErr: true,
		err:       "Should have an error of missing `Every` field",
	},
	{
		name:      "`domain_or_ip` field is MISSING",
		file:      "./testdata/config_domain_missing.yml",
		Type:      Validatetype,
		shouldErr: true,
		err:       "Should have an error of missing `domain_or_ip` field",
	},
	{
		name:      "`Route` field is MISSING",
		file:      "./testdata/config_route_missing.yml",
		Type:      Validatetype,
		shouldErr: true,
		err:       "Should have an error of missing `Route` field",
	},
	{
		name:      "`Method` field is MISSING",
		file:      "./testdata/config_method_missing.yml",
		Type:      Validatetype,
		shouldErr: true,
		err:       "Should have an error of missing `Method` field",
	},
	{
		name:      "`Method` field is INVALID",
		file:      "./testdata/config_method_invalid.yml",
		Type:      Validatetype,
		shouldErr: true,
		err:       "Should have an error of INVALID `Method` field",
	},
	{
		name:      "`Domain` field is INVALID",
		file:      "./testdata/config_domain_invalid.yml",
		Type:      Validatetype,
		shouldErr: true,
		err:       "Should have an error of INVALID `Domain` field",
	},
	{
		name:      "Adding valid API to config file",
		file:      "./testdata/config_API_valid.yml",
		Type:      AddAPIType,
		shouldErr: false,
		api: API{
			Name: "API_name_1",
		},
	},
}

func TestLoadAndValidate(t *testing.T) {
	for _, s := range tests {
		t.Run(s.name, func(t *testing.T) {
			inst := Config{
				Address: s.file,
			}

			switch s.Type {
			case 1:
				_, err := inst.Reload()
				if (err != nil) != (s.shouldErr) {
					if err == nil {
						err = errors.New(s.err)
					}
					t.Error(err)
				}
			case 2:
				conf, err := inst.Reload()

				if err != nil {
					t.Error(err)
					break
				}
				err = conf.Validate()

				if (err != nil) != (s.shouldErr) {
					if err == nil {
						err = errors.New(s.err)
					}
					t.Error(err)
				}
			case 3:
				conf, err := inst.Reload()
				if (err != nil) != (s.shouldErr) {
					if err == nil {
						err = errors.New(s.err)
					}
					t.Error(err)
					break
				}
				conf, err = conf.AddAPI(s.api)
				if (err != nil) != (s.shouldErr) {
					if err == nil {
						err = errors.New(s.err)
					}
					t.Error(err)
				}
				apis := conf.Root.APIs
				api := apis[len(apis)-1]
				if api.Name != s.api.Name {
					t.Error(errors.New("ERROR occured in adding API"))
				}
			}
		})
	}
}

func BenchmarkLoadAndValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range tests {
			b.Run(s.name, func(b *testing.B) {
				inst := Config{
					Address: s.file,
				}

				switch s.Type {
				case 1:
					_, err := inst.Reload()
					if (err != nil) != (s.shouldErr) {
						if err == nil {
							err = errors.New(s.err)
						}
						b.Error(err)
					}
				case 2:
					conf, err := inst.Reload()

					if err != nil {
						b.Error(err)
						break
					}
					err = conf.Validate()

					if (err != nil) != (s.shouldErr) {
						if err == nil {
							err = errors.New(s.err)
						}
						b.Error(err)
					}
				case 3:
					conf, err := inst.Reload()
					if (err != nil) != (s.shouldErr) {
						if err == nil {
							err = errors.New(s.err)
						}
						b.Error(err)
						break
					}
					conf, err = conf.AddAPI(s.api)
					if (err != nil) != (s.shouldErr) {
						if err == nil {
							err = errors.New(s.err)
						}
						b.Error(err)
					}
					apis := conf.Root.APIs
					api := apis[len(apis)-1]
					if api.Name != s.api.Name {
						b.Error(errors.New("ERROR occured in adding API"))
					}
				}
			})
		}
	}

}
