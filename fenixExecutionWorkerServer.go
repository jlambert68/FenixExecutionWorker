package main

import (
	"FenixExecutionWorker/gRPCServer"
	"github.com/sirupsen/logrus"
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

	// Start Backend GrpcServer-server
	FenixExecutionWorkerObject.GrpcServer.InitGrpcServer(FenixExecutionWorkerObject.logger)

}
