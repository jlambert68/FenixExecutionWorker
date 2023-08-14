package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/messagesToExecutionServer"
	"context"
	"fmt"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ConnectorProcessTestInstructionExecutionResponse
// Response from execution client to execution Worker using direct gRPC call that Client(Connector) has taken care of TestInstructionExecution
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorProcessTestInstructionExecutionResponse(ctx context.Context, processTestInstructionExecutionResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse) (ackNackResponse *fenixExecutionWorkerGrpcApi.AckNackResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "47e23f54-bf01-4504-af26-a27243e7c51b",
		"processTestInstructionExecutionResponse": processTestInstructionExecutionResponse,
	}).Debug("Incoming 'gRPCServer - ConnectorProcessTestInstructionExecutionResponse'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "c0d50288-fdac-4dfc-a5a5-1cbbdc04d1bc",
	}).Debug("Outgoing 'gRPCServer - ConnectorProcessTestInstructionExecutionResponse'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, processTestInstructionExecutionResponse.AckNackResponse.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Set up instance to use for execution gPRC
	var fenixExecutionWorkerObject *messagesToExecutionServer.MessagesToExecutionServerObjectStruct
	fenixExecutionWorkerObject = &messagesToExecutionServer.MessagesToExecutionServerObjectStruct{Logger: s.logger}

	// Create 'ProcessTestInstructionExecutionResponseStatus'
	var processTestInstructionExecutionResponseStatusToExecutionServerMessage *fenixExecutionServerGrpcApi.ProcessTestInstructionExecutionResponseStatus
	processTestInstructionExecutionResponseStatusToExecutionServerMessage = &fenixExecutionServerGrpcApi.ProcessTestInstructionExecutionResponseStatus{
		AckNackResponse: &fenixExecutionServerGrpcApi.AckNackResponse{
			AckNack:                      processTestInstructionExecutionResponse.AckNackResponse.AckNack,
			Comments:                     processTestInstructionExecutionResponse.AckNackResponse.Comments,
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionServerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(common_config.GetHighestFenixExecutionServerProtoFileVersion()),
		},
		TestInstructionExecutionUuid:   processTestInstructionExecutionResponse.TestInstructionExecutionUuid,
		ExpectedExecutionDuration:      processTestInstructionExecutionResponse.ExpectedExecutionDuration,
		TestInstructionCanBeReExecuted: false,
	}

	succeededToSend, responseMessage := fenixExecutionWorkerObject.SendProcessResponseTestInstructionExecution(processTestInstructionExecutionResponseStatusToExecutionServerMessage)

	if succeededToSend == false {
		s.logger.WithFields(logrus.Fields{
			"id":              "dd3749c7-9bd9-435b-8966-f6bd134b1dd2",
			"responseMessage": responseMessage,
		}).Error("Got some error when sending 'processTestInstructionExecutionResponseStatusToExecutionServerMessage'")
	}

	// Create Error Codes
	var errorCodes []fenixExecutionWorkerGrpcApi.ErrorCodesEnum

	// Generate response
	ackNackResponse = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      succeededToSend,
		Comments:                     fmt.Sprintf("Messagage from ExecutionServer: '%s'", responseMessage),
		ErrorCodes:                   errorCodes,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return ackNackResponse, nil

}
