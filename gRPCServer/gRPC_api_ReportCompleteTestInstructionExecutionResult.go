package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ReportCompleteTestInstructionExecutionResult
// Execution Server ask Worker (client) to report the final results of the execution result to the Server
func (s *fenixExecutionWorkerGrpcServicesServer) ReportCompleteTestInstructionExecutionResult(ctx context.Context, testInstructionExecutionRequestMessage *fenixExecutionWorkerGrpcApi.TestInstructionExecutionRequestMessage) (*fenixExecutionWorkerGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id": "b5fcd623-b81e-41b3-ad20-d89351dc0235",
	}).Debug("Incoming 'gRPC - ReportCompleteTestInstructionExecutionResult'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "8b38ac5e-0bf5-490c-a339-2e618d04983d",
	}).Debug("Outgoing 'gRPC - ReportCompleteTestInstructionExecutionResult'")

	// Calling system
	userId := "Execution Server"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(testInstructionExecutionRequestMessage.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		// Exiting
		return returnMessage, nil
	}

	//

	returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return returnMessage, nil
}
