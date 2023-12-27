package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/messagesToGuiBuilderServer"
	"context"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TypeAndStructs"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/shared_code"
	"github.com/sirupsen/logrus"
	"os"
)

// ConnectorReportCompleteTestInstructionExecutionResult
// When a TestInstruction has been fully executed the Execution Connector use this to inform the results of the execution result to the Worker
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	ctx context.Context,
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage *fenixExecutionWorkerGrpcApi.
		SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage) (
	returnMessage *fenixExecutionWorkerGrpcApi.AckNackResponse,
	err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "38b45573-c71e-4059-afeb-cd2deef237fb",
		//"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage,
	}).Debug("Incoming 'gRPCWorker- ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "1e5128bf-4a60-477e-a88e-ef08efc5642d",
	}).Debug("Outgoing 'gRPCWorker - ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage = common_config.IsCallerUsingCorrectWorkerProtoFileVersion(
		userId,
		fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage.ClientSystemIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Set up instance to use for execution gPRC
	var fenixGuiBuilderObject *messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct
	fenixGuiBuilderObject = &messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct{Logger: s.logger}

	// Convert back supported TestInstructions, TestInstructionContainers and Allowed Users message from a gRPC-Worker version of the message and check correctness of Hashes
	var testInstructionsAndTestInstructionContainersFromGrpcWorkerMessage *TestInstructionAndTestInstuctionContainerTypes.TestInstructionsAndTestInstructionsContainersStruct
	testInstructionsAndTestInstructionContainersFromGrpcWorkerMessage, err = shared_code.
		GenerateStandardFromGrpcWorkerMessageForTestInstructionsAndUsers(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":    "ac669e80-28ea-4002-97f9-a413063e83c3",
			"error": err,
		}).Fatalln("Problem when Convert back supported TestInstructions, TestInstructionContainers and " +
			"Allowed Users message from a gRPC-Worker version of the message and check correctness of Hashes " +
			"in 'ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")
	}

	// Verify recreated Hashes from gRPC-Worker-message
	var errorSliceWorker []error
	errorSliceWorker = shared_code.VerifyTestInstructionAndTestInstructionContainerAndUsersMessageHashesAndDomain(
		TypeAndStructs.DomainUUIDType(common_config.ThisDomainsUuid),
		testInstructionsAndTestInstructionContainersFromGrpcWorkerMessage)
	if errorSliceWorker != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":               "c05d5d90-cf78-4979-ae72-2c53a1aa12c9",
			"errorSliceWorker": errorSliceWorker,
		}).Error("Problem when recreated Hashes from gRPC-Worker-message " +
			"in 'ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

		// Loop error messages and concatenate into one string
		var errorMessageBackToConnector string
		for _, errorFromWorker := range errorSliceWorker {
			if len(errorMessageBackToConnector) == 0 {
				errorMessageBackToConnector = errorFromWorker.Error()
			} else {
				errorMessageBackToConnector = errorMessageBackToConnector + "; " + errorFromWorker.Error()
			}
		}

		// Create return message
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     errorMessageBackToConnector,
			ErrorCodes:                   []fenixExecutionWorkerGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil
	}

	// Create gRPC-message towards GuiBuilderServer for 'SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'

	// First
	// Convert back supported TestInstructions, TestInstructionContainers and Allowed Users message from a gRPC-Worker version of the message and check correctness of Hashes
	var testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.TestInstructionsAndTestInstructionsContainersStruct
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err = shared_code.
		GenerateStandardFromGrpcWorkerMessageForTestInstructionsAndUsers(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage)

	if err != nil {
		// Create return message
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     err.Error(),
			ErrorCodes:                   []fenixExecutionWorkerGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil
	}

	// Second
	// Verify recreated Hashes from gRPC-Builder-message
	var errorSliceBuilder []error
	errorSliceBuilder = shared_code.VerifyTestInstructionAndTestInstructionContainerAndUsersMessageHashesAndDomain(
		TypeAndStructs.DomainUUIDType(common_config.ThisDomainsUuid),
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	// If there are error then loop and concatenate error message to be sent to user
	if errorSliceBuilder != nil {
		var errToReturn string
		for _, errFromBuilder := range errorSliceBuilder {
			if len(errToReturn) == 0 {
				errToReturn = errFromBuilder.Error()
			} else {
				errToReturn = errToReturn + " - " + errFromBuilder.Error()
			}

		}

		// Create return message
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     errToReturn,
			ErrorCodes:                   []fenixExecutionWorkerGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil
	}

	// Third
	// Convert supported TestInstructions, TestInstructionContainers and Allowed Users message into a gRPC-Builder version of the message
	var supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcBuilderMessage *fenixTestCaseBuilderServerGrpcApi.SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcBuilderMessage, err = shared_code.
		GenerateTestInstructionAndTestInstructionContainerAndUserGrpcBuilderMessage(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage.GetConnectorDomain().GetConnectorsDomainUUID(),
			testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get Message to sign to prove identity
	succeededToSend, responseMessage, messageToSign := fenixGuiBuilderObject.
		SendGetMessageToSignToProveCallerIdentity()

	if succeededToSend == false {
		s.logger.WithFields(logrus.Fields{
			"id":              "7fb5a94a-7984-47f8-8c08-0d7da1446b94",
			"responseMessage": responseMessage,
		}).Error("Got some error when sending 'GetMessageToSignToProveCallerIdentity'")

		// Generate response
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      succeededToSend,
			Comments:                     fmt.Sprintf("Messagage from BuilderServer: '%s'", responseMessage),
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil

	}

	// Sign Message to prove Identity to BuilderServer
	var signedMessageAsByteSlice []byte
	signedMessageAsByteSlice, err = fenixGuiBuilderObject.SignMessageToProveIdentityToBuilderServer(messageToSign)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"id":  "532dff93-5786-4350-96a2-ddf977ee5ec5",
			"err": err,
		}).Error("Got some error when signing message")

		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     fmt.Sprintf("Got some error when signing message. '%s'", err.Error()),
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, err
	}

	// Add Signed message
	var signedMessageByWorkerServiceAccountMessage *fenixTestCaseBuilderServerGrpcApi.SignedMessageByWorkerServiceAccountMessage
	signedMessageByWorkerServiceAccountMessage = &fenixTestCaseBuilderServerGrpcApi.SignedMessageByWorkerServiceAccountMessage{
		MessageToBeSigned:          messageToSign,
		Signature:                  string(signedMessageAsByteSlice),
		PublicKeyForServiceAccount: "",
	}
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcBuilderMessage.
		SignedMessageByWorkerServiceAccount = signedMessageByWorkerServiceAccountMessage

	// Publish Supported TestInstructions, TestInstructionContainers And AllowedUsers To FenixGuiBuilderServer
	succeededToSend, responseMessage = fenixGuiBuilderObject.
		SendPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersToFenixGuiBuilderServer(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcBuilderMessage)

	if succeededToSend == false {
		s.logger.WithFields(logrus.Fields{
			"id":              "532dff93-5786-4350-96a2-ddf977ee5ec5",
			"responseMessage": responseMessage,
		}).Error("Got some error when sending 'CompleteTestInstructionExecutionResultToFenixExecutionServer'")

	}

	// Create Error Codes
	var errorCodes []fenixExecutionWorkerGrpcApi.ErrorCodesEnum

	// Generate response
	returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      succeededToSend,
		Comments:                     fmt.Sprintf("Messagage from ExecutionServer: '%s'", responseMessage),
		ErrorCodes:                   errorCodes,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return returnMessage, nil

}
