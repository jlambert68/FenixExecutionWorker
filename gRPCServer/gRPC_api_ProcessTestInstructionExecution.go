package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// ProcessTestInstructionExecution
// Fenix Execution Server send a request to Execution Worker to initiate a execution of a TestInstruction
func (s *fenixExecutionWorkerGrpcServicesServer) ProcessTestInstructionExecution(ctx context.Context, processTestInstructionExecutionRequest *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest) (processTestInstructionExecutionResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id":                                     "37bc2356-33a2-4e2c-9420-122df581d757",
		"processTestInstructionExecutionRequest": processTestInstructionExecutionRequest,
	}).Debug("Incoming 'gRPCServer - ProcessTestInstructionExecution'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "f3fd3e50-5770-48ad-8524-85f34d28545e",
	}).Debug("Outgoing 'gRPCServer - ProcessTestInstructionExecution'")

	// Calling system
	userId := "Execution Server"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(processTestInstructionExecutionRequest.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		// Exiting
		processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
			AckNackResponse:                returnMessage,
			TestInstructionExecutionUuid:   "",
			ExpectedExecutionDuration:      nil,
			TestInstructionCanBeReExecuted: false,
		}

		s.logger.WithFields(logrus.Fields{
			"id":            "de3c6e57-f1d1-4a51-9f31-56d7ded167f9",
			"returnMessage": returnMessage,
		}).Debug("Client using wrong proto-file version")

		return processTestInstructionExecutionResponse, nil
	}

	// slice with sleep time, in milliseconds, between each attempt to do gRPC-call to Worker
	var sleepTimeBetweenConnectorIsConnectedCheckAttempts []int
	sleepTimeBetweenConnectorIsConnectedCheckAttempts = []int{100, 200, 300, 300, 500, 500, 1000, 1000, 1000, 1000} // Total: 5.9 seconds

	// Do multiple attempts to do gRPC-call to Execution Worker, when it fails
	//var numberOfConnectorIsConnectedCheckAttempts int
	var connectorIsConnectedCheckAttemptCounter int
	//numberOfConnectorIsConnectedCheckAttempts = len(sleepTimeBetweenConnectorIsConnectedCheckAttempts)
	connectorIsConnectedCheckAttemptCounter = 0

	for {

		// Add to counter for how many gRPC-call-attempts to Worker that have been done
		connectorIsConnectedCheckAttemptCounter = connectorIsConnectedCheckAttemptCounter + 1

		s.logger.WithFields(logrus.Fields{
			"id": "23fb6d78-d4ce-4b39-8091-e9bfeb59f354",
			"connectorIsConnectedCheckAttemptCounter": connectorIsConnectedCheckAttemptCounter,
		}).Debug("Check if there is an active connection from a Connector")

		// Wait up to 30 seconds, when there has been at least one connection from the Connector
		if connectorHasConnectedAtLeastOnce == true {

			for {

				// Wait if there is no connection to Connector
				if connectorHasConnected == false {
					s.logger.WithFields(logrus.Fields{
						"id": "1b949202-5ba1-49c5-b158-a0ecc33314e6",
						"connectorIsConnectedCheckAttemptCounter": connectorIsConnectedCheckAttemptCounter,
					}).Debug("Connector has connected at least one time so do multiple checks if there is an active connection from a Connector")

					// Wait before checking connection again
					time.Sleep(time.Second * 5)

					s.logger.WithFields(logrus.Fields{
						"id": "2cff6ae7-438e-4585-8601-3899b9610fe8",
						"connectorIsConnectedCheckAttemptCounter": connectorIsConnectedCheckAttemptCounter,
					}).Debug("Slept 5 seconds")

					// Max wait time is 30 seconds
					if time.Now().After(connectorConnectionTime.Add(30*time.Second)) == true {
						s.logger.WithFields(logrus.Fields{
							"id": "2d07fa75-672e-4f3d-b61d-5cc1859c1c20",
							"connectorIsConnectedCheckAttemptCounter": connectorIsConnectedCheckAttemptCounter,
						}).Debug("No connector after 30 seconds, so exiting")
						break
					}

				} else {

					break

				}
			}
		}

		// If there isn't an active connection to the Connector then  report that back
		if connectorHasConnected == false {

			// Only return the that no Connector has connected after last attempt
			if true { // Not used for now.... connectorIsConnectedCheckAttemptCounter >= numberOfConnectorIsConnectedCheckAttempts {

				// Create Error Code
				var errorCodes []fenixExecutionWorkerGrpcApi.ErrorCodesEnum
				var errorCode fenixExecutionWorkerGrpcApi.ErrorCodesEnum
				errorCode = fenixExecutionWorkerGrpcApi.ErrorCodesEnum_ERROR_CONNECTOR_NOT_RESPONDING
				errorCodes = append(errorCodes, errorCode)

				// Generate response
				processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
					AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
						AckNack:                      false,
						Comments:                     fmt.Sprintf("Message couldn't be sent to Connector, due to no active Connector was found"),
						ErrorCodes:                   errorCodes,
						ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
					},
					TestInstructionExecutionUuid:   processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid,
					ExpectedExecutionDuration:      nil,
					TestInstructionCanBeReExecuted: false,
				}

				s.logger.WithFields(logrus.Fields{
					"id": "47763d5d-6c20-4515-8109-da9069756c51",
				}).Error("No active connection from a Connector")

				// Return Response to Execution Server
				return processTestInstructionExecutionResponse, nil
			}

			// Sleep for some time before checking if Connector has Connected to Worker
			time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenConnectorIsConnectedCheckAttempts[connectorIsConnectedCheckAttemptCounter-1]))

		} else {
			// Connector has connected to Worker
			break
		}
	}

	s.logger.WithFields(logrus.Fields{
		"id":                                     "0909cb27-ab05-446b-9fe3-c36b05a6137b",
		"processTestInstructionExecutionRequest": processTestInstructionExecutionRequest,
	}).Debug("Received 'processTestInstructionExecutionRequest' from Execution Server")

	//  Check that TestInstructionExecutionUuid isn't already is in Map
	//_, existsInMap :=  processTestInstructionExecutionReversedResponseChannelMap[processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid]
	_, existsInMap := loadFromTestInstructionExecutionReversedResponseChannelMap(processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid)

	// Shouldn't exist in map
	if existsInMap == true {
		s.logger.WithFields(logrus.Fields{
			"id":                                     "df3ddde1-f55d-4d47-86bf-88626a6bb3ea",
			"processTestInstructionExecutionRequest": processTestInstructionExecutionRequest,
		}).Error("TestInstructionExecutionUuid already exists i 'processTestInstructionExecutionReversedResponseChannelMap'")

		// Generate response
		processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     fmt.Sprintf("TestInstructionExecutionUuid already exists i 'processTestInstructionExecutionReversedResponseChannelMap'"),
				ErrorCodes:                   nil,
				ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
			},
			TestInstructionExecutionUuid:   processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid,
			ExpectedExecutionDuration:      nil,
			TestInstructionCanBeReExecuted: false,
		}

		// Return Response to Execution Server
		return processTestInstructionExecutionResponse, nil

	}

	// Create response channel to be able to get the 'processTestInstructionExecutionReversedResponse' back from Connector
	var processTestInstructionExecutionReversedResponseChannel processTestInstructionExecutionReversedResponseChannelType
	processTestInstructionExecutionReversedResponseChannel = make(chan *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse)

	// Create data for 'processTestInstructionExecutionReversedResponseChannelMap'
	var processTestInstructionExecutionReversedResponseMapData *processTestInstructionExecutionReversedResponseStruct
	processTestInstructionExecutionReversedResponseMapData = &processTestInstructionExecutionReversedResponseStruct{
		testInstructionExecutionUuid:                                    processTestInstructionExecutionRequest.TestInstruction.TestInstructionUuid,
		processTestInstructionExecutionReversedResponseChannelReference: &processTestInstructionExecutionReversedResponseChannel,
		savedInMapTimeStamp:                                             time.Now(),
	}

	// Save 'processTestInstructionExecutionReversedResponseChannelData' in Map
	//processTestInstructionExecutionReversedResponseChannelMap[processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid] = processTestInstructionExecutionReversedResponseMapData
	saveToTestInstructionExecutionReversedResponseChannelMap(processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid, processTestInstructionExecutionReversedResponseMapData)

	// Handle reversed response from Connector
	var testInstructionExecutionReversedResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse
	var gotReveresedResponseFromConnector bool
	go func() {
		testInstructionExecutionReversedResponse = <-processTestInstructionExecutionReversedResponseChannel
		gotReveresedResponseFromConnector = true
	}()

	// Create response channel to be able to get response when TestInstructionExecution is sent to Connector
	var executionResponseChannel executionResponseChannelType
	executionResponseChannel = make(chan executionResponseChannelStruct)

	// Create message to be sent to stream-server
	executionForwardChannelMessage := executionForwardChannelStruct{
		processTestInstructionExecutionReveredRequest: processTestInstructionExecutionRequest,
		executionResponseChannelReference:             &executionResponseChannel,
		isKeepAliveMessage:                            false,
	}

	// Send TestInstructionExecution to Stream-server, to later be sent to Connector, over channel
	executionForwardChannel <- executionForwardChannelMessage

	// Wait for response from stream-server that message has been sent TODO create some maximum time before clearing channel
	executionResponseChannelMessage := <-executionResponseChannel

	if executionResponseChannelMessage.testInstructionExecutionIsSentToConnector == false ||
		executionResponseChannelMessage.err != nil {
		// Message failed to be sent to Connector

		// Generate response
		processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     fmt.Sprintf("Message couldn't be sent to Connector, error: '%s'", executionResponseChannelMessage.err.Error()),
				ErrorCodes:                   nil,
				ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
			},
			TestInstructionExecutionUuid:   processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid,
			ExpectedExecutionDuration:      nil,
			TestInstructionCanBeReExecuted: false,
		}

	} else {
		// Message succeeded to be sent to Connector

		// Wait for response from 'processTestInstructionExecutionReversedResponseChannel' which is run as go-routine stated above
		for {
			if gotReveresedResponseFromConnector == true {
				break
			}

		}

		// Generate response
		processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
				AckNack:                      testInstructionExecutionReversedResponse.AckNackResponse.AckNack,
				Comments:                     testInstructionExecutionReversedResponse.AckNackResponse.Comments,
				ErrorCodes:                   testInstructionExecutionReversedResponse.AckNackResponse.ErrorCodes,
				ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
			},
			TestInstructionExecutionUuid:   processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid,
			ExpectedExecutionDuration:      testInstructionExecutionReversedResponse.ExpectedExecutionDuration,
			TestInstructionCanBeReExecuted: testInstructionExecutionReversedResponse.TestInstructionCanBeReExecuted,
		}

	}

	// Remove message from Map
	//delete(processTestInstructionExecutionReversedResponseChannelMap, processTestInstructionExecutionRequest.TestInstruction.TestInstructionUuid)
	deleteFromTestInstructionExecutionReversedResponseChannelMap(processTestInstructionExecutionRequest.TestInstruction.TestInstructionExecutionUuid)

	return processTestInstructionExecutionResponse, nil

}
