package main

import (
	"FenixExecutionWorker/common_config"
	"github.com/sirupsen/logrus"
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

func init() {
	//executionLocationForWorker := flag.String("startupType", "0", "The application should be started with one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
	//flag.Parse()

	var err error

	// Get Environment variable to tell how/were this worker is  running
	var executionLocationForWorker = mustGetenv("ExecutionLocationForWorker")

	switch executionLocationForWorker {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForWorker = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForWorker = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForWorker = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for Worker: " + executionLocationForWorker + ". Expected one of the following: 'LOCALHOST_NODOCKER', 'LOCALHOST_DOCKER', 'GCP'")
		os.Exit(0)

	}

	// Get Environment variable to tell were Fenix Execution Server is running
	var executionLocationForExecutionServer = mustGetenv("ExecutionLocationForFenixTestExecutionServer")

	switch executionLocationForExecutionServer {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForFenixExecutionServer = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForFenixExecutionServer = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForFenixExecutionServer = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for Fenix Execution Server: " + executionLocationForWorker + ". Expected one of the following: 'LOCALHOST_NODOCKER', 'LOCALHOST_DOCKER', 'GCP'")
		os.Exit(0)

	}

	// Address to Fenix Execution Server
	common_config.FenixExecutionServerAddress = mustGetenv("FenixExecutionServerAddress")

	// Port for Fenix Execution Server
	common_config.FenixExecutionServerPort, err = strconv.Atoi(mustGetenv("FenixExecutionServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixExecutionServerPort' to an integer, error: ", err)
		os.Exit(0)
	}

	// Address when Execution Server is not in GCP
	common_config.FenixExecutionServerAddressToDial = common_config.FenixExecutionServerAddress + ":" + strconv.Itoa(common_config.FenixExecutionServerPort)

	// Port for Execution Worker
	common_config.ExecutionWorkerServerPort, err = strconv.Atoi(mustGetenv("ExecutionWorkerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'ExecutionWorkerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Extract Debug level
	var loggingLevel = mustGetenv("LoggingLevel")

	switch loggingLevel {

	case "DebugLevel":
		common_config.LoggingLevel = logrus.DebugLevel

	case "InfoLevel":
		common_config.LoggingLevel = logrus.InfoLevel

	default:
		fmt.Println("Unknown LoggingLevel '" + loggingLevel + "'. Expected one of the following: 'DebugLevel', 'InfoLevel'")
		os.Exit(0)

	}

	// Extract OAuth 2.0 Client ID
	common_config.AuthClientId = mustGetenv("AuthClientId")

	// Extract OAuth 2.0 Client Secret
	common_config.AuthClientSecret = mustGetenv("AuthClientSecret")

	// Extract the GCP-project
	common_config.GcpProject = mustGetenv("GcpProject")

	// Should PubSub be used for sending 'TestInstructionExecutions' to Connector
	common_config.UsePubSubWhenSendingTestInstructionExecutions, err = strconv.ParseBool(mustGetenv("UsePubSubWhenSendingTestInstructionExecutions"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'UsePubSubWhenSendingTestInstructionExecutions' to a boolean, error: ", err)
		os.Exit(0)
	}

	// Extract local path to Service-Account file
	common_config.LocalServiceAccountPath = mustGetenv("LocalServiceAccountPath")
	// The only way have an OK space is to replace an existing character
	if common_config.LocalServiceAccountPath == "#" {
		common_config.LocalServiceAccountPath = ""
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", common_config.LocalServiceAccountPath)

	// Extract environment variable for 'ThisDomainsUuid'
	common_config.ThisDomainsUuid = mustGetenv("ThisDomainsUuid")

	// Extract environment variable for 'TestInstructionExecutionPubSubTopicBase'
	common_config.TestInstructionExecutionPubSubTopicBase = mustGetenv("TestInstructionExecutionPubSubTopicBase")

	// Extract environment variable for 'TestInstructionExecutionPubSubTopicSchema'
	common_config.TestInstructionExecutionPubSubTopicSchema = mustGetenv("TestInstructionExecutionPubSubTopicSchema")

}
