package workerEngine

import (
	"fmt"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// Channel reader which is used for reading out commands to CommandEngine
func (executionEngine *TestInstructionExecutionEngineStruct) startCommandChannelReader() {

	var incomingChannelCommand ChannelCommandStruct

	for {
		// Wait for incoming command over channel
		incomingChannelCommand = <-*executionEngine.CommandChannelReference

		switch incomingChannelCommand.ChannelCommand {

		case ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServer:
			executionEngine.SendReportCompleteTestInstructionExecutionResultToFenixExecutionServer(incomingChannelCommand)

		// No other command is supported
		default:
			executionEngine.logger.WithFields(logrus.Fields{
				"Id":                     "6bf37452-da99-4e7e-aa6a-4627b05d1bdb",
				"incomingChannelCommand": incomingChannelCommand,
			}).Fatalln("Unknown command in CommandChannel for Worker Engine")
		}
	}

}

// Check ExecutionQueue for TestInstructions and move them to ongoing Executions-table
func (executionEngine *TestInstructionExecutionEngineStruct) initiateExecutionsForTestInstructionsOnExecutionQueue() {

	fmt.Println("initiateExecutionsForTestInstructionsOnExecutionQueue")

}

// Check ongoing executions  for TestInstructions for change in status that should be propagated to other places
func (executionEngine *TestInstructionExecutionEngineStruct) checkOngoingExecutionsForTestInstructions() {

}

// SendReportCompleteTestInstructionExecutionResultToFenixExecutionServer
// Forward the final result of a TestInstructionExecution done by domains own execution engine
func (executionEngine *TestInstructionExecutionEngineStruct) SendReportCompleteTestInstructionExecutionResultToFenixExecutionServer(channelCommand ChannelCommandStruct) {
	var finalTestInstructionExecutionResultMessageFromExecutionWorker *fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage
	var finalTestInstructionExecutionResultMessageToExecutionServer *fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage

	// Convert message into Worker-message-structure-type
	finalTestInstructionExecutionResultMessageFromExecutionWorker = channelCommand.ReportCompleteTestInstructionExecutionResultParameter.finalTestInstructionExecutionResultMessage

	// Convert from Worker-message into ExecutionServer-message
	finalTestInstructionExecutionResultMessageToExecutionServer = &fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage{
		ClientSystemIdentification: &fenixExecutionServerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid:                   finalTestInstructionExecutionResultMessageFromExecutionWorker.ClientSystemIdentification.DomainUuid,
			ProtoFileVersionUsedByClient: fenixExecutionServerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(finalTestInstructionExecutionResultMessageFromExecutionWorker.ClientSystemIdentification.ProtoFileVersionUsedByClient),
		},
		TestInstructionExecutionUuid:   finalTestInstructionExecutionResultMessageFromExecutionWorker.TestInstructionExecutionUuid,
		TestInstructionExecutionStatus: fenixExecutionServerGrpcApi.TestInstructionExecutionStatusEnum(finalTestInstructionExecutionResultMessageFromExecutionWorker.TestInstructionExecutionStatus),
	}

	// Send the result using a go-routine to be able to process next command on command-queue
	go func() {
		sendResult, errorMessage := executionEngine.messagesToExecutionServerObjectReference.SendReportCompleteTestInstructionExecutionResultToFenixExecutionServer(finalTestInstructionExecutionResultMessageToExecutionServer)

		if sendResult == false {
			executionEngine.logger.WithFields(logrus.Fields{
				"id":             "e9aae7c6-8a14-4da2-8001-2029d5bbac8d",
				"errorMessage":   errorMessage,
				"channelCommand": channelCommand,
			}).Error("Couldn't do gRPC-call to Execution Server ('SendReportCompleteTestInstructionExecutionResultToFenixExecutionServer')")
		}
	}()

}
