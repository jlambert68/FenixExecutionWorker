package messagesToExecutionServer

import (
	"FenixExecutionWorker/common_config"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendAreYouAliveToFenixTestDataServer - Send the client's TestDataHeaders to Fenix by calling Fenix's gPRC server
func (fenixExecutionWorkerObject *fenixExecutionWorkerObject_struct) SendAreYouAliveToFenixTestDataServer() (bool, string) {

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string

	// Set up connection to Server
	fenixExecutionWorkerObject.SetConnectionToFenixTestDataSyncServer()

	// Create the message with all test data to be sent to Fenix
	emptyParameter := &fenixExecutionWorkerGrpcApi.EmptyParameter{

		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixTestDataProtoFileVersionEnum(fenixExecutionWorkerObject.getHighestFenixProtoFileVersion()),
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		fenixExecutionWorkerObject.logger.WithFields(logrus.Fields{
			"ID": "c5ba19bd-75ff-4366-818d-745d4d7f1a52",
		}).Error("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixTestDataServer == common_config.GCP {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = fenixExecutionWorkerObject.generateGCPAccessToken(ctx)
		if returnMessageAckNack == false {
			return false, returnMessageString
		}

	}

	returnMessage, err := fenixGuiTestCaseBuilderServerClient.AreYouAlive(ctx, emptyParameter)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.logger.WithFields(logrus.Fields{
			"ID":    "818aaf0b-4112-4be4-97b9-21cc084c7b8b",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixTestDataSyncServer for 'SendAreYouAliveToFenixTestDataServer'")

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		fenixExecutionWorkerObject.logger.WithFields(logrus.Fields{
			"ID": "2ecbc800-2fb6-4e88-858d-a421b61c5529",
			"Message from FenixTestDataSyncServerObject": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixTestDataSyncServer for 'SendAreYouAliveToFenixTestDataServer'")
	}

	return returnMessage.AckNack, returnMessage.Comments

}
