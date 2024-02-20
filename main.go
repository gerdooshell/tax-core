package main

import (
	"flag"
	"runtime"
	"time"

	restApi "github.com/gerdooshell/tax-core/controller/rest_api"
	"github.com/gerdooshell/tax-core/environment"
)

func readEnvironment() environment.Environment {
	isProdEnvPtr := flag.Bool("prod", false, "is environment prod")
	flag.Parse()
	env := environment.Dev
	if *isProdEnvPtr {
		env = environment.Prod
	}
	return env
}

func main() {
	env := readEnvironment()
	if err := environment.SetEnvironment(env); err != nil {
		panic(err)
	}
	go runGC()
	restApi.ServeHTTP()
}

func runGC() {
	ticker := time.NewTicker(time.Minute)
	for {
		<-ticker.C
		runtime.GC()
	}
}
