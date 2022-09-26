package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	//"FenixExecutionWorker"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// AreYouAlive - *********************************************************************
//Anyone can check if Fenix Execution Worker server is alive with this service
func (s *fenixExecutionWorkerGrpcServicesServer) AreYouAlive(ctx context.Context, emptyParameter *fenixExecutionWorkerGrpcApi.EmptyParameter) (*fenixExecutionWorkerGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id": "1ff67695-9a8b-4821-811d-0ab8d33c4d8b",
	}).Debug("Incoming 'gRPCServer - AreYouAlive'")

	s.logger.WithFields(logrus.Fields{
		"id": "9c7f0c3d-7e9f-4c91-934e-8d7a22926d84",
	}).Debug("Outgoing 'gRPCServer - AreYouAlive'")

	return &fenixExecutionWorkerGrpcApi.AckNackResponse{AckNack: true, Comments: "I'am alive."}, nil
}

// ReportProcessingCapability
// Ask Client to inform Execution Server of Clients capability to execute requests in parallell, serial or no processing at all
func (s *fenixExecutionWorkerGrpcServicesServer) ReportProcessingCapability(ctx context.Context, emptyParameter *fenixExecutionWorkerGrpcApi.EmptyParameter) (*fenixExecutionWorkerGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id": "37bc2356-33a2-4e2c-9420-122df581d757",
	}).Debug("Incoming 'gRPCServer - ReportProcessingCapability'")

	s.logger.WithFields(logrus.Fields{
		"id": "f3fd3e50-5770-48ad-8524-85f34d28545e",
	}).Debug("Outgoing 'gRPCServer - ReportProcessingCapability'")

	// Calling system
	userId := "Execution Server"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(emptyParameter.ProtoFileVersionUsedByClient))
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

// ProcessTestInstructionExecution
// Fenix Execution Server send a request to Execution Worker to initiate a execution of a TestInstruction
func (s *fenixExecutionWorkerGrpcServicesServer) ProcessTestInstructionExecution(ctx context.Context, processTestInstructionExecutionRequest *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionRequest) (processTestInstructionExecutionResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "37bc2356-33a2-4e2c-9420-122df581d757",
	}).Debug("Incoming 'gRPCServer - ProcessTestInstructionExecution'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "f3fd3e50-5770-48ad-8524-85f34d28545e",
	}).Debug("Outgoing 'gRPCServer - ProcessTestInstructionExecution'")

	// Calling system
	userId := "Execution Server"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(processTestInstructionExecutionRequest.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		// Exiting
		processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
			AckNackResponse:                returnMessage,
			TestInstructionExecutionUuid:   "",
			ExpectedExecutionDuration:      nil,
			TestInstructionCanBeReExecuted: false,
		}

		return processTestInstructionExecutionResponse, nil
	}

	// Generate response
	processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
		AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(common_config.GetHighestFenixTestDataProtoFileVersion()),
		},
		TestInstructionExecutionUuid:   "",
		ExpectedExecutionDuration:      nil,
		TestInstructionCanBeReExecuted: false,
	}

	return processTestInstructionExecutionResponse, nil

}
