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
	go runGC()
	go setupLogging(env)
	restApi.ServeHTTP()
}

func setupLogging(env environment.Environment) {
	var err error
	loggingConfFilePath, err := getLoggingConfigFilePath()
	if err != nil {
		fmt.Println(err)
		return
	}
	<-time.After(time.Second * 5)
	if err = logger.SetUpLogger(context.Background(), string(env), loggingConfFilePath); err != nil {
		fmt.Println(err)
	}
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
