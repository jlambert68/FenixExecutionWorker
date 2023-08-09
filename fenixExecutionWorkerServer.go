package main

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/gRPCServer"
	"FenixExecutionWorker/outgoingPubSubMessages"
	"fmt"
	uuidGenerator "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

// Used for only process cleanup once
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// Cleanup before close down application
		FenixExecutionWorkerObject.logger.WithFields(logrus.Fields{}).Info("Clean up and shut down servers")

		// Stop Backend GrpcServer Server
		FenixExecutionWorkerObject.GrpcServer.StopGrpcServer()

	}
}

func fenixExecutionWorkerMain() {

	// Create Unique Uuid for run time instance used as identification when communication with ExecutionWorker
	common_config.ApplicationRunTimeUuid = uuidGenerator.New().String()
	fmt.Println("sharedCode.ApplicationRunTimeUuid: " + common_config.ApplicationRunTimeUuid)

	// Set up BackendObject
	FenixExecutionWorkerObject = &fenixExecutionWorkerObjectStruct{
		logger:     nil,
		GrpcServer: &gRPCServer.FenixExecutionWorkerGrpcObjectStruct{},
	}

	// Init logger
	FenixExecutionWorkerObject.InitLogger("")

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	// Initiate Logger for gRPC-server
	FenixExecutionWorkerObject.GrpcServer.InitiateLogger(FenixExecutionWorkerObject.logger)

	msg := "Hello World"
	result, returnMessageString, err := outgoingPubSubMessages.Publish(os.Stdout, msg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if result {
		fmt.Printf("Message published successfully: %s\n", returnMessageString)
	} else {
		fmt.Printf("Message publish failed: %s\n", returnMessageString)
	}

	// Start Backend GrpcServer-server
	FenixExecutionWorkerObject.GrpcServer.InitGrpcServer(FenixExecutionWorkerObject.logger)

}
