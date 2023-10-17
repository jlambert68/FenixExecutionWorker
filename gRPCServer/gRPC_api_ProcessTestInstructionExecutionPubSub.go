package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixSyncShared/pubSubHelpers"
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

	// Convert PubSub-message back into proto-message
	var processTestInstructionExecutionPubSubRequest2 fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionPubSubRequest
	err2 := protojson.Unmarshal([]byte(processTestInstructionExecutionRequestAsJsonString), &processTestInstructionExecutionPubSubRequest2)
	if err2 != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "5be52325-5862-45f4-8ef8-ced518e11c8c",
			"Error": err2,
			"processTestInstructionExecutionRequestAsJsonString": processTestInstructionExecutionRequestAsJsonString,
			"processTestInstructionExecutionPubSubRequest":       processTestInstructionExecutionPubSubRequest,
			"processTestInstructionExecutionPubSubRequest2":      processTestInstructionExecutionPubSubRequest2,
		}).Error("Something went wrong when converting 'PubSub-message into proto-message")

		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   err2.Error(),
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
				common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil
	}
	common_config.Logger.WithFields(logrus.Fields{
		"Id": "01a55d6b-311d-4bdd-8284-f2f4ac8e582a",
		"processTestInstructionExecutionRequestAsJsonString": processTestInstructionExecutionRequestAsJsonString,
		"processTestInstructionExecutionPubSubRequest":       processTestInstructionExecutionPubSubRequest,
		"processTestInstructionExecutionPubSubRequest2":      processTestInstructionExecutionPubSubRequest2,
	}).Debug("Message before and after json-convert")

	var (
		err                  error
		returnMessageAckNack bool
		returnMessageString  string
	)

	// Create PubSub-Topic
	var pubSubTopicToLookFor string
	pubSubTopicToLookFor = GeneratePubSubTopicForTestInstructionExecutions()

	// Publish TestInstructionExecution on PubSub
	//returnMessageAckNack, returnMessageString, err = outgoingPubSubMessages.Publish(
	//	processTestInstructionExecutionRequestAsJsonString)
	returnMessageAckNack, returnMessageString, err = pubSubHelpers.PublishExecutionStatusOnPubSub(
		pubSubTopicToLookFor, processTestInstructionExecutionRequestAsJsonString)

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

}

// Create the PubSub-topic from Domain-Uuid
func GeneratePubSubTopicForTestInstructionExecutions() (statusExecutionTopic string) {

	var pubSubTopicBase string
	pubSubTopicBase = common_config.TestExecutionStatusPubSubTopicBase

	var testerGuiApplicationUuid string
	testerGuiApplicationUuid = common_config.ThisDomainsUuid

	// Get the first 8 characters from TesterGui-ApplicationUuid
	var shortedAppUuid string
	shortedAppUuid = testerGuiApplicationUuid[0:8]

	// Build PubSub-topic
	statusExecutionTopic = pubSubTopicBase + "-" + shortedAppUuid

	return statusExecutionTopic
}
