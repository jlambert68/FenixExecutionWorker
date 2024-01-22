package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
)

// SendConnectorInformsItIsAlive - *********************************************************************
// Connector informs Worker that Connector is up and running and ready to receive work
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorInformsItIsAlive(
	ctx context.Context,
	connectorIsReadyMessage *fenixExecutionWorkerGrpcApi.ConnectorIsReadyMessage) (
	connectorIsReadyResponseMessage *fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage, err error) {

	/*
		s.logger.WithFields(logrus.Fields{
			"id":                      "b9612e9a-0b94-4113-ac65-eae9a9732ec7",
			"connectorIsReadyMessage": connectorIsReadyMessage,
		}).Debug("Incoming 'gRPCServer - ConnectorInformsItIsAlive'")

		s.logger.WithFields(logrus.Fields{
			"id": "0e0f68ac-fde0-4c93-b47d-73bfb44a68bd",
		}).Debug("Outgoing 'gRPCServer - ConnectorInformsItIsAlive'")
	*/

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	var ackNackResponse *fenixExecutionWorkerGrpcApi.AckNackResponse
	ackNackResponse = common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, connectorIsReadyMessage.
		ClientSystemIdentification.ProtoFileVersionUsedByClient)
	if ackNackResponse != nil {

		connectorIsReadyResponseMessage = &fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage{
			AckNackResponse: ackNackResponse,
		}

		return connectorIsReadyResponseMessage, nil

	}

	// Create response message to Connector
	ackNackResponse = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	if common_config.ExecutionLocationForWorker == common_config.LocalhostNoDocker {
		// Running Locally
		connectorIsReadyResponseMessage = &fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage{
			AckNackResponse: ackNackResponse,
		}

	} else {
		// Running in GCP
		connectorIsReadyResponseMessage = &fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage{
			AckNackResponse: ackNackResponse,
		}

	}

	return connectorIsReadyResponseMessage, nil
}
