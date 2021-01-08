package main

import (
	"errors"
	"flag"
	"os"

	"github.com/TheSlipper/ParkindStreamer/logging"
	"github.com/TheSlipper/ParkindStreamer/streaming"
)

// Runtime arguments of the parkind client
type runtimeArgs struct {
	verbosity bool
	config    string
	login     string
	password  string
}

// Load the command line and environment arguments into the args singleton
func setUpRuntimeArgs() (args runtimeArgs, err error) {
	flag.BoolVar(&args.verbosity, "verbose", false, "defines how much information should be printed out")
	flag.StringVar(&args.config, "config", "config.json", "path to the configuration file")
	flag.Parse()

	args.login = os.Getenv("LOGIN")
	args.password = os.Getenv("PASSWORD")
	if args.login == "" || args.password == "" {
		return args, errors.New("LOGIN or PASSWORD were not provided")
	}

	return args, nil
}

func main() {
	// Set up runtime arguments or stop execution if not satisfied
	args, err := setUpRuntimeArgs()
	if err != nil {
		logging.ErrorLog(err.Error())
		os.Exit(1)
	}

	// Set up a camera session and a streaming goroutine
	cs, err := streaming.NewCameraSession(1)
	if err != nil {
		logging.ErrorLog(err.Error())
		os.Exit(3)
	}
	defer cs.Close()
	go cs.Stream()

	// Set up local http server instance
	server, err := streaming.CreateHttpServer(args.verbosity)
	if err != nil {
		logging.ErrorLog(err.Error())
		os.Exit(2)
	}

	// Start the server
	err = server.ListenAndServe()
	if err != nil {
		logging.ErrorLog(err.Error())
		os.Exit(3)
	}
}
