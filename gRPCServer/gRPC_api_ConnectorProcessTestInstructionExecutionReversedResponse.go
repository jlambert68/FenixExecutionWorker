package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ConnectorReportCompleteTestInstructionExecutionResult
// When a TestInstruction has been fully executed the Execution Connector use this to inform the results of the execution result to the Worker
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorProcessTestInstructionExecutionReversedResponse(ctx context.Context, processTestInstructionExecutionReversedResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse) (ackNackResponse *fenixExecutionWorkerGrpcApi.AckNackResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "a0b241e2-1fa8-4b2d-990a-e238a182b869",
	}).Debug("Incoming 'gRPCServer - ConnectorProcessTestInstructionExecutionReversedResponse'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "125973c6-7a89-481c-b4ff-21501a689eba",
	}).Debug("Outgoing 'gRPCServer - ConnectorProcessTestInstructionExecutionReversedResponse'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(processTestInstructionExecutionReversedResponse.AckNackResponse.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Extract response channel
	processTestInstructionExecutionReversedResponseData, existInMap := processTestInstructionExecutionReversedResponseChannelMap[processTestInstructionExecutionReversedResponse.TestInstructionExecutionUuid]

	// Shouldn't exist in map
	if existInMap == false {
		s.logger.WithFields(logrus.Fields{
			"id": "7ab10f1b-266a-4684-abc9-94362b6b8fdb",
			"processTestInstructionExecutionReversedResponse": processTestInstructionExecutionReversedResponse,
		}).Error("TestInstructionExecutionUuid couldn't be found in 'processTestInstructionExecutionReversedResponseChannelMap'")

		// Generate response
		ackNackResponse = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     fmt.Sprintf("TestInstructionExecutionUuid '%s' couldn't be found in 'processTestInstructionExecutionReversedResponseChannelMap'", processTestInstructionExecutionReversedResponse.TestInstructionExecutionUuid),
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		// Return Response to Connector
		return ackNackResponse, nil

	}

	// Extract Channel
	testInstructionExecutionReversedResponseChannelReference := processTestInstructionExecutionReversedResponseData.processTestInstructionExecutionReversedResponseChannelReference
	var testInstructionExecutionReversedResponseChannel processTestInstructionExecutionReversedResponseChannelType
	testInstructionExecutionReversedResponseChannel = *testInstructionExecutionReversedResponseChannelReference

	// Send response over channel to gRPC-function 'ProcessTestInstructionExecution' where it is used
	testInstructionExecutionReversedResponseChannel <- processTestInstructionExecutionReversedResponse

	// Generate response to Connector
	ackNackResponse = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return ackNackResponse, nil

}
