package messagesToExecutionServer

import (
	"FenixClientServer/common_config"
	"crypto/tls"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ********************************************************************************************************************

// SetConnectionToFenixTestExecutionServer - Set upp connection and Dial to FenixExecutionServer
func (messagesToExecutionServerObject *messagesToExecutionServerObjectStruct) SetConnectionToFenixTestExecutionServer() {

	var err error
	var opts []grpc.DialOption

	//When running on GCP then use credential otherwise not
	if common_config.ExecutionLocationForFenixTestDataServer == common_config.GCP {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})

		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}
	}

	// Set up connection to FenixTestDataSyncServer
	// When run on GCP, use credentials
	if common_config.ExecutionLocationForFenixTestDataServer == common_config.GCP {
		// Run on GCP
		remoteFenixTestDataSyncServerConnection, err = grpc.Dial(fenixGuiTestCaseBuilderServer_address_to_dial, opts...)
	} else {
		// Run Local
		remoteFenixTestDataSyncServerConnection, err = grpc.Dial(fenixGuiTestCaseBuilderServer_address_to_dial, grpc.WithInsecure())
	}
	if err != nil {
		fenixExecutionWorkerObject.logger.WithFields(logrus.Fields{
			"ID": "50b59b1b-57ce-4c27-aa84-617f0cde3100",
			"fenixGuiTestCaseBuilderServer_address_to_dial": fenixGuiTestCaseBuilderServer_address_to_dial,
			"error message": err,
		}).Error("Did not connect to FenixTestDataSyncServer via gRPC")
		//os.Exit(0)
	} else {
		fenixExecutionWorkerObject.logger.WithFields(logrus.Fields{
			"ID": "0c650bbc-45d0-4029-bd25-4ced9925a059",
			"fenixGuiTestCaseBuilderServer_address_to_dial": fenixGuiTestCaseBuilderServer_address_to_dial,
		}).Info("gRPC connection OK to FenixTestDataSyncServer")

		// Creates a new Clients
		fenixGuiTestCaseBuilderServerClient = fenixExecutionWorkerGrpcApi.NewFenixTestDataGrpcServicesClient(remoteFenixTestDataSyncServerConnection)

	}
}

// ********************************************************************************************************************

// Get the highest FenixProtoFileVersionEnumeration
func (messagesToExecutionServerObject *messagesToExecutionServerObjectStruct) getHighestFenixProtoFileVersion() int32 {

	// Check if there already is a 'highestFenixProtoFileVersion' saved, if so use that one
	if highestFenixProtoFileVersion != -1 {
		return highestFenixProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixExecutionWorkerGrpcApi.CurrentFenixTestDataProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestFenixProtoFileVersion = maxValue

	return highestFenixProtoFileVersion
}

// ********************************************************************************************************************
// Get the highest ClientProtoFileVersionEnumeration
func (messagesToExecutionServerObject *messagesToExecutionServerObjectStruct) getHighestClientProtoFileVersion() int32 {

	// Check if there already is a 'highestclientProtoFileVersion' saved, if so use that one
	if highestClientProtoFileVersion != -1 {
		return highestClientProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixExecutionWorkerGrpcApi.CurrentFenixClientTestDataProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestClientProtoFileVersion = maxValue

	return highestClientProtoFileVersion
}
