package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"errors"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// ConnectorRequestForProcessTestInstructionExecution
// Used to send TestInstructions for Execution to Connector. Worker Stream tasks as response, and it is done this way due to it is impossible to call Connector on SEB network
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

	// Recreate channel for incoming TestInstructionExecution from Execution Server
	executionForwardChannel = make(chan executionForwardChannelStruct)

	// Local channel to decide when Server stopped sending
	done := make(chan bool)

	go func() {

		for {
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

				// Only send back response over response channel if it wasn't a 'keep alive' message
				if testInstructionExecution.TestInstruction.TestInstructionName != "KeepAlive" {

					// Create response message to be sent on response channel
					executionResponseChannelMessage := executionResponseChannelStruct{
						testInstructionExecutionIsSentToConnector: false,
						err: err,
					}

					// Send message back over response channel that message was failed to be sent to Connector
					*executionForwardChannelMessage.executionResponseChannelReference <- executionResponseChannelMessage
				}

				// Have the gRPC-call be continued, end stream server
				done <- true //close(done)

				return

			}

			// Check if message only was a keep alive message to Connector
			if executionForwardChannelMessage.isKeepAliveMessage == false {

				// Is a standard TestInstructionExecution that was sent to Connector
				s.logger.WithFields(logrus.Fields{
					"id":                       "6f5e6dc7-cef5-4008-a4ea-406be80ded4c",
					"testInstructionExecution": testInstructionExecution,
				}).Debug("Success in reversed streaming TestInstructionExecution to Connector")

			} else {

				// Is a keep alive message that was sent to Connector
				s.logger.WithFields(logrus.Fields{
					"id":                       "c1d5a756-b7fa-48ae-953e-59dedd0671f4",
					"testInstructionExecution": testInstructionExecution,
				}).Debug("Success in reversed streaming TestInstructionExecution to Connector")
			}
		}

	}()

	// Feed 'executionForwardChannel' with messages every 15 seconds to keep connection open to Connector
	go func() {

		// Create keep alive message
		ProcessTestInstructionExecutionReveredRequest_TestInstructionExecutionMessage := fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest_TestInstructionExecutionMessage{
			TestInstructionExecutionUuid: "KeepAlive",
			TestInstructionUuid:          "KeepAlive",
			TestInstructionName:          "KeepAlive",
			MajorVersionNumber:           0,
			MinorVersionNumber:           0,
			TestInstructionAttributes:    nil,
		}
		processTestInstructionExecutionReveredRequest := fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest{
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
			TestInstruction:              &ProcessTestInstructionExecutionReveredRequest_TestInstructionExecutionMessage,
			TestData:                     nil,
		}
		keepAliveMessageToConnector := executionForwardChannelStruct{
			processTestInstructionExecutionReveredRequest: &processTestInstructionExecutionReveredRequest,
			executionResponseChannelReference:             nil,
			isKeepAliveMessage:                            true,
		}

		for {

			// Send Keep Alive message on channel to be sent to Connector
			executionForwardChannel <- keepAliveMessageToConnector

			// Sleep for 15 seconds before continue
			time.Sleep(time.Second * 15)

		}
	}()

	// Server stopped so wait for new connection
	<-done

	return err

}