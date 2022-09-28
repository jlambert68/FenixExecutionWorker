package gRPCServer

import (
	"FenixExecutionWorker/testInstructionExecutionEngine"
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
	registerFenixExecutionWorkerGrpcServicesServer *grpc.Server
	lis                                            net.Listener
)

// gRPCServer Server used for register clients Name, Ip and Por and Clients Test Enviroments and Clients Test Commandst
type fenixExecutionWorkerGrpcServicesServer struct {
	logger                  *logrus.Logger
	CommandChannelReference *testInstructionExecutionEngine.ExecutionEngineChannelType
	fenixExecutionWorkerGrpcApi.UnimplementedFenixExecutionWorkerGrpcServicesServer
}
