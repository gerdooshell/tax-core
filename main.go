package main

import (
	"flag"
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
	restApi.ServeHTTP()
}
