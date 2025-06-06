package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/messagesToExecutionServer"
	"context"
	"fmt"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ConnectorReportCompleteTestInstructionExecutionResult
// When a TestInstruction has been fully executed the Execution Connector use this to inform the results of the execution result to the Worker
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorReportCompleteTestInstructionExecutionResult(
	ctx context.Context,
	finalTestInstructionExecutionResultMessage *fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage) (
	ackNackResponse *fenixExecutionWorkerGrpcApi.AckNackResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "44addf9e-2027-4b0d-9502-787194903e06",
		"finalTestInstructionExecutionResultMessage": finalTestInstructionExecutionResultMessage,
	}).Debug("Incoming 'gRPCServer - ConnectorReportCompleteTestInstructionExecutionResult'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "e658f679-be43-4427-9669-71d711223643",
	}).Debug("Outgoing 'gRPCServer - ConnectorReportCompleteTestInstructionExecutionResult'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(finalTestInstructionExecutionResultMessage.ClientSystemIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Set up instance to use for execution gPRC
	var fenixExecutionWorkerObject *messagesToExecutionServer.MessagesToExecutionServerObjectStruct
	fenixExecutionWorkerObject = &messagesToExecutionServer.MessagesToExecutionServerObjectStruct{Logger: s.logger}

	// Create ResponseVariables
	var responseVariablesGrpc []*fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage_ResponseVariableMessage
	for _, tempResponseVariable := range finalTestInstructionExecutionResultMessage.GetResponseVariables() {

		var responseVariableGrpc *fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage_ResponseVariableMessage
		responseVariableGrpc = &fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage_ResponseVariableMessage{
			ResponseVariableUuid:          tempResponseVariable.ResponseVariableUuid,
			ResponseVariableName:          tempResponseVariable.GetResponseVariableName(),
			ResponseVariableTypeUuid:      tempResponseVariable.GetResponseVariableTypeUuid(),
			ResponseVariableTypeName:      tempResponseVariable.GetResponseVariableTypeName(),
			ResponseVariableValueAsString: tempResponseVariable.GetResponseVariableValueAsString(),
		}

		// Append to list of Response Variables for ExecutionServer-gRPC
		responseVariablesGrpc = append(responseVariablesGrpc, responseVariableGrpc)
	}

	// Create LogPosts
	var logPostsGrpc []*fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage_LogPostMessage
	for _, tempLogPost := range finalTestInstructionExecutionResultMessage.GetLogPosts() {

		// Create FoundVersusExpectedValueForVariable
		var foundVersusExpectedValueForVariablesGrpc []*fenixExecutionServerGrpcApi.
			FinalTestInstructionExecutionResultMessage_LogPostMessage_FoundVersusExpectedValueForVariableMessage

		for _, tempFoundVersusExpectedValueForVariable := range tempLogPost.GetFoundVersusExpectedValueForVariable() {

			var foundVersusExpectedValueForVariableGrpc = &fenixExecutionServerGrpcApi.
				FinalTestInstructionExecutionResultMessage_LogPostMessage_FoundVersusExpectedValueForVariableMessage{
				VariableName:        tempFoundVersusExpectedValueForVariable.GetVariableName(),
				VariableDescription: tempFoundVersusExpectedValueForVariable.GetVariableDescription(),
				FoundVersusExpectedValue: &fenixExecutionServerGrpcApi.
					FinalTestInstructionExecutionResultMessage_LogPostMessage_FoundVersusExpectedValueMessage{
					VariableName:        tempFoundVersusExpectedValueForVariable.GetVariableName(),
					VariableDescription: tempFoundVersusExpectedValueForVariable.GetVariableDescription(),
					FoundValue:          tempFoundVersusExpectedValueForVariable.GetFoundVersusExpectedValue().GetFoundValue(),
					ExpectedValue:       tempFoundVersusExpectedValueForVariable.GetFoundVersusExpectedValue().GetExpectedValue(),
				},
			}

			// Append to list of Expected vs Found values for variable
			foundVersusExpectedValueForVariablesGrpc = append(foundVersusExpectedValueForVariablesGrpc, foundVersusExpectedValueForVariableGrpc)

		}

		var logPostGrpc *fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage_LogPostMessage
		logPostGrpc = &fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage_LogPostMessage{
			LogPostUuid:                         tempLogPost.GetLogPostUuid(),
			LogPostTimeStamp:                    tempLogPost.GetLogPostTimeStamp(),
			LogPostStatus:                       fenixExecutionServerGrpcApi.LogPostStatusEnum(tempLogPost.GetLogPostStatus()),
			LogPostText:                         tempLogPost.GetLogPostText(),
			FoundVersusExpectedValueForVariable: foundVersusExpectedValueForVariablesGrpc,
		}

		// Append to list of Response Variables for ExecutionServer-gRPC
		logPostsGrpc = append(logPostsGrpc, logPostGrpc)
	}

	// Create 'FinalTestInstructionExecutionResultMessage'
	var finalTestInstructionExecutionResultToServerMessage *fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage
	finalTestInstructionExecutionResultToServerMessage = &fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage{
		ClientSystemIdentification: &fenixExecutionServerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid:                   finalTestInstructionExecutionResultMessage.ClientSystemIdentification.DomainUuid,
			ProtoFileVersionUsedByClient: fenixExecutionServerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(common_config.GetHighestFenixExecutionServerProtoFileVersion()),
		},
		TestInstructionExecutionUuid:    finalTestInstructionExecutionResultMessage.GetTestInstructionExecutionUuid(),
		TestInstructionExecutionVersion: finalTestInstructionExecutionResultMessage.GetTestInstructionExecutionVersion(),
		MatureTestInstructionUuid:       finalTestInstructionExecutionResultMessage.GetMatureTestInstructionUuid(),
		TestInstructionExecutionStatus: fenixExecutionServerGrpcApi.TestInstructionExecutionStatusEnum(
			finalTestInstructionExecutionResultMessage.GetTestInstructionExecutionStatus()),
		TestInstructionExecutionStartTimeStamp: finalTestInstructionExecutionResultMessage.GetTestInstructionExecutionStartTimeStamp(),
		TestInstructionExecutionEndTimeStamp:   finalTestInstructionExecutionResultMessage.GetTestInstructionExecutionEndTimeStamp(),
		ResponseVariables:                      responseVariablesGrpc,
		LogPosts:                               logPostsGrpc,
	}

	succeededToSend, responseMessage := fenixExecutionWorkerObject.SendReportCompleteTestInstructionExecutionResultToFenixExecutionServer(finalTestInstructionExecutionResultToServerMessage)

	if succeededToSend == false {
		s.logger.WithFields(logrus.Fields{
			"id":              "3f2aab0a-d482-457c-845f-7c8537ee203d",
			"responseMessage": responseMessage,
		}).Error("Got some error when sending 'CompleteTestInstructionExecutionResultToFenixExecutionServer'")
	}

	// Create Error Codes
	var errorCodes []fenixExecutionWorkerGrpcApi.ErrorCodesEnum

	// Generate response
	ackNackResponse = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      succeededToSend,
		Comments:                     fmt.Sprintf("Messagage from ExecutionServer: '%s'", responseMessage),
		ErrorCodes:                   errorCodes,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return ackNackResponse, nil

}
