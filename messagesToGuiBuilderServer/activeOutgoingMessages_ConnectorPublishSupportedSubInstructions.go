package messagesToGuiBuilderServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendConnectorPublishSupportedSubInstructionsToFenixGuiBuilderServer
// Connector send Supported SubInstructions to GuiBuilderServer
func (fenixExecutionWorkerObject *MessagesToGuiBuilderServerObjectStruct) SendConnectorPublishSupportedSubInstructionsToFenixGuiBuilderServer(
	supportedSubInstructions *fenixTestCaseBuilderServerGrpcApi.SupportedSubInstructions) (bool, string) {

	fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "6dc6da76-9ddc-479b-8934-8105f1942cf4",
		//"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage,
	}).Debug("Incoming 'SendConnectorPublishSupportedSubInstructionsToFenixGuiBuilderServer'")

	defer fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "14dc5e94-6356-404b-89a1-e7bb426d0cca",
	}).Debug("Outgoing 'SendConnectorPublishSupportedSubInstructionsToFenixGuiBuilderServer'")

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
			"ID": "ef99b115-673d-4823-a48a-788fa1c6d511",
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
	returnMessage, err := tempFenixGuiBuilderServerGrpcClient.ConnectorPublishSupportedSubInstructions(
		ctx,
		supportedSubInstructions)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":    "daa8e257-fbda-44f5-ad10-e461048bdee3",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendConnectorPublishSupportedSubInstructionsToFenixGuiBuilderServer'")

		// Set that a new connection needs to be done next time
		fenixExecutionWorkerObject.connectionToGuiBuilderServerInitiated = false

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":                                 "b940b569-badf-48e8-b0bf-b0666ef7de34",
			"Message from FenixGuiBuilderServer": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendConnectorPublishSupportedSubInstructionsToFenixGuiBuilderServer'")

		return false, returnMessage.Comments
	}

	return returnMessage.AckNack, returnMessage.Comments

}
