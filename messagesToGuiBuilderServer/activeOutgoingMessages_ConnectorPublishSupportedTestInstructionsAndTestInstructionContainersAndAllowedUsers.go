package messagesToGuiBuilderServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersToFenixGuiBuilderServer
// Connector send available TestInstructions, TestInstructionContainers and Allowed Users and this function forwards them to GuiBuilderServer
func (fenixExecutionWorkerObject *MessagesToGuiBuilderServerObjectStruct) SendPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersToFenixGuiBuilderServer(
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage *fenixTestCaseBuilderServerGrpcApi.
		SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage) (bool, string) {

	fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "7a464b61-200e-418a-bf6a-a63d1f6608ab",
		"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage,
	}).Debug("Incoming 'SendPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersToFenixGuiBuilderServer'")

	defer fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "afad0715-be5f-4bb9-931e-dd760a07ca03",
	}).Debug("Outgoing 'SendPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersToFenixGuiBuilderServer'")

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string

	// Set up connection to ExecutionServer, if that is not already done
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
			"ID": "43856814-8dfe-4ba3-95da-3bac2b35f42e",
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
	returnMessage, err := tempFenixGuiBuilderServerGrpcClient.PublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
		ctx,
		supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":    "ebe601e0-14b9-42c5-8f8f-960acec80433",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersToFenixGuiBuilderServer'")

		// Set that a new connection needs to be done next time
		fenixExecutionWorkerObject.connectionToGuiBuilderServerInitiated = false

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":                                  "24ebaa96-fbeb-4c51-bcd0-87179c887b15",
			"Message from Fenix Execution Server": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersToFenixGuiBuilderServer'")

		return false, returnMessage.Comments
	}

	return returnMessage.AckNack, returnMessage.Comments

}
