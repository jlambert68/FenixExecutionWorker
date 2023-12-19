package messagesToGuiBuilderServer

import (
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type MessagesToGuiBuilderServerObjectStruct struct {
	Logger                                *logrus.Logger
	gcpAccessToken                        *oauth2.Token
	connectionToGuiBuilderServerInitiated bool
}

// Variables used for contacting Fenix Execution Server
var (
	remoteFenixGuiBuilderServerConnection              *grpc.ClientConn
	fenixTestCaseBuilderServerGrpcWorkerServicesClient fenixTestCaseBuilderServerGrpcApi.FenixTestCaseBuilderServerGrpcWorkerServicesClient
)
