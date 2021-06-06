package configparser

import (
	"errors"
	"strconv"
	"strings"
)

func validateAPI(index int,api API) error{
	if api.Name==""{
		return errors.New("`Name` property of API #"+strconv.Itoa(index+1)+" is missing")
	}
	if api.Every.String()=="0s"{
		return errors.New("`Every` property of API #"+strconv.Itoa(index+1)+" is invalid")
	}
	if api.Domain==""{
		return errors.New("`Domain_or_Ip` property of API #"+strconv.Itoa(index+1)+" is missing")
	}
	if api.Route==""{
		return errors.New("`Route` property of API #"+strconv.Itoa(index+1)+" is missing")
	}
	method := strings.ToLower(api.Method)
	if (method!="get"&&method!="post"&&method!="put"&&method!="delete"&&method!="patch"){
		return errors.New("`Method` property of API #"+strconv.Itoa(index+1)+" is invalid")
	}
	return nil
}

func(c *Config) Validate() error{
	apis := c.Root.APIs
	
	for i,a := range apis{
		err := validateAPI(i,a);
		if err != nil {
			return err
		}
	}
	return nil
}