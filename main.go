package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	restApi "github.com/gerdooshell/tax-core/controller/rest_api"
	"github.com/gerdooshell/tax-core/environment"
	logger "github.com/gerdooshell/tax-logger-client-go"
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

	if loggingConfFilePath, err := getLoggingConfigFilePath(); err != nil {
		fmt.Println(err)
	} else {
		if err = logger.SetUpLogger(context.Background(), string(env), loggingConfFilePath); err != nil {
			fmt.Println(err)
		}
	}
	go runGC()
	restApi.ServeHTTP()
}

func getLoggingConfigFilePath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v/logging-config.json", path), nil
}

func runGC() {
	ticker := time.NewTicker(time.Minute)
	for {
		<-ticker.C
		runtime.GC()
	}
}
