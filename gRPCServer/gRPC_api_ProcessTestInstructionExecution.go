package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// ProcessTestInstructionExecution
// Fenix Execution Server send a request to Execution Worker to initiate a execution of a TestInstruction
func (s *fenixExecutionWorkerGrpcServicesServer) ProcessTestInstructionExecution(ctx context.Context, processTestInstructionExecutionRequest *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest) (processTestInstructionExecutionResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "37bc2356-33a2-4e2c-9420-122df581d757",
	}).Debug("Incoming 'gRPCServer - ProcessTestInstructionExecution'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "f3fd3e50-5770-48ad-8524-85f34d28545e",
	}).Debug("Outgoing 'gRPCServer - ProcessTestInstructionExecution'")

	// Calling system
	userId := "Execution Server"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(processTestInstructionExecutionRequest.ProtoFileVersionUsedByClient))
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

	fmt.Println(processTestInstructionExecutionRequest)

	// Generate duration for Execution:: TODO This is only for test and should be done in another way lator
	executionDuration := time.Minute * 5
	timeAtDurationEnd := time.Now().Add(executionDuration)

	// Generate response
	processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
		AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		},
		TestInstructionExecutionUuid:   processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid,
		ExpectedExecutionDuration:      timestamppb.New(timeAtDurationEnd),
		TestInstructionCanBeReExecuted: false,
	}

	return processTestInstructionExecutionResponse, nil

}
