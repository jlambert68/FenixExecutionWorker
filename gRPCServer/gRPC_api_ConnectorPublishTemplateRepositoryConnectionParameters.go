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

// ConnectorPublishTemplateRepositoryConnectionParameters
// Connector publish template repository connection parameters to Worker
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorPublishTemplateRepositoryConnectionParameters(
	ctx context.Context,
	allTemplateRepositoryConnectionParameters *fenixExecutionWorkerGrpcApi.
		AllTemplateRepositoryConnectionParameters) (
	returnMessage *fenixExecutionWorkerGrpcApi.AckNackResponse,
	err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "5d5bbb93-fa60-4d29-81ae-eb492441dd35",
		//"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage,
	}).Debug("Incoming 'gRPCWorker- ConnectorPublishTemplateRepositoryConnectionParameters'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "9a1a13e3-8912-417d-b150-8615c682db5f",
	}).Debug("Outgoing 'gRPCWorker - ConnectorPublishTemplateRepositoryConnectionParameters'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage = common_config.IsCallerUsingCorrectWorkerProtoFileVersion(
		userId,
		fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
			allTemplateRepositoryConnectionParameters.ClientSystemIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Set up instance to use for execution gPRC
	var fenixGuiBuilderObject *messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct
	fenixGuiBuilderObject = &messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct{Logger: s.logger}

	// Create gRPC-message towards GuiBuilderServer for 'AllTemplateRepositoryConnectionParameters'
	var allTemplateRepositoryConnectionParametersToBuilderServer *fenixTestCaseBuilderServerGrpcApi.
		AllTemplateRepositoryConnectionParameters
	allTemplateRepositoryConnectionParametersToBuilderServer = &fenixTestCaseBuilderServerGrpcApi.
		AllTemplateRepositoryConnectionParameters{
		ClientSystemIdentification: &fenixTestCaseBuilderServerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid: allTemplateRepositoryConnectionParameters.GetClientSystemIdentification().
				GetDomainUuid(),
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestBuilderServerProtoFileVersion()),
		},
		AllTemplateRepositories: nil,
	}

	// Convert incoming message to be used for outgoing 'AllTemplateRepositories'
	var allTemplateRepositories []*fenixTestCaseBuilderServerGrpcApi.TemplateRepositoryConnectionParameters
	for _, templateRepositoryConnectionParameters := range allTemplateRepositoryConnectionParameters.GetAllTemplateRepositories() {

		var templateRepositoryConnectionParametersToTestCaseBuilderServer *fenixTestCaseBuilderServerGrpcApi.
			TemplateRepositoryConnectionParameters

		templateRepositoryConnectionParametersToTestCaseBuilderServer = &fenixTestCaseBuilderServerGrpcApi.
			TemplateRepositoryConnectionParameters{
			RepositoryApiUrlName: templateRepositoryConnectionParameters.RepositoryApiUrlName,
			RepositoryApiUrl:     templateRepositoryConnectionParameters.RepositoryApiUrl,
			RepositoryOwner:      templateRepositoryConnectionParameters.RepositoryOwner,
			RepositoryName:       templateRepositoryConnectionParameters.RepositoryName,
			RepositoryPath:       templateRepositoryConnectionParameters.RepositoryPath,
			GitHubApiKey:         templateRepositoryConnectionParameters.GitHubApiKey,
		}

		allTemplateRepositories = append(allTemplateRepositories, templateRepositoryConnectionParametersToTestCaseBuilderServer)
	}
	// Add converted messages to outgoing message to BuilderServer
	allTemplateRepositoryConnectionParametersToBuilderServer.AllTemplateRepositories = allTemplateRepositories

	// Create signature message
	var messageSignatureData *fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage
	messageSignatureData = &fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage{
		HashToBeSigned: allTemplateRepositoryConnectionParameters.GetMessageSignatureData().GetHashToBeSigned(),
		Signature:      allTemplateRepositoryConnectionParameters.GetMessageSignatureData().GetSignature(),
	}

	// Save the Signature data in message to TestCaseBuilderServer
	allTemplateRepositoryConnectionParametersToBuilderServer.MessageSignatureData = messageSignatureData

	// Publish Supported template repository connection parameters to TestCaseBuilderServer
	var succeededToSend bool
	var responseMessage string
	succeededToSend, responseMessage = fenixGuiBuilderObject.
		SendConnectorPublishTemplateRepositoryConnectionParametersToFenixGuiBuilderServer(
			allTemplateRepositoryConnectionParametersToBuilderServer)

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
