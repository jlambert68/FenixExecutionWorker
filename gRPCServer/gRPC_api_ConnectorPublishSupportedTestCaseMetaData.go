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

// ConnectorPublishSupportedMetaData
// Connector publish supported TestCaseMetaData to Worker
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorPublishSupportedMetaData(
	ctx context.Context,
	supportedTestCaseMetaDataMessage *fenixExecutionWorkerGrpcApi.SupportedTestCaseMetaData) (
	returnMessage *fenixExecutionWorkerGrpcApi.AckNackResponse,
	err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "96038f1a-a58f-4025-9bd7-9fe47d700d2a",
		//"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage,
	}).Debug("Incoming 'gRPCWorker- ConnectorPublishSupportedMetaData'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "f1862521-1349-48e1-ab9f-601317675710",
	}).Debug("Outgoing 'gRPCWorker - ConnectorPublishSupportedMetaData'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage = common_config.IsCallerUsingCorrectWorkerProtoFileVersion(
		userId,
		fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
			supportedTestCaseMetaDataMessage.ClientSystemIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Set up instance to use for execution gPRC
	var fenixGuiBuilderObject *messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct
	fenixGuiBuilderObject = &messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct{Logger: s.logger}

	// Create gRPC-message towards GuiBuilderServer for 'SupportedTestCaseMetaDataMessage'
	var supportedTestCaseMetaDataMessageMessageToBuilderServer *fenixTestCaseBuilderServerGrpcApi.
		SupportedTestCaseMetaData
	supportedTestCaseMetaDataMessageMessageToBuilderServer = &fenixTestCaseBuilderServerGrpcApi.
		SupportedTestCaseMetaData{
		ClientSystemIdentification: &fenixTestCaseBuilderServerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid: supportedTestCaseMetaDataMessage.GetClientSystemIdentification().
				GetDomainUuid(),
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestBuilderServerProtoFileVersion()),
		},
		SupportedMetaDataAsJson: supportedTestCaseMetaDataMessage.GetSupportedMetaDataAsJson(),
		MessageSignatureData:    nil,
	}

	// Create signature message
	var messageSignatureData *fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage
	messageSignatureData = &fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage{
		HashToBeSigned: supportedTestCaseMetaDataMessage.GetMessageSignatureData().GetHashToBeSigned(),
		Signature:      supportedTestCaseMetaDataMessage.GetMessageSignatureData().GetSignature(),
	}

	// Save the Signature data in message to TestCaseBuilderServer
	supportedTestCaseMetaDataMessageMessageToBuilderServer.MessageSignatureData = messageSignatureData

	// Publish Supported template repository connection parameters to TestCaseBuilderServer
	var succeededToSend bool
	var responseMessage string
	succeededToSend, responseMessage = fenixGuiBuilderObject.
		SendConnectorPublishSupportedTestCaseMetaDataToFenixGuiBuilderServer(
			supportedTestCaseMetaDataMessageMessageToBuilderServer)

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
