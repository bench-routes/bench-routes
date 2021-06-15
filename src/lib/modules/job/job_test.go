package job

import (
	"log"
	"testing"
	"time"

	config "github.com/bench-routes/bench-routes/src/lib/config_v2"
)

func TestJob(t *testing.T){
	ch := make(chan struct{})
	api := config.API{
		Name: "API_1",
		Every: time.Second*5,
		Domain: "https://www.google.com",
		// Route: "/watch?v=aaDaIGHTGT8",
		Method: "GET",
		Headers: map[string]string{
			"Content-Type" : "application/json",
		},
	}
	exec,err :=NewJob("monitor",ch,&api);
	if err != nil {
		t.Fatalf("Error: %s",err)
	}

	go exec.Execute()
	ch<-struct{}{}
	time.Sleep(time.Second*5)

	ch<-struct{}{}
	time.Sleep(time.Second*5)

	ch<-struct{}{}
	time.Sleep(time.Second*5)

	ch<-struct{}{}
	time.Sleep(time.Second*5)

	log.Println(exec.Info())
	exec.Abort()

}