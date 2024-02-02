package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

// AreYouAlive - *********************************************************************
// Anyone can check if Fenix Execution Worker server is alive with this service, ushould be used to check serves for Connector
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorAreYouAlive(ctx context.Context, emptyParameter *fenixExecutionWorkerGrpcApi.EmptyParameter) (*fenixExecutionWorkerGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id": "5c2d4e0c-904a-41d8-81bc-3123641aa6db",
	}).Debug("Incoming 'gRPCServer - ConnectorAreYouAlive'")

	s.logger.WithFields(logrus.Fields{
		"id": "b9003ecf-b686-429b-b603-261f78e9c787",
	}).Debug("Outgoing 'gRPCServer - ConnectorAreYouAlive'")

	ackNackResponse := &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     fmt.Sprintf("I'am alive and the time is %s", time.Now().String()),
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return ackNackResponse, nil
}
