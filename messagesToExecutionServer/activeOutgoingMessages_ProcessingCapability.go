package messagesToExecutionServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendReportProcessingCapabilityToFenixExecutionServer - Worker send the execution capabilities regrading parallell executions
func (fenixExecutionWorkerObject *MessagesToExecutionServerObjectStruct) SendReportProcessingCapabilityToFenixExecutionServer(processingCapabilityMessage *fenixExecutionServerGrpcApi.ProcessingCapabilityMessage) (bool, string) {

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string

	// Set up connection to Server
	err := fenixExecutionWorkerObject.SetConnectionToFenixTestExecutionServer()
	if err != nil {
		return false, err.Error()
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID": "f3aa9000-c175-407f-bdd8-96624c087a39",
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

	returnMessage, err := fenixExecutionServerGrpcClient.ReportProcessingCapability(ctx, processingCapabilityMessage)

	// Shouldn't happen
	if err != nil {
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":    "864d7750-d387-49e7-8eed-286650e52036",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixExecutionServer for 'SendReportProcessingCapabilityToFenixExecutionServer'")

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID":                                  "d8abb6a3-d152-42ed-9e99-051e90d59c91",
			"Message from Fenix Execution Server": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixExecutionServer for 'SendReportProcessingCapabilityToFenixExecutionServer'")

		return false, err.Error()
	}

	return returnMessage.AckNack, returnMessage.Comments

}
