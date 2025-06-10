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
	supportedMetaDataMessage *fenixExecutionWorkerGrpcApi.SupportedTestCaseAndTestSuiteMetaData) (
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
			supportedMetaDataMessage.ClientSystemIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Set up instance to use for execution gPRC
	var fenixGuiBuilderObject *messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct
	fenixGuiBuilderObject = &messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct{Logger: s.logger}

	// Create gRPC-message towards GuiBuilderServer for 'SupportedTestCaseMetaDataMessage'
	var supportedMetaDataMessageMessageToBuilderServer *fenixTestCaseBuilderServerGrpcApi.
		SupportedTestCaseAndTestSuiteMetaData
	supportedMetaDataMessageMessageToBuilderServer = &fenixTestCaseBuilderServerGrpcApi.
		SupportedTestCaseAndTestSuiteMetaData{
		ClientSystemIdentification: &fenixTestCaseBuilderServerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid: supportedMetaDataMessage.GetClientSystemIdentification().
				GetDomainUuid(),
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestBuilderServerProtoFileVersion()),
		},
		SupportedTestCaseMetaDataAsJson:  supportedMetaDataMessage.GetSupportedTestCaseMetaDataAsJson(),
		SupportedTestSuiteMetaDataAsJson: supportedMetaDataMessage.GetSupportedTestSuiteMetaDataAsJson(),
		MessageSignatureData:             nil,
	}

	// Create signature message
	var messageSignatureData *fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage
	messageSignatureData = &fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage{
		HashToBeSigned: supportedMetaDataMessage.GetMessageSignatureData().GetHashToBeSigned(),
		Signature:      supportedMetaDataMessage.GetMessageSignatureData().GetSignature(),
	}

	// Save the Signature data in message to TestCaseBuilderServer
	supportedMetaDataMessageMessageToBuilderServer.MessageSignatureData = messageSignatureData

	// Publish Supported template repository connection parameters to TestCaseBuilderServer
	var succeededToSend bool
	var responseMessage string
	succeededToSend, responseMessage = fenixGuiBuilderObject.
		SendConnectorPublishSupportedMetaDataToFenixGuiBuilderServer(
			supportedMetaDataMessageMessageToBuilderServer)

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
