package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
)

// InitGrpcServer - Set up and start Backend gRPCServer-server
func (fenixExecutionWorkerGrpcObject *FenixExecutionWorkerGrpcObjectStruct) InitGrpcServer(logger *logrus.Logger) {

	var err error

	// Find first non allocated port from defined start port
	fenixExecutionWorkerGrpcObject.logger.WithFields(logrus.Fields{
		"Id": "054bc0ef-93bb-4b75-8630-74e3823f71da",
	}).Info("Backend Server tries to start")

	fenixExecutionWorkerGrpcObject.logger.WithFields(logrus.Fields{
		"Id": "ca3593b1-466b-4536-be91-5e038de178f4",
		"common_config.FenixExecutionWorkerServerPort: ": common_config.FenixExecutionWorkerServerPort,
	}).Info("Start listening on:")
	lis, err = net.Listen("tcp", ":"+strconv.Itoa(common_config.FenixExecutionWorkerServerPort))

	if err != nil {
		fenixExecutionWorkerGrpcObject.logger.WithFields(logrus.Fields{
			"Id":    "ad7815b3-63e8-4ab1-9d4a-987d9bd94c76",
			"err: ": err,
		}).Error("failed to listen:")
	} else {
		fenixExecutionWorkerGrpcObject.logger.WithFields(logrus.Fields{
			"Id": "ba070b9b-5d57-4c0a-ab4c-a76247a50fd3",
			"common_config.FenixExecutionWorkerServerPort: ": common_config.FenixExecutionWorkerServerPort,
		}).Info("Success in listening on port:")

	}

	// Creates a new RegisterWorkerServer gRPCServer server
	//go func() {
	fenixExecutionWorkerGrpcObject.logger.WithFields(logrus.Fields{
		"Id": "b0ccffb5-4367-464c-a3bc-460cafed16cb",
	}).Info("Starting Backend gRPCServer Server")

	registerFenixExecutionWorkerGrpcServicesServer = grpc.NewServer()
	fenixExecutionWorkerGrpcApi.RegisterFenixExecutionWorkerGrpcServicesServer(registerFenixExecutionWorkerGrpcServicesServer, &fenixExecutionWorkerGrpcServicesServer{logger: logger})

	// Register RouteGuide on the same server.
	reflection.Register(registerFenixExecutionWorkerGrpcServicesServer)

	fenixExecutionWorkerGrpcObject.logger.WithFields(logrus.Fields{
		"Id": "e843ece9-b707-4c60-b1d8-14464305e68f",
		"common_config.FenixExecutionWorkerServerPort: ": common_config.FenixExecutionWorkerServerPort,
	}).Info("registerFenixExecutionWorkerGrpcServicesServer for TestExecution-Worker Server started")
	registerFenixExecutionWorkerGrpcServicesServer.Serve(lis)
	//}()

}

// StopGrpcServer - Stop Backend gRPCServer-server
func (fenixExecutionWorkerGrpcObject *FenixExecutionWorkerGrpcObjectStruct) StopGrpcServer() {

	fenixExecutionWorkerGrpcObject.logger.WithFields(logrus.Fields{}).Info("Gracefully stop for: registerFenixExecutionWorkerGrpcServicesServer")
	registerFenixExecutionWorkerGrpcServicesServer.GracefulStop()

	fenixExecutionWorkerGrpcObject.logger.WithFields(logrus.Fields{
		"common_config.FenixExecutionWorkerServerPort: ": common_config.FenixExecutionWorkerServerPort,
	}).Info("Close net.Listing")
	_ = lis.Close()

}