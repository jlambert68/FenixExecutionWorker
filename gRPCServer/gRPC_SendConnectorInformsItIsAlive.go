package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// SendConnectorInformsItIsAlive - *********************************************************************
// Connector informs Worker that Connector is up and running and ready to receive work
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorInformsItIsAlive(
	ctx context.Context,
	connectorIsReadyMessage *fenixExecutionWorkerGrpcApi.ConnectorIsReadyMessage) (
	*fenixExecutionWorkerGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id":                      "b9612e9a-0b94-4113-ac65-eae9a9732ec7",
		"connectorIsReadyMessage": connectorIsReadyMessage,
	}).Debug("Incoming 'gRPCServer - SendConnectorInformsItIsAlive'")

	s.logger.WithFields(logrus.Fields{
		"id": "87e4f4b4-ddd7-4d11-8369-4cf1ef93c131",
	}).Debug("Outgoing 'gRPCServer - SendConnectorInformsItIsAlive'")

	ackNackResponse := &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return ackNackResponse, nil
}
