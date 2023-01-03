package messagesToExecutionServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendAreYouAliveToFenixExecutionServer - Ask Fenix Execution Server to check if it's up and running
func (fenixExecutionWorkerObject *MessagesToExecutionServerObjectStruct) SendAreYouAliveToFenixExecutionServer() (bool, string) {

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string

	// Set up connection to Server
	err := fenixExecutionWorkerObject.SetConnectionToFenixTestExecutionServer()
	if err != nil {
		return false, err.Error()
	}

	// Create the message with all test data to be sent to Fenix
	emptyParameter := &fenixExecutionServerGrpcApi.EmptyParameter{

		ProtoFileVersionUsedByClient: fenixExecutionServerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(common_config.GetHighestFenixExecutionServerProtoFileVersion()),
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID": "2c4fb00e-bcc0-4232-963e-df0c4d2f010d",
		}).Error("Running Defer Cancel function")
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

	returnMessage, err := fenixExecutionServerGrpcClient.AreYouAlive(ctx, emptyParameter)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":    "818aaf0b-4112-4be4-97b9-21cc084c7b8b",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixExecutionServer for 'SendAreYouAliveToFenixExecutionServer'")

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":                                  "2ecbc800-2fb6-4e88-858d-a421b61c5529",
			"Message from Fenix Execution Server": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixExecutionServer for 'SendAreYouAliveToFenixExecutionServer'")

		return false, err.Error()
	}

	return returnMessage.AckNack, returnMessage.Comments

}
