package messagesToGuiBuilderServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendConnectorPublishSupportedMetaDataToFenixGuiBuilderServer
// Connector send Supported TestCaseMetaData to GuiBuilderServer
func (fenixExecutionWorkerObject *MessagesToGuiBuilderServerObjectStruct) SendConnectorPublishSupportedMetaDataToFenixGuiBuilderServer(
	supportedTestCaseMetaData *fenixTestCaseBuilderServerGrpcApi.
		SupportedTestCaseMetaData) (bool, string) {

	fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "f3192329-b4a2-4ecf-bda3-d16386cadbc1",
		//"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage,
	}).Debug("Incoming 'SendConnectorPublishSupportedMetaDataToFenixGuiBuilderServer'")

	defer fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "0d90b530-a3a3-4892-9499-49edc42a240e",
	}).Debug("Outgoing 'SendConnectorPublishSupportedMetaDataToFenixGuiBuilderServer'")

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
			"ID": "84f571a9-dcbb-4509-a9aa-89551eb27f81",
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
	returnMessage, err := tempFenixGuiBuilderServerGrpcClient.ConnectorPublishSupportedMetaData(
		ctx,
		supportedTestCaseMetaData)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":    "64cc1219-d237-42e3-a2b7-d09aa2e4b7ea",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendConnectorPublishSupportedMetaDataToFenixGuiBuilderServer'")

		// Set that a new connection needs to be done next time
		fenixExecutionWorkerObject.connectionToGuiBuilderServerInitiated = false

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":                                 "9287098f-192c-4907-9783-de4c9538366d",
			"Message from FenixGuiBuilderServer": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendConnectorPublishSupportedMetaDataToFenixGuiBuilderServer'")

		return false, returnMessage.Comments
	}

	return returnMessage.AckNack, returnMessage.Comments

}
