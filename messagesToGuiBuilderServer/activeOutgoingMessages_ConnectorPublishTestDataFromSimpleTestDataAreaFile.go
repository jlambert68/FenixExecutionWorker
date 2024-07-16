package messagesToGuiBuilderServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendConnectorPublishTestDataFromSimpleTestDataAreaFileMessageToBuilderServer
// Connector send TestData from 'simple' file to GuiBuilderServer
func (fenixExecutionWorkerObject *MessagesToGuiBuilderServerObjectStruct) SendConnectorPublishTestDataFromSimpleTestDataAreaFileMessageToBuilderServer(
	testDataFromSimpleTestDataAreaFileMessage *fenixTestCaseBuilderServerGrpcApi.
		TestDataFromSimpleTestDataAreaFileMessage) (bool, string) {

	fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "81402ded-b99a-4600-98b4-66239736693c",
	}).Debug("Incoming 'SendConnectorPublishTestDataFromSimpleTestDataAreaFileMessageToBuilderServer'")

	defer fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "3c02d124-6f5e-402c-aa5c-6bce0d2566bf",
	}).Debug("Outgoing 'SendConnectorPublishTestDataFromSimpleTestDataAreaFileMessageToBuilderServer'")

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
			"ID": "91e0fc6e-3396-4a15-81e2-9c623a56bdf0",
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
	returnMessage, err := tempFenixGuiBuilderServerGrpcClient.ConnectorPublishTestDataFromSimpleTestDataAreaFile(
		ctx,
		testDataFromSimpleTestDataAreaFileMessage)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":    "98909b2c-1981-4135-a0b9-f56bdf4461ff",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendConnectorPublishTestDataFromSimpleTestDataAreaFileMessageToBuilderServer'")

		// Set that a new connection needs to be done next time
		fenixExecutionWorkerObject.connectionToGuiBuilderServerInitiated = false

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":                                 "c1b838d2-159a-440b-b8cb-391800691ef8",
			"Message from FenixGuiBuilderServer": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendConnectorPublishTestDataFromSimpleTestDataAreaFileMessageToBuilderServer'")

		return false, returnMessage.Comments
	}

	return returnMessage.AckNack, returnMessage.Comments

}
