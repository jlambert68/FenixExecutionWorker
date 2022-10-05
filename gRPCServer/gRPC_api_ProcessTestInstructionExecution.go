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

	// If there isn't an active connection to the Connector then  report that back
	if connectorHasConnected == false {

		// Generate response
		processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     fmt.Sprintf("Message couldn't be sent to Connector, due to no active Connector was found"),
				ErrorCodes:                   nil,
				ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
			},
			TestInstructionExecutionUuid:   processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid,
			ExpectedExecutionDuration:      nil,
			TestInstructionCanBeReExecuted: false,
		}

		// Return Response to Execution Server
		return processTestInstructionExecutionResponse, nil

	}

	fmt.Println(processTestInstructionExecutionRequest) //TODO Remove
	s.logger.WithFields(logrus.Fields{
		"id":                                     "0909cb27-ab05-446b-9fe3-c36b05a6137b",
		"processTestInstructionExecutionRequest": processTestInstructionExecutionRequest,
	}).Debug("Received 'processTestInstructionExecutionRequest' from Execution Server")

	// Create response channel to be able to get response when TestInstructionExecution is sent to Connector
	var executionResponseChannel executionResponseChannelType

	// Create message to be sent to stream-server
	executionForwardChannelMessage := executionForwardChannelStruct{
		processTestInstructionExecutionReveredRequest: processTestInstructionExecutionRequest,
		executionResponseChannelReference:             &executionResponseChannel,
		isKeepAliveMessage:                            false,
	}

	// Send TestInstructionExecution to Stream-server, to later be sent to Connector, over channel
	executionForwardChannel <- executionForwardChannelMessage

	// Wait for response from stream-server that message has been sent TODO create some maximum time before clearing channel
	executionResponseChannelMessage := <-executionResponseChannel

	if executionResponseChannelMessage.testInstructionExecutionIsSentToConnector == false ||
		executionResponseChannelMessage.err != nil {
		// Message failed to be sent to Connector

		// Generate response
		processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     fmt.Sprintf("Message couldn't be sent to Connector, error: '%s'", err.Error()),
				ErrorCodes:                   nil,
				ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
			},
			TestInstructionExecutionUuid:   processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid,
			ExpectedExecutionDuration:      nil,
			TestInstructionCanBeReExecuted: false,
		}

	} else {
		// Message succeeded to be sent to Connector

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
	}

	return processTestInstructionExecutionResponse, nil

}
