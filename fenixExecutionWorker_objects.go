package main

import (
	"FenixExecutionWorker/gRPCServer"
	"github.com/sirupsen/logrus"
)

type fenixExecutionWorkerObjectStruct struct {
	logger     *logrus.Logger
	GrpcServer *gRPCServer.FenixExecutionWorkerGrpcObjectStruct
}

// Variable holding everything together
var FenixExecutionWorkerObject *fenixExecutionWorkerObjectStruct
