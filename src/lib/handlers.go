package lib

import (
	"log"
)

const (
	// ConfigurationFilePath is the constant path to the configuration file needed to start the application
	// written from root file since the application starts from `make run`
	ConfigurationFilePath = "storage/local-config.yml"
)

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	// load configuration file

	// keep the below line to the end of file so that we ensure that we give a confirmation message only when all the
	// required resources for the application is up and healthy
	log.Println("Bench-routes is up and running")
}

// HandlerPingGeneral handles the ping route
func HandlerPingGeneral() {
}
