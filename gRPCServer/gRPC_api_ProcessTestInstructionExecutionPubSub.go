package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/outgoingPubSubMessages"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

// ProcessTestInstructionExecutionPubSub
// Fenix Execution Server send a request to Execution Worker to initiate an execution of a TestInstruction
func (s *fenixExecutionWorkerGrpcServicesServer) ProcessTestInstructionExecutionPubSub(
	ctx context.Context,
	processTestInstructionExecutionPubSubRequest *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionPubSubRequest) (
	*fenixExecutionWorkerGrpcApi.AckNackResponse,
	error) {

	s.logger.WithFields(logrus.Fields{
		"id": "89844e62-4389-4ae8-aefd-8b1f44a1a3fc",
		"processTestInstructionExecutionPubSubRequest": processTestInstructionExecutionPubSubRequest,
	}).Debug("Incoming 'gRPC - ProcessTestInstructionExecutionPubSub'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "1ef7f394-6a6d-49be-a398-20b69ec58594",
	}).Debug("Outgoing 'gRPC - ProcessTestInstructionExecutionPubSub'")

	// Calling system
	userId := "Execution Server"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectWorkerProtoFileVersion(
		userId,
		fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
			processTestInstructionExecutionPubSubRequest.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		// Exiting
		return returnMessage, nil
	}

	// Convert gRPC-message into json-string
	var processTestInstructionExecutionRequestAsJsonString string
	processTestInstructionExecutionRequestAsJsonString = protojson.Format(processTestInstructionExecutionPubSubRequest)

	var (
		returnMessageAckNack bool
		returnMessageString  string
		err                  error
	)
	// Publish TestInstructionExecution on PubSub
	returnMessageAckNack, returnMessageString, err = outgoingPubSubMessages.Publish(
		processTestInstructionExecutionRequestAsJsonString)

	// Some problem when sending over PubSub
	if returnMessageAckNack == false || err != nil {
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   returnMessageString + "-" + err.Error(),
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
				common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil
	}

	// Sending over PubSub succeeded
	returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:    true,
		Comments:   "",
		ErrorCodes: nil,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
			common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return returnMessage, nil

	return returnMessage, nil
}
