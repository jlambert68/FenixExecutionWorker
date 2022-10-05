package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
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

	// fmt.Println(processTestInstructionExecutionRequest) //TODO Remove

	s.logger.WithFields(logrus.Fields{
		"id":                                     "0909cb27-ab05-446b-9fe3-c36b05a6137b",
		"processTestInstructionExecutionRequest": processTestInstructionExecutionRequest,
	}).Debug("Received 'processTestInstructionExecutionRequest' from Execution Server")

	//  Check that TestInstructionExecutionUuid already is in Map
	_, existsInMap := processTestInstructionExecutionReversedResponseChannelMap[processTestInstructionExecutionRequest.TestInstruction.TestInstructionUuid]

	// Shouldn't exist in map
	if existsInMap == true {
		s.logger.WithFields(logrus.Fields{
			"id":                                     "df3ddde1-f55d-4d47-86bf-88626a6bb3ea",
			"processTestInstructionExecutionRequest": processTestInstructionExecutionRequest,
		}).Error("TestInstructionExecutionUuid already exists i 'processTestInstructionExecutionReversedResponseChannelMap'")

		// Generate response
		processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     fmt.Sprintf("TestInstructionExecutionUuid already exists i 'processTestInstructionExecutionReversedResponseChannelMap'"),
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

	// Create response channel to be able to get the 'processTestInstructionExecutionReversedResponse' back from Connector
	var processTestInstructionExecutionReversedResponseChannel processTestInstructionExecutionReversedResponseChannelType
	processTestInstructionExecutionReversedResponseChannel = make(chan *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse)

	// Create data for 'processTestInstructionExecutionReversedResponseChannelMap'
	var processTestInstructionExecutionReversedResponseMapData *processTestInstructionExecutionReversedResponseStruct
	processTestInstructionExecutionReversedResponseMapData = &processTestInstructionExecutionReversedResponseStruct{
		testInstructionExecutionUuid:                                    processTestInstructionExecutionRequest.TestInstruction.TestInstructionUuid,
		processTestInstructionExecutionReversedResponseChannelReference: processTestInstructionExecutionReversedResponseChannel,
		savedInMapTimeStamp:                                             time.Now(),
	}

	// Save 'processTestInstructionExecutionReversedResponseChannelData' in Map
	processTestInstructionExecutionReversedResponseChannelMap[processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid] = processTestInstructionExecutionReversedResponseMapData

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

		// Wait for response from 'processTestInstructionExecutionReversedResponseChannel'
		var testInstructionExecutionReversedResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse
		testInstructionExecutionReversedResponse = <-processTestInstructionExecutionReversedResponseChannel

		// Generate response
		processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
				AckNack:                      testInstructionExecutionReversedResponse.AckNackResponse.AckNack,
				Comments:                     testInstructionExecutionReversedResponse.AckNackResponse.Comments,
				ErrorCodes:                   testInstructionExecutionReversedResponse.AckNackResponse.ErrorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
			},
			TestInstructionExecutionUuid:   processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid,
			ExpectedExecutionDuration:      testInstructionExecutionReversedResponse.ExpectedExecutionDuration,
			TestInstructionCanBeReExecuted: testInstructionExecutionReversedResponse.TestInstructionCanBeReExecuted,
		}

	}

	// Remove message from Map
	delete(processTestInstructionExecutionReversedResponseChannelMap, processTestInstructionExecutionRequest.TestInstruction.TestInstructionUuid)

	return processTestInstructionExecutionResponse, nil

}
