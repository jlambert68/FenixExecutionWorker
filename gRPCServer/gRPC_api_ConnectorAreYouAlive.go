package gRPCServer

import (
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// AreYouAlive - *********************************************************************
//Anyone can check if Fenix Execution Worker server is alive with this service, ushould be used to check serves for Connector
func (s *fenixExecutionWorkerGrpcServicesServer) ConnectorAreYouAlive(ctx context.Context, emptyParameter *fenixExecutionWorkerGrpcApi.EmptyParameter) (*fenixExecutionWorkerGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id": "5c2d4e0c-904a-41d8-81bc-3123641aa6db",
	}).Debug("Incoming 'gRPCServer - ConnectorAreYouAlive'")

	s.logger.WithFields(logrus.Fields{
		"id": "b9003ecf-b686-429b-b603-261f78e9c787",
	}).Debug("Outgoing 'gRPCServer - ConnectorAreYouAlive'")

	return &fenixExecutionWorkerGrpcApi.AckNackResponse{AckNack: true, Comments: "I'am alive."}, nil
}
