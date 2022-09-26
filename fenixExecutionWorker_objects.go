package main

import (
	"FenixExecutionWorker/testInstructionExecutionEngine"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"net"
)

type fenixExecutionWorkerObjectStruct struct {
	logger                    *logrus.Logger
	gcpAccessToken            *oauth2.Token
	executionEngineChannelRef *testInstructionExecutionEngine.ExecutionEngineChannelType
	executionEngine           *testInstructionExecutionEngine.TestInstructionExecutionEngineStruct
}

// Variable holding everything together
var fenixExecutionWorkerObject *fenixExecutionWorkerObjectStruct

// gRPC variables
var (
	registerFenixExecutionWorkerGrpcServicesServer *grpc.Server
	lis                                            net.Listener
)

// gRPC Server used for register clients Name, Ip and Por and Clients Test Enviroments and Clients Test Commandst
type fenixExecutionWorkerGrpcServicesServer struct {
	fenixExecutionWorkerGrpcApi.UnimplementedFenixExecutionWorkerGrpcServicesServer
}

//TODO FIXA DENNA PATH, HMMM borde köra i DB framöver
// For now hardcoded MerklePath
//var merkleFilterPath string = //"AccountEnvironment/ClientJuristictionCountryCode/MarketSubType/MarketName/" //SecurityType/"
