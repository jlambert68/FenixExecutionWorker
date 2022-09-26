package main

import (
	"FenixExecutionWorker/testInstructionExecutionEngine"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// Used for only process cleanup once
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// Cleanup before close down application
		fenixExecutionWorkerObject.logger.WithFields(logrus.Fields{}).Info("Clean up and shut down servers")

		// Stop Backend gRPC Server
		fenixExecutionWorkerObject.StopGrpcServer()

		//log.Println("Close DB_session: %v", DB_session)
		//DB_session.Close()
	}
}

func fenixExecutionWorkerMain() {

	// Connect to CloudDB
	fenixSyncShared.ConnectToDB()

	// Set up BackendObject
	fenixExecutionWorkerObject = &fenixExecutionWorkerObjectStruct{
		logger:                    nil,
		gcpAccessToken:            nil,
		executionEngineChannelRef: nil,
		executionEngine:           &testInstructionExecutionEngine.TestInstructionExecutionEngineStruct{},
	}

	// Init logger
	fenixExecutionWorkerObject.InitLogger("")

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	// Create Channel used for sending Commands to TestInstructionExecutionCommandsEngine
	testInstructionExecutionEngine.ExecutionEngineCommandChannel = make(chan testInstructionExecutionEngine.ChannelCommandStruct)
	myCommandChannelRef := &testInstructionExecutionEngine.ExecutionEngineCommandChannel
	fenixExecutionWorkerObject.executionEngineChannelRef = myCommandChannelRef

	// Initiate logger in TestInstructionEngine
	fenixExecutionWorkerObject.executionEngine.SetLogger(fenixExecutionWorkerObject.logger)

	// Start Receiver channel for Commands
	fenixExecutionWorkerObject.executionEngine.InitiateTestInstructionExecutionEngineCommandChannelReader(*myCommandChannelRef)

	// Start Backend gRPC-server
	fenixExecutionWorkerObject.InitGrpcServer()

}
