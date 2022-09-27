package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// SendAllLogPostForExecution
// Execution Server ask Worker (client) to report all log posts of the execution result to the Server
func (s *fenixExecutionWorkerGrpcServicesServer) SendAllLogPostForExecution(ctx context.Context, testInstructionExecutionRequestMessage *fenixExecutionWorkerGrpcApi.TestInstructionExecutionRequestMessage) (*fenixExecutionWorkerGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id": "7c98ab59-eb2a-4cfb-83da-668ca8c1e5c2",
	}).Debug("Incoming 'gRPC - SendAllLogPostForExecution'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "eeef51aa-30f6-49cd-ad8e-40e231f90398",
	}).Debug("Outgoing 'gRPC - SendAllLogPostForExecution'")

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
