package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/messagesToGuiBuilderServer"
	"context"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/shared_code"
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
			RepositoryApiUrl: templateRepositoryConnectionParameters.RepositoryApiUrl,
			RepositoryOwner:  templateRepositoryConnectionParameters.RepositoryOwner,
			RepositoryName:   templateRepositoryConnectionParameters.RepositoryName,
			RepositoryPath:   templateRepositoryConnectionParameters.RepositoryPath,
			GitHubApiKey:     templateRepositoryConnectionParameters.GitHubApiKey,
		}

		allTemplateRepositories = append(allTemplateRepositories, templateRepositoryConnectionParametersToTestCaseBuilderServer)
	}
	// Add converted messages to outgoing message to BuilderServer
	allTemplateRepositoryConnectionParametersToBuilderServer.AllTemplateRepositories = allTemplateRepositories

	// Get Message to sign to prove identity
	succeededToSend, responseMessage, messageToSign := fenixGuiBuilderObject.
		SendGetMessageToSignToProveCallerIdentity()

	if succeededToSend == false {
		s.logger.WithFields(logrus.Fields{
			"id":              "4199ea85-debd-4b49-93d1-2e42d7611122",
			"responseMessage": responseMessage,
		}).Error("Got some error when sending 'GetMessageToSignToProveCallerIdentity'")

		// Generate response
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      succeededToSend,
			Comments:                     fmt.Sprintf("Messagage from BuilderServer: '%s'", responseMessage),
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil

	}

	// Specify the service account to be used when signing
	var serviceAccountUsedWhenSigning string
	serviceAccountUsedWhenSigning = fmt.Sprintf("projects/-/serviceAccounts/%s",
		common_config.ServiceAccountUsedForSigningMessage)

	// Sign Message to prove Identity to BuilderServer
	var hashOfSignature string
	var hashedKeyId string
	if common_config.ExecutionLocationForWorker == common_config.GCP {

		hashOfSignature, hashedKeyId, err = shared_code.SignMessageToProveIdentityToBuilderServer(
			messageToSign,
			serviceAccountUsedWhenSigning,
			true)

	} else {

		hashOfSignature, hashedKeyId, err = shared_code.SignMessageToProveIdentityToBuilderServer(
			messageToSign,
			serviceAccountUsedWhenSigning,
			false)
	}

	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"id":  "bb8284b2-cc4f-440f-8474-a2ea12adceb1",
			"err": err,
		}).Error("Got some error when signing message")

		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     fmt.Sprintf("Got some error when signing message. '%s'", err.Error()),
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, err
	}

	// Add Signed message
	var signedMessageByWorkerServiceAccountMessage *fenixTestCaseBuilderServerGrpcApi.SignedMessageByWorkerServiceAccountMessage
	signedMessageByWorkerServiceAccountMessage = &fenixTestCaseBuilderServerGrpcApi.SignedMessageByWorkerServiceAccountMessage{
		MessageToBeSigned: messageToSign,
		HashOfSignature:   hashOfSignature,
		HashedKeyId:       hashedKeyId,
	}
	allTemplateRepositoryConnectionParametersToBuilderServer.
		SignedMessageByWorkerServiceAccount = signedMessageByWorkerServiceAccountMessage

	// Publish Supported template repository connection parameters to TestCaseBuilderServer
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
