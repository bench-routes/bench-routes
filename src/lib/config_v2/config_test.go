package configparser

import (
	"errors"
	"testing"
)

const (
	LoadType uint8 = 1
	Validatetype uint8 = 2
	AddAPIType uint8 = 3
)

type Test struct{
	name 		string	
	file 		string
	Type		uint8
	err 		string
	api			API
	shouldErr	bool
}

var tests = []Test{
	{
	name: "testing domain name missing",
	file: "./testdata/config_domain_missing.yml",
	Type: Validatetype,
	shouldErr: true,
	err: "Should have an error of missing domain",
	},
	{
	name: "testing domain name missing",
	file: "./testdata/config_domain_missing.yml",
	Type: Validatetype,
	shouldErr: true,
	err: "Should have an error of missing domain",
	},
}

func TestLoadAndValidate(t *testing.T) {
	for _ ,s := range tests{
		t.Run(s.name,func (t *testing.T)  {
			inst := Config{
				Address: s.file,
			}

			switch s.Type {
			case 1:
				_,err := inst.Load()
				if (err != nil) != (!s.shouldErr) {
					if(err == nil){
						err = errors.New(s.err) 
					}
					t.Error(err);
				} 
			case 2:
				conf,err := inst.Load()
				if (err != nil) != (!s.shouldErr) {
					if(err == nil){
						err = errors.New(s.err) 
					}
					t.Error(err);
					break
				} 
				err = conf.Validate()
				if (err != nil) != (!s.shouldErr) {
					if(err == nil){
						err = errors.New(s.err) 
					}
					t.Error(err);
				} 
			case 3:
				conf,err := inst.Load()
				if (err != nil) != (!s.shouldErr) {
					if(err == nil){
						err = errors.New(s.err) 
					}
					t.Error(err);
					break
				} 
				conf,err = conf.AddAPI(s.api)
				if (err != nil) != (!s.shouldErr) {
					if(err == nil){
						err = errors.New(s.err) 
					}
					t.Error(err);
				}
				apis := conf.Root.APIs
				api := apis[len(apis)-1]
				if api.Name != s.api.Name{
					t.Error(errors.New("ERROR occured in adding API"));
				}
			}
		})
	}
}
