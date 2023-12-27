package messagesToGuiBuilderServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendGetMessageToSignToProveCallerIdentity
// Worker ask BuilderServer for a message to sign and use the signature to prove identity when sending 'SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage'
func (fenixExecutionWorkerObject *MessagesToGuiBuilderServerObjectStruct) SendGetMessageToSignToProveCallerIdentity() (
	returnMessageAckNack bool,
	returnMessageString string,
	messageToSignToProveIdentity string) {

	fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "d702fb24-7ff1-42d9-877f-ac2b1717b262",
	}).Debug("Incoming 'SendGetMessageToSignToProveCallerIdentity'")

	defer fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "e49cae85-cc36-4fc8-bef2-3e78ccf15bb0",
	}).Debug("Outgoing 'SendGetMessageToSignToProveCallerIdentity'")

	var ctx context.Context

	// Set up connection to BuilderServer, if that is not already done
	if fenixExecutionWorkerObject.connectionToGuiBuilderServerInitiated == false {
		err := fenixExecutionWorkerObject.SetConnectionToFenixGuiBuilderServer()
		if err != nil {
			return false, err.Error(), ""
		}
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID": "70a1b92b-d520-4a5c-b36c-c8c2a12258f7",
		}).Debug("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixGuiBuilderServer == common_config.GCP {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = fenixExecutionWorkerObject.generateGCPAccessToken(ctx)
		if returnMessageAckNack == false {
			return false, returnMessageString, ""
		}

	}

	// Creates a new temporary client only to be used for this call
	var tempFenixGuiBuilderServerGrpcClient fenixTestCaseBuilderServerGrpcApi.FenixTestCaseBuilderServerGrpcWorkerServicesClient
	tempFenixGuiBuilderServerGrpcClient = fenixTestCaseBuilderServerGrpcApi.NewFenixTestCaseBuilderServerGrpcWorkerServicesClient(remoteFenixGuiBuilderServerConnection)

	// Create empty message to send
	var emptyParameter *fenixTestCaseBuilderServerGrpcApi.EmptyParameter
	emptyParameter = &fenixTestCaseBuilderServerGrpcApi.EmptyParameter{}

	// Do gRPC-call
	returnMessage, err := tempFenixGuiBuilderServerGrpcClient.GetMessageToSignToProveCallerIdentity(
		ctx, emptyParameter)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":    "a04b0d74-2158-44f4-b8b6-bac531eb9a75",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendGetMessageToSignToProveCallerIdentity'")

		// Set that a new connection needs to be done next time
		fenixExecutionWorkerObject.connectionToGuiBuilderServerInitiated = false

		return false, err.Error(), ""

	} else if returnMessage.AckNack.AckNack == false {
		// FenixBuilderServer couldn't handle gPRC call
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":                                "24ebaa96-fbeb-4c51-bcd0-87179c887b15",
			"Message from Fenix Builder Server": returnMessage.GetAckNack().GetComments(),
		}).Error("Problem to do gRPC-call to FenixGuiBuilderServer for 'SendGetMessageToSignToProveCallerIdentity'")

		return false, returnMessage.GetAckNack().GetComments(), ""
	}

	return returnMessage.GetAckNack().GetAckNack(), returnMessage.GetAckNack().GetComments(), returnMessage.GetMessageToSign()

}
