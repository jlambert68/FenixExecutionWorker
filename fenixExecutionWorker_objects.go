package main

import (
	"FenixExecutionWorker/gRPCServer"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type fenixExecutionWorkerObjectStruct struct {
	logger         *logrus.Logger
	gcpAccessToken *oauth2.Token
	GrpcServer     *gRPCServer.FenixExecutionWorkerGrpcObjectStruct
}

// Variable holding everything together
var FenixExecutionWorkerObject *fenixExecutionWorkerObjectStruct
