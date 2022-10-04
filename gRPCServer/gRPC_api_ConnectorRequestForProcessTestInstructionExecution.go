package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/workerEngine"
	"errors"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ConnectorRequestForProcessTestInstructionExecution
// Used to send TestInstructions for Execution to Connector. Worker Stream tasks as response and it is done this way due to it is impossible to call Connector on SEB network
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorRequestForProcessTestInstructionExecution(emptyParameter *fenixExecutionWorkerGrpcApi.EmptyParameter, streamServer fenixExecutionWorkerGrpcApi.FenixExecutionWorkerConnectorGrpcServices_ConnectorRequestForProcessTestInstructionExecutionServer) (err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "d986194e-ec8c-4198-8160-bd7eb9838aca",
	}).Debug("Incoming 'gRPCServer - ConnectorRequestForProcessTestInstructionExecution'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "1b9fb882-f3aa-4ffa-b575-910569aec6c4",
	}).Debug("Outgoing 'gRPCServer - ConnectorRequestForProcessTestInstructionExecution'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(emptyParameter.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return errors.New(returnMessage.Comments)
	}

	// Local channel to decide when Server stopped sending
	done := make(chan bool)

	go func() {

		// Wait for incoming TestInstructionExecution from Execution Server
		executionForwardChannelMessage := <-executionForwardChannel

		testInstructionExecution := executionForwardChannelMessage.processTestInstructionExecutionReveredRequest

		err = streamServer.Send(testInstructionExecution)
		if err != nil {

			s.logger.WithFields(logrus.Fields{
				"id":                       "70ab1dcb-0be3-49b6-b49a-694bab529ed4",
				"err":                      err,
				"testInstructionExecution": testInstructionExecution,
			}).Error("Got some problem when doing reversed streaming of TestInstructionExecution to Connector. Stopping Reversed Streaming")

			// Create response message to be sent on response channel
			executionResponseChanneMessage := executionResponseChannelStruct{
				testInstructionExecutionIsSentToConnector: false,
				err: err,
			}

			// Send message back over response channel that message was failed to be sent to Connector
			*executionForwardChannelMessage.executionResponseChannelReference <- executionResponseChanneMessage

			// Have the gRPC-call be continued
			done <- true //close(done)

			return
		}

		// Send message back over response channel that message was sent to Connector

		s.logger.WithFields(logrus.Fields{
			"id":                       "6f5e6dc7-cef5-4008-a4ea-406be80ded4c",
			"testInstructionExecution": testInstructionExecution,
		}).Debug("Success in reversed streaming TestInstructionExecution to Connector")
	}()

	// Server stopped so wait for new connection
	<-done

	// Send Message on CommandChannel to be able to send Result back to Fenix Execution Server
	channelCommand := workerEngine.ChannelCommandStruct{
		ChannelCommand: workerEngine.ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServer,
		ReportCompleteTestInstructionExecutionResultParameter: workerEngine.ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServerStruct{
			FinalTestInstructionExecutionResultMessage: finalTestInstructionExecutionResultMessage},
	}

	*s.CommandChannelReference <- channelCommand

	// Generate response
	ackNackResponse = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return ackNackResponse, nil

}
