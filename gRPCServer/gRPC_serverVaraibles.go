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

// Variable holding everything together
//var ExecutionWorkerGrpcObject *FenixExecutionWorkerGrpcObjectStruct

// gRPCServer variables
var (
	fenixExecutionWorkerGrpcServer                          *grpc.Server
	registerFenixExecutionWorkerGrpcServicesServer          *grpc.Server
	registerFenixExecutionWorkerConnectorGrpcServicesServer *grpc.Server
	lis                                                     net.Listener
)

// gRPCServer Server used for register clients Name, Ip and Por and Clients Test Enviroments and Clients Test Commandst
type fenixExecutionWorkerGrpcServicesServer struct {
	logger                  *logrus.Logger
	CommandChannelReference *workerEngine.ExecutionEngineChannelType
	fenixExecutionWorkerGrpcApi.UnimplementedFenixExecutionWorkerGrpcServicesServer
}

// gRPCServer Server used for register clients Name, Ip and Por and Clients Test Enviroments and Clients Test Commandst
type fenixExecutionWorkerConnectorGrpcServicesServer struct {
	logger                  *logrus.Logger
	CommandChannelReference *workerEngine.ExecutionEngineChannelType
	fenixExecutionWorkerGrpcApi.UnimplementedFenixExecutionWorkerConnectorGrpcServicesServer
}
