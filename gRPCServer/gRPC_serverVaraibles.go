package gRPCServer

import (
	"FenixExecutionWorker/workerEngine"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type FenixExecutionWorkerGrpcObjectStruct struct {
	logger                    *logrus.Logger
	ExecutionWorkerGrpcObject *FenixExecutionWorkerGrpcObjectStruct
}

// gRPCServer variables
var (
	fenixExecutionWorkerGrpcServer *grpc.Server
	lis                            net.Listener
)

// gRPCServer Server used from Execution Server
type fenixExecutionWorkerGrpcServicesServer struct {
	logger                  *logrus.Logger
	CommandChannelReference *workerEngine.ExecutionEngineChannelType
	fenixExecutionWorkerGrpcApi.UnimplementedFenixExecutionWorkerGrpcServicesServer
}

// gRPCServer Server used from Connector
type fenixExecutionWorkerConnectorGrpcServicesServer struct {
	logger                  *logrus.Logger
	CommandChannelReference *workerEngine.ExecutionEngineChannelType
	fenixExecutionWorkerGrpcApi.UnimplementedFenixExecutionWorkerConnectorGrpcServicesServer
}

// Channel used for forwarding TestInstructionExecutions to stream-server which then forwards it to the Connector
var executionForwardChannel executionForwardChannelType

type executionForwardChannelType chan executionForwardChannelStruct

type executionForwardChannelStruct struct {
	processTestInstructionExecutionReveredRequest *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest
	executionResponseChannelReference             *executionResponseChannelType
}

// Channel used for response from Stream server (from Worker to Connector) that message has been sent
type executionResponseChannelType chan executionResponseChannelStruct

type executionResponseChannelStruct struct {
	testInstructionExecutionIsSentToConnector bool
	err                                       error
}
