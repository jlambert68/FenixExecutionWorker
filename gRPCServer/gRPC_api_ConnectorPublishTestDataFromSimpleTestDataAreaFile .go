package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/messagesToGuiBuilderServer"
	"context"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/shared_code"
	"github.com/sirupsen/logrus"
)

// ConnectorPublishTestDataFromSimpleTestDataAreaFile
// Connector publish TestData from 'simple' file structure to Worker
func (s *fenixExecutionWorkerConnectorGrpcServicesServer) ConnectorPublishTestDataFromSimpleTestDataAreaFile(
	ctx context.Context,
	testDataFromSimpleTestDataAreaFileMessage *fenixExecutionWorkerGrpcApi.
		TestDataFromSimpleTestDataAreaFileMessage) (
	returnMessage *fenixExecutionWorkerGrpcApi.AckNackResponse,
	err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "5f547644-662b-4824-92a6-cb850cce0d99",
		//"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage,
	}).Debug("Incoming 'gRPCWorker- ConnectorPublishTestDataFromSimpleTestDataAreaFile'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "a87b449d-b548-4786-9e18-15322908190d",
	}).Debug("Outgoing 'gRPCWorker - ConnectorPublishTestDataFromSimpleTestDataAreaFile'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage = common_config.IsCallerUsingCorrectWorkerProtoFileVersion(
		userId,
		fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
			testDataFromSimpleTestDataAreaFileMessage.ClientSystemIdentification.ProtoFileVersionUsedByClient))
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Set up instance to use for execution gPRC
	var fenixGuiBuilderObject *messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct
	fenixGuiBuilderObject = &messagesToGuiBuilderServer.MessagesToGuiBuilderServerObjectStruct{Logger: s.logger}

	// Create gRPC-message towards GuiBuilderServer for 'TestDataFromSimpleTestDataAreaFileMessage'
	var testDataFromSimpleTestDataAreaFileMessageToBuilderServer *fenixTestCaseBuilderServerGrpcApi.
		TestDataFromSimpleTestDataAreaFileMessage
	testDataFromSimpleTestDataAreaFileMessageToBuilderServer = &fenixTestCaseBuilderServerGrpcApi.
		TestDataFromSimpleTestDataAreaFileMessage{
		ClientSystemIdentification: &fenixTestCaseBuilderServerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid: testDataFromSimpleTestDataAreaFileMessage.GetClientSystemIdentification().
				GetDomainUuid(),
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestBuilderServerProtoFileVersion()),
		},
		TestDataFromSimpleTestDataAreaFiles: nil,
	}

	// Convert incoming message to be used for outgoing 'TestDataFromSimpleTestDataAreaFileMessage'
	var testDataFromOneSimpleTestDataAreaFileMessage []*fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage
	for _, testDataFromSimpleTestDataAreaFile := range testDataFromSimpleTestDataAreaFileMessage.GetTestDataFromSimpleTestDataAreaFiles() {

		// Generate Headers for gRPC-message
		var headersForTestDataFromOneSimpleTestDataAreaFileForGrpc []*fenixTestCaseBuilderServerGrpcApi.
			HeaderForTestDataFromOneSimpleTestDataAreaFileMessage
		for _, header := range testDataFromSimpleTestDataAreaFile.HeadersForTestDataFromOneSimpleTestDataAreaFile {
			var headerForTestDataFromOneSimpleTestDataAreaFileForGrpc *fenixTestCaseBuilderServerGrpcApi.
				HeaderForTestDataFromOneSimpleTestDataAreaFileMessage
			headerForTestDataFromOneSimpleTestDataAreaFileForGrpc = &fenixTestCaseBuilderServerGrpcApi.
				HeaderForTestDataFromOneSimpleTestDataAreaFileMessage{
				ShouldHeaderActAsFilter: header.ShouldHeaderActAsFilter,
				HeaderName:              header.HeaderName,
				HeaderUiName:            header.HeaderName,
			}

			headersForTestDataFromOneSimpleTestDataAreaFileForGrpc = append(
				headersForTestDataFromOneSimpleTestDataAreaFileForGrpc,
				headerForTestDataFromOneSimpleTestDataAreaFileForGrpc)
		}

		// Generate the TestData-rows for gRPC-message
		var simpleTestDataRowMessageAsGrpc []*fenixTestCaseBuilderServerGrpcApi.SimpleTestDataRowMessage
		for _, tempTestDataRow := range testDataFromSimpleTestDataAreaFile.SimpleTestDataRows {

			// Convert one row of data into gRPC-version
			var testDataValuesAsStringSlice []string
			for _, testDataValue := range tempTestDataRow.GetTestDataValue() {
				testDataValuesAsStringSlice = append(testDataValuesAsStringSlice, testDataValue)
			}

			var tempTestDataRowAsGrpc *fenixTestCaseBuilderServerGrpcApi.SimpleTestDataRowMessage
			tempTestDataRowAsGrpc = &fenixTestCaseBuilderServerGrpcApi.SimpleTestDataRowMessage{TestDataValue: testDataValuesAsStringSlice}

			// Add row to slice of rows
			simpleTestDataRowMessageAsGrpc = append(simpleTestDataRowMessageAsGrpc, tempTestDataRowAsGrpc)
		}

		var oneSimpleTestDataAreaFileMessage *fenixTestCaseBuilderServerGrpcApi.
			TestDataFromOneSimpleTestDataAreaFileMessage
		oneSimpleTestDataAreaFileMessage = &fenixTestCaseBuilderServerGrpcApi.
			TestDataFromOneSimpleTestDataAreaFileMessage{
			TestDataDomainUuid:         testDataFromSimpleTestDataAreaFile.TestDataDomainUuid,
			TestDataDomainName:         testDataFromSimpleTestDataAreaFile.TestDataDomainName,
			TestDataDomainTemplateName: testDataFromSimpleTestDataAreaFile.TestDataDomainTemplateName,
			TestDataAreaUuid:           testDataFromSimpleTestDataAreaFile.TestDataAreaUuid,
			TestDataAreaName:           testDataFromSimpleTestDataAreaFile.TestDataAreaName,
			HeadersForTestDataFromOneSimpleTestDataAreaFile: nil,
			SimpleTestDataRows:            nil,
			TestDataFileSha256Hash:        testDataFromSimpleTestDataAreaFile.TestDataFileSha256Hash,
			ImportantDataInFileSha256Hash: testDataFromSimpleTestDataAreaFile.ImportantDataInFileSha256Hash,
		}

		// Add 'HeadersForTestDataFromOneSimpleTestDataAreaFile'
		var tempHeadersForTestDataFromOneSimpleTestDataAreaFiles []*fenixTestCaseBuilderServerGrpcApi.
			HeaderForTestDataFromOneSimpleTestDataAreaFileMessage

		for _, tempHeadersForTestDataFromOneSimpleTestDataAreaFileFromConnector := range testDataFromSimpleTestDataAreaFile.
			GetHeadersForTestDataFromOneSimpleTestDataAreaFile() {

			// Create the Header towards GuiServer
			var tempHeadersForTestDataFromOneSimpleTestDataAreaFile *fenixTestCaseBuilderServerGrpcApi.
				HeaderForTestDataFromOneSimpleTestDataAreaFileMessage

			tempHeadersForTestDataFromOneSimpleTestDataAreaFile = &fenixTestCaseBuilderServerGrpcApi.
				HeaderForTestDataFromOneSimpleTestDataAreaFileMessage{
				ShouldHeaderActAsFilter: tempHeadersForTestDataFromOneSimpleTestDataAreaFileFromConnector.
					GetShouldHeaderActAsFilter(),
				HeaderName: tempHeadersForTestDataFromOneSimpleTestDataAreaFileFromConnector.
					GetHeaderName(),
				HeaderUiName: tempHeadersForTestDataFromOneSimpleTestDataAreaFileFromConnector.
					GetHeaderUiName(),
			}

			// Add Header to slice of header
			tempHeadersForTestDataFromOneSimpleTestDataAreaFiles = append(
				tempHeadersForTestDataFromOneSimpleTestDataAreaFiles,
				tempHeadersForTestDataFromOneSimpleTestDataAreaFile)

		}

		// Add Headers to message towards GuiServer
		oneSimpleTestDataAreaFileMessage.HeadersForTestDataFromOneSimpleTestDataAreaFile =
			tempHeadersForTestDataFromOneSimpleTestDataAreaFiles

		// Add 'SimpleTestDataRows'
		var tempSimpleTestDataRowForTestDataFromOneSimpleTestDataAreaFiles []*fenixTestCaseBuilderServerGrpcApi.
			SimpleTestDataRowMessage

		for _, tempSimpleTestDataRowsFromConnector := range testDataFromSimpleTestDataAreaFile.GetSimpleTestDataRows() {

			// Create the Header towards GuiServer
			var tempSimpleTestDataRowForTestDataFromOneSimpleTestDataAreaFile *fenixTestCaseBuilderServerGrpcApi.
				SimpleTestDataRowMessage

			tempSimpleTestDataRowForTestDataFromOneSimpleTestDataAreaFile = &fenixTestCaseBuilderServerGrpcApi.
				SimpleTestDataRowMessage{
				TestDataValue: tempSimpleTestDataRowsFromConnector.GetTestDataValue()}

			// Add TestDataRow to slice of TestDataRows
			tempSimpleTestDataRowForTestDataFromOneSimpleTestDataAreaFiles = append(
				tempSimpleTestDataRowForTestDataFromOneSimpleTestDataAreaFiles,
				tempSimpleTestDataRowForTestDataFromOneSimpleTestDataAreaFile)

		}

		// Add Headers to message towards GuiServer
		oneSimpleTestDataAreaFileMessage.SimpleTestDataRows =
			tempSimpleTestDataRowForTestDataFromOneSimpleTestDataAreaFiles

		// Add one full TestDataFile to slice of all TestDataFiles
		testDataFromOneSimpleTestDataAreaFileMessage = append(testDataFromOneSimpleTestDataAreaFileMessage,
			oneSimpleTestDataAreaFileMessage)
	}
	// Add converted messages to outgoing message to BuilderServer
	testDataFromSimpleTestDataAreaFileMessageToBuilderServer.
		TestDataFromSimpleTestDataAreaFiles = testDataFromOneSimpleTestDataAreaFileMessage

	// Get Message to sign to prove identity
	succeededToSend, responseMessage, messageToSign := fenixGuiBuilderObject.
		SendGetMessageToSignToProveCallerIdentity()

	if succeededToSend == false {
		s.logger.WithFields(logrus.Fields{
			"id":              "f178612d-423c-4cab-8ec5-46f495baabc5",
			"responseMessage": responseMessage,
		}).Error("Got some error when sending 'GetMessageToSignToProveCallerIdentity'")

		// Generate response
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:    succeededToSend,
			Comments:   fmt.Sprintf("Messagage from BuilderServer: '%s'", responseMessage),
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
				common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil

	}

	// Specify the service account to be used when signing
	var serviceAccountUsedWhenSigning string
	serviceAccountUsedWhenSigning = fmt.Sprintf("projects/-/serviceAccounts/%s",
		common_config.ServiceAccountUsedForSigningMessage)

	// Sign Message to prove Identity to BuilderServer
	var hashOfSignature string
	var hashedKeyId string
	if common_config.ExecutionLocationForWorker == common_config.GCP {

		hashOfSignature, hashedKeyId, err = shared_code.SignMessageToProveIdentityToBuilderServer(
			messageToSign,
			serviceAccountUsedWhenSigning,
			true)

	} else {

		hashOfSignature, hashedKeyId, err = shared_code.SignMessageToProveIdentityToBuilderServer(
			messageToSign,
			serviceAccountUsedWhenSigning,
			false)
	}

	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"id":  "bb8284b2-cc4f-440f-8474-a2ea12adduce1",
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
		MessageToBeSigned: messageToSign,
		HashOfSignature:   hashOfSignature,
		HashedKeyId:       hashedKeyId,
	}
	testDataFromSimpleTestDataAreaFileMessageToBuilderServer.
		SignedMessageByWorkerServiceAccount = signedMessageByWorkerServiceAccountMessage

	// Publish Supported template repository connection parameters to TestCaseBuilderServer
	succeededToSend, responseMessage = fenixGuiBuilderObject.
		SendConnectorPublishTestDataFromSimpleTestDataAreaFileMessageToBuilderServer(
			testDataFromSimpleTestDataAreaFileMessageToBuilderServer)

	// Create Error Codes
	var errorCodes []fenixExecutionWorkerGrpcApi.ErrorCodesEnum

	// Generate response
	returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      succeededToSend,
		Comments:                     fmt.Sprintf("Messagage from BuilderServer: '%s'", responseMessage),
		ErrorCodes:                   errorCodes,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return returnMessage, nil

}
