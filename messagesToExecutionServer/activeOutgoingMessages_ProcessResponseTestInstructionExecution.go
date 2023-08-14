package messagesToExecutionServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendProcessResponseTestInstructionExecution - When a TestInstruction has been fully executed the Client use this to inform the results of the execution result to the Server
func (fenixExecutionWorkerObject *MessagesToExecutionServerObjectStruct) SendProcessResponseTestInstructionExecution(processTestInstructionExecutionResponseStatus *fenixExecutionServerGrpcApi.ProcessTestInstructionExecutionResponseStatus) (bool, string) {

	fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "1222d00f-17bc-449a-8de8-a029a9f989f6",
		"processTestInstructionExecutionResponseStatus": processTestInstructionExecutionResponseStatus,
	}).Debug("Incoming 'SendProcessResponseTestInstructionExecution'")

	defer fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"id": "f309fce2-2210-469c-8cf9-8fc023bc8929",
	}).Debug("Outgoing 'SendProcessResponseTestInstructionExecution'")

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string

	// Set up connection to ExecutionServer, if that is not already done
	if fenixExecutionWorkerObject.connectionToExecutionServerInitiated == false {
		err := fenixExecutionWorkerObject.SetConnectionToFenixTestExecutionServer()
		if err != nil {
			return false, err.Error()
		}
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID": "5f02b94f-b07d-4bd7-9607-89cf712824c9",
		}).Debug("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixExecutionServer == common_config.GCP {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = fenixExecutionWorkerObject.generateGCPAccessToken(ctx)
		if returnMessageAckNack == false {
			return false, returnMessageString
		}

	}

	// Creates a new temporary client only to be used for this call
	var tempFenixExecutionServerGrpcClient fenixExecutionServerGrpcApi.FenixExecutionServerGrpcServicesClient
	tempFenixExecutionServerGrpcClient = fenixExecutionServerGrpcApi.NewFenixExecutionServerGrpcServicesClient(remoteFenixExecutionServerConnection)

	// Do gRPC-call
	returnMessage, err := tempFenixExecutionServerGrpcClient.ProcessResponseTestInstructionExecution(ctx, processTestInstructionExecutionResponseStatus)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":    "33134d62-5666-40cd-9894-b8a30a603277",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixExecutionServer for 'SendProcessResponseTestInstructionExecution'")

		// Set that a new connection needs to be done next time
		fenixExecutionWorkerObject.connectionToExecutionServerInitiated = false

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":                                  "7d1d4021-5e0f-4be7-8a7e-0554833cf0b9",
			"Message from Fenix Execution Server": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixExecutionServer for 'SendProcessResponseTestInstructionExecution'")

		return false, returnMessage.Comments
	}

	return returnMessage.AckNack, returnMessage.Comments

}
