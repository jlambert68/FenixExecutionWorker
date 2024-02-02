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
// Anyone can check if Fenix Execution Worker server is alive with this service
func (s *fenixExecutionWorkerGrpcServicesServer) AreYouAlive(ctx context.Context, emptyParameter *fenixExecutionWorkerGrpcApi.EmptyParameter) (*fenixExecutionWorkerGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id": "1ff67695-9a8b-4821-811d-0ab8d33c4d8b",
	}).Debug("Incoming 'gRPCServer - AreYouAlive'")

	s.logger.WithFields(logrus.Fields{
		"id": "9c7f0c3d-7e9f-4c91-934e-8d7a22926d84",
	}).Debug("Outgoing 'gRPCServer - AreYouAlive'")

	ackNackResponse := &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     fmt.Sprintf("I'am alive and the time is %s", time.Now().String()),
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return ackNackResponse, nil
}
