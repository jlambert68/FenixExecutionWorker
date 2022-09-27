package main

import (
	"FenixExecutionWorker/common_config"
	"strconv"

	//"flag"
	"fmt"
	"log"
	"os"
)

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func main() {
	//time.Sleep(15 * time.Second)
	fenixExecutionWorkerMain()
}

func init() {
	//executionLocationForWorker := flag.String("startupType", "0", "The application should be started with one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
	//flag.Parse()

	var err error

	// Get Environment variable to tell how this program was started
	var executionLocationForWorker = mustGetenv("ExecutionLocationForWorker")

	switch executionLocationForWorker {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForWorker = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForWorker = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForWorker = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for FenixGuiServer: " + executionLocationForWorker + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	// Address to Fenix Execution Server
	common_config.FenixExecutionServerAddress = mustGetenv("FenixExecutionServerAddress")

	// Port for Fenix Execution Server
	common_config.FenixExecutionWorkerServerPort, err = strconv.Atoi(mustGetenv("FenixExecutionServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixGuiBuilderServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

}
