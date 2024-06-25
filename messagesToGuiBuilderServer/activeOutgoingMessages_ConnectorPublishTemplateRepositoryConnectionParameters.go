package messagesToGuiBuilderServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendConnectorPublishTemplateRepositoryConnectionParametersToFenixGuiBuilderServer
// Connector send template repository connection parameters to GuiBuilderServer
func (fenixExecutionWorkerObject *MessagesToGuiBuilderServerObjectStruct) SendConnectorPublishTemplateRepositoryConnectionParametersToFenixGuiBuilderServer(
	allTemplateRepositoryConnectionParameters *fenixTestCaseBuilderServerGrpcApi.
		AllTemplateRepositoryConnectionParameters) (bool, string) {

	fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "6837603a-06a8-4288-a6dc-02f832f7e3aa",
		//"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage,
	}).Debug("Incoming 'SendConnectorPublishTemplateRepositoryConnectionParametersToFenixGuiBuilderServer'")

	defer fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "63b5b963-c2fe-43f1-aed3-5576ab571d52",
	}).Debug("Outgoing 'SendConnectorPublishTemplateRepositoryConnectionParametersToFenixGuiBuilderServer'")

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string

	// Set up connection to BuilderServer, if that is not already done
	if fenixExecutionWorkerObject.connectionToGuiBuilderServerInitiated == false {
		err := fenixExecutionWorkerObject.SetConnectionToFenixGuiBuilderServer()
		if err != nil {
			return false, err.Error()
		}
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID": "06ab9268-8694-4807-bd80-920293c60664",
		}).Debug("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixGuiBuilderServer == common_config.GCP {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = fenixExecutionWorkerObject.generateGCPAccessToken(ctx)
		if returnMessageAckNack == false {
			return false, returnMessageString
		}

	}

	// Creates a new temporary client only to be used for this call
	var tempFenixGuiBuilderServerGrpcClient fenixTestCaseBuilderServerGrpcApi.FenixTestCaseBuilderServerGrpcWorkerServicesClient
	tempFenixGuiBuilderServerGrpcClient = fenixTestCaseBuilderServerGrpcApi.NewFenixTestCaseBuilderServerGrpcWorkerServicesClient(remoteFenixGuiBuilderServerConnection)

	// Do gRPC-call
	returnMessage, err := tempFenixGuiBuilderServerGrpcClient.ConnectorPublishTemplateRepositoryConnectionParameters(
		ctx,
		allTemplateRepositoryConnectionParameters)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":    "2ef2ed3d-4b24-4ab3-9240-16089939480c",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendConnectorPublishTemplateRepositoryConnectionParametersToFenixGuiBuilderServer'")

		// Set that a new connection needs to be done next time
		fenixExecutionWorkerObject.connectionToGuiBuilderServerInitiated = false

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":                                 "91bf5e4a-1824-46df-be9d-c55c3ac7bbca",
			"Message from FenixGuiBuilderServer": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendConnectorPublishTemplateRepositoryConnectionParametersToFenixGuiBuilderServer'")

		return false, returnMessage.Comments
	}

	return returnMessage.AckNack, returnMessage.Comments

}
