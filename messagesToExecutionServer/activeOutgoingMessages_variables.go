package messagesToExecutionServer

import (
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type MessagesToExecutionServerObjectStruct struct {
	Logger         *logrus.Logger
	gcpAccessToken *oauth2.Token
}

// Variables used for contacting Fenix Execution Server
var (
	remoteFenixExecutionServerConnection *grpc.ClientConn
	//FenixExecutionServerAddressToDial    string
	fenixExecutionServerGrpcClient fenixExecutionServerGrpcApi.FenixExecutionServerGrpcServicesClient
)
