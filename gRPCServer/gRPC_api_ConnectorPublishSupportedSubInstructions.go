package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/messagesToGuiBuilderServer"
	"context"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ConnectorPublishSupportedSubInstructions
// Connector publish supported SubInstructions to Worker
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorPublishSupportedSubInstructions(
	ctx context.Context,
	supportedSubInstructionsMessage *fenixExecutionWorkerGrpcApi.SupportedSubInstructions) (
	returnMessage *fenixExecutionWorkerGrpcApi.AckNackResponse,
	err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "a1ecfa84-fc73-43b7-b644-869f4bed0d6c",
		//"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage,
	}).Debug("Incoming 'gRPCWorker- ConnectorPublishSupportedSubInstructions'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "85dc9444-27e3-407f-a480-bc9965c46b12",
	}).Debug("Outgoing 'gRPCWorker - ConnectorPublishSupportedSubInstructions'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage = common_config.IsCallerUsingCorrectWorkerProtoFileVersion(
		userId,
		fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
			supportedSubInstructionsMessage.ClientSystemIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Set up instance to use for execution gPRC
	var fenixGuiBuilderObject *messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct
	fenixGuiBuilderObject = &messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct{Logger: s.logger}

	// Create gRPC-message towards GuiBuilderServer for 'SupportedTestCaseMetaDataMessage'
	var supportedSubInstructionsMessageToBuilderServer *fenixTestCaseBuilderServerGrpcApi.
		SupportedSubInstructions
	supportedSubInstructionsMessageToBuilderServer = &fenixTestCaseBuilderServerGrpcApi.
		SupportedSubInstructions{
		ClientSystemIdentification: &fenixTestCaseBuilderServerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid: supportedSubInstructionsMessage.GetClientSystemIdentification().
				GetDomainUuid(),
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestBuilderServerProtoFileVersion()),
		},
		SupportedSubInstructionsAsJson: supportedSubInstructionsMessage.
			GetSupportedSubInstructionsAsJson(),
		SupportedSubInstructionsPerTestInstructionAsJson: supportedSubInstructionsMessage.
			GetSupportedSubInstructionsPerTestInstructionAsJson(),
		MessageSignatureData: nil,
	}

	// Create signature message
	var messageSignatureData *fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage
	messageSignatureData = &fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage{
		HashToBeSigned: supportedSubInstructionsMessage.GetMessageSignatureData().GetHashToBeSigned(),
		Signature:      supportedSubInstructionsMessage.GetMessageSignatureData().GetSignature(),
	}

	// Save the Signature data in message to TestCaseBuilderServer
	supportedSubInstructionsMessageToBuilderServer.MessageSignatureData = messageSignatureData

	// Publish Supported template repository connection parameters to TestCaseBuilderServer
	var succeededToSend bool
	var responseMessage string
	succeededToSend, responseMessage = fenixGuiBuilderObject.
		SendConnectorPublishSupportedSubInstructionsToFenixGuiBuilderServer(
			supportedSubInstructionsMessageToBuilderServer)

	// Create Error Codes
	var errorCodes []fenixExecutionWorkerGrpcApi.ErrorCodesEnum

	// Generate response
	returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      succeededToSend,
		Comments:                     fmt.Sprintf("Messagage from BuilderServer: '%s'", responseMessage),
		ErrorCodes:                   errorCodes,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return returnMessage, nil

}
