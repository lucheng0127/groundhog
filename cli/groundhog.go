package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lucheng0127/groundhog/common"
	"github.com/lucheng0127/groundhog/gh"
)

var debug, showVersion bool
var cfgFile, logFile string

// VERSION : Version of groundhog
var VERSION = "0.0.1-dev"

func main() {
	// Initialize comand line parameters
	flag.BoolVar(&showVersion, "version", false, "Show version info")
	flag.BoolVar(&debug, "debug", false, "Provide debug info")
	flag.StringVar(&cfgFile, "config", "", "Config file")
	flag.StringVar(&logFile, "logfile", "/var/log/groundhog.log", "Log file")
	flag.Parse()

	if showVersion {
		fmt.Printf("Groundhog Version: %s\n", VERSION)
		os.Exit(0)
	}
	if cfgFile == "" {
		fmt.Println("Configure file is required")
		fmt.Println("Try groundhog -h for more information")
		os.Exit(1)
	}

	// Setup and get logger
	logger, err := common.SetupLogger(debug, logFile)
	if err != nil {
		fmt.Printf("Failed to setup log: %s", err)
		os.Exit(1)
	}
	checkErr := func(err error) {
		if err != nil {
			logger.Error(fmt.Sprintf("%+v\n", err))
			os.Exit(1)
		}
	}

	// Parse config and run
	cfgInterface, err := common.ParseConfig(cfgFile)
	checkErr(err)
	switch cfg := cfgInterface.(type) {
	case common.ServerConfig:
		server, err := gh.NewServer(&cfg)
		checkErr(err)
		checkErr(server.Launch())
	case common.ClientConfig:
		fmt.Printf("Client configure: %v\n", cfg)
	}
}
