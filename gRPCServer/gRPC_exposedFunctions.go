package gRPCServer

import "github.com/sirupsen/logrus"

// InitiateLogger - Initiate local logger object
func (fenixExecutionWorkerGrpcObject *FenixExecutionWorkerGrpcObjectStruct) InitiateLogger(logger *logrus.Logger) {

	fenixExecutionWorkerGrpcObject.logger = logger
}

// InitiateLocalObject - Initiate local 'ExecutionWorkerGrpcObject'
func (fenixExecutionWorkerGrpcObject *FenixExecutionWorkerGrpcObjectStruct) InitiateLocalObject(inFenixExecutionWorkerGrpcObject *FenixExecutionWorkerGrpcObjectStruct) {

	fenixExecutionWorkerGrpcObject.ExecutionWorkerGrpcObject = inFenixExecutionWorkerGrpcObject
}
