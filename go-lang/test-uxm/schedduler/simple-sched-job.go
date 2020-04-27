
// GoLang: Schedduler (a cron like thing ...)

package main

import (
	"fmt"
	"time"
	"github.com/onatm/clockwerk"
)

type DummyJob1 struct{}
func (d DummyJob1) Run() {
	fmt.Println("Every 5 seconds")
}

type DummyJob2 struct{}
func (d DummyJob2) Run() {
	fmt.Println("Every 10 seconds")
}

func main() {

	c := clockwerk.New()

	var job1 DummyJob1
	c.Every(5 * time.Second).Do(job1)

	var job2 DummyJob2
	c.Every(10 * time.Second).Do(job2)

	c.Start()
	for i := 0; i >= 0; i++ { // main loop
	}

}

// #END
