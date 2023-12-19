package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/gcp"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// SendConnectorInformsItIsAlive - *********************************************************************
// Connector informs Worker that Connector is up and running and ready to receive work
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorInformsItIsAlive(
	ctx context.Context,
	connectorIsReadyMessage *fenixExecutionWorkerGrpcApi.ConnectorIsReadyMessage) (
	connectorIsReadyResponseMessage *fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage, err error) {

	s.logger.WithFields(logrus.Fields{
		"id":                      "b9612e9a-0b94-4113-ac65-eae9a9732ec7",
		"connectorIsReadyMessage": connectorIsReadyMessage,
	}).Debug("Incoming 'gRPCServer - ConnectorInformsItIsAlive'")

	s.logger.WithFields(logrus.Fields{
		"id": "87e4f4b4-ddd7-4d11-8369-4cf1ef93c131",
	}).Debug("Outgoing 'gRPCServer - ConnectorInformsItIsAlive'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	var ackNackResponse *fenixExecutionWorkerGrpcApi.AckNackResponse
	ackNackResponse = common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, connectorIsReadyMessage.
		ClientSystemIdentification.ProtoFileVersionUsedByClient)
	if ackNackResponse != nil {

		connectorIsReadyResponseMessage = &fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage{
			AckNackResponse:          ackNackResponse,
			PubSubAuthorizationToken: "",
		}

		return connectorIsReadyResponseMessage, nil

	}

	// Get latest authorization token or create a new one to access PubSub for Connector
	//var appendedCtx context.Context
	var returnAckNack bool
	var returnMessage string

	_, returnAckNack, returnMessage = gcp.Gcp.GenerateGCPAccessToken(context.Background(), gcp.GenerateTokenForPubSub)
	if returnAckNack == false {
		ackNackResponse = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     returnMessage,
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		connectorIsReadyResponseMessage = &fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage{
			AckNackResponse:          ackNackResponse,
			PubSubAuthorizationToken: "",
		}
	}

	// Extract Token to be used by Connector to do a PubSub-Subscription
	/*var pubSubTokenAsAny any
	var pubSubTokenAsString string

	pubSubTokenAsAny = appendedCtx.Value("authorization")

	pubSubTokenAsString, ok := pubSubTokenAsAny.(string)
	if ok == false {
		ackNackResponse = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     returnMessage,
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		connectorIsReadyResponseMessage = &fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage{
			AckNackResponse:          ackNackResponse,
			PubSubAuthorizationToken: "",
		}
	}

	*/

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
			AckNackResponse:          ackNackResponse,
			PubSubAuthorizationToken: "Not used for now",
		}

	} else {
		// Running in GCP
		connectorIsReadyResponseMessage = &fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage{
			AckNackResponse:          ackNackResponse,
			PubSubAuthorizationToken: "Not used for now",
		}

	}

	return connectorIsReadyResponseMessage, nil
}
