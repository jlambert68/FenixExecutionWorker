package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/workerEngine"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ConnectorReportCompleteTestInstructionExecutionResult
// When a TestInstruction has been fully executed the Execution Connector use this to inform the results of the execution result to the Worker
func (s *fenixExecutionWorkerGrpcServicesServer) ConnectorReportCompleteTestInstructionExecutionResult(ctx context.Context, finalTestInstructionExecutionResultMessage *fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage) (ackNackResponse *fenixExecutionWorkerGrpcApi.AckNackResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "d85d5be5-33e8-4b8e-9577-50e4b84df389",
	}).Debug("Incoming 'gRPCServer - ConnectorReportCompleteTestInstructionExecutionResult'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "0a46c193-d37a-40bc-8c7b-43c1c2e02898",
	}).Debug("Outgoing 'gRPCServer - ConnectorReportCompleteTestInstructionExecutionResult'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(finalTestInstructionExecutionResultMessage.ClientSystemIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return returnMessage, nil
	}

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
