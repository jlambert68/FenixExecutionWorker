package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ReportCurrentTestInstructionExecutionResult
// Execution Server ask Worker (client) to report the ongoing results of the execution result to the Server
func (s *fenixExecutionWorkerGrpcServicesServer) ReportCurrentTestInstructionExecutionResult(ctx context.Context, testInstructionExecutionRequestMessage *fenixExecutionWorkerGrpcApi.TestInstructionExecutionRequestMessage) (*fenixExecutionWorkerGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id": "19b9dfce-8f53-4ff4-9558-f1cf8d871f9e",
	}).Debug("Incoming 'gRPC - ReportCurrentTestInstructionExecutionResult'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "f3cb8f70-8239-4dbe-a02d-3602b875cbf8",
	}).Debug("Outgoing 'gRPC - ReportCurrentTestInstructionExecutionResult'")

	// Calling system
	userId := "Execution Server"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(testInstructionExecutionRequestMessage.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		// Exiting
		return returnMessage, nil
	}

	returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(common_config.GetHighestFenixTestDataProtoFileVersion()),
	}

	return returnMessage, nil
}
