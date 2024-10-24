package gRPCServer

import (
	"FenixExecutionWorker/common_config"
	"context"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/shared_code"
	"github.com/sirupsen/logrus"
)

// BuilderServerAskWorkerToSignMessage
// BuilderServer ask Worker to sign a message to prove that 'SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers' was sent from Worker
func (s *fenixExecutionWorkerBuilderGrpcServicesServer) BuilderServerAskWorkerToSignMessage(
	ctx context.Context,
	signMessageRequest *fenixExecutionWorkerGrpcApi.SignMessageRequest) (
	signMessageResponse *fenixExecutionWorkerGrpcApi.SignMessageResponse,
	err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "38b45573-c71e-4059-afeb-cd2deef237fb",
	}).Debug("Incoming 'gRPCWorker- BuilderServerAskWorkerToSignMessage'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "1e5128bf-4a60-477e-a88e-ef08efc5642d",
	}).Debug("Outgoing 'gRPCWorker - BuilderServerAskWorkerToSignMessage'")

	// Calling system
	userId := "BuilderServer"

	// Check if Client is using correct proto files version
	var ackNackMessage *fenixExecutionWorkerGrpcApi.AckNackResponse
	ackNackMessage = common_config.IsCallerUsingCorrectWorkerProtoFileVersion(
		userId,
		signMessageRequest.GetProtoFileVersionUsedByClient())

	if ackNackMessage != nil {
		signMessageResponse = &fenixExecutionWorkerGrpcApi.SignMessageResponse{
			AckNackResponse:                     ackNackMessage,
			SignedMessageByWorkerServiceAccount: nil,
		}

		return signMessageResponse, nil
	}

	// Specify the service account to be used when signing
	var serviceAccountUsedWhenSigning string
	serviceAccountUsedWhenSigning = fmt.Sprintf("projects/-/serviceAccounts/%s",
		common_config.ServiceAccountUsedForSigningMessage)

	// Sign Message to prove Identity to BuilderServer
	var hashOfSignature string
	var hashedKeyId string
	if common_config.ExecutionLocationForWorker == common_config.GCP {
		// Worker is running in GCP
		hashOfSignature, hashedKeyId, err = shared_code.SignMessageToProveIdentityToBuilderServer(
			signMessageRequest.GetMessageToBeSigned(),
			serviceAccountUsedWhenSigning, true)

	} else {
		// Worker is running locally
		hashOfSignature, hashedKeyId, err = shared_code.SignMessageToProveIdentityToBuilderServer(
			signMessageRequest.GetMessageToBeSigned(),
			serviceAccountUsedWhenSigning,
			false)
	}

	if err != nil {

		s.logger.WithFields(logrus.Fields{
			"id":  "c2501fe3-f9f6-4bb9-a475-cb9d865df9d7",
			"err": err,
		}).Error("Got some error when signing Message from BuilderServer")

		signMessageResponse = &fenixExecutionWorkerGrpcApi.SignMessageResponse{
			AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     err.Error(),
				ErrorCodes:                   nil,
				ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
			},
			SignedMessageByWorkerServiceAccount: nil,
		}

		return signMessageResponse, nil
	}

	// Generate response
	signMessageResponse = &fenixExecutionWorkerGrpcApi.SignMessageResponse{
		AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		},
		SignedMessageByWorkerServiceAccount: &fenixExecutionWorkerGrpcApi.SignedMessageByWorkerServiceAccountMessage{
			MessageToBeSigned: signMessageRequest.GetMessageToBeSigned(),
			HashOfSignature:   hashOfSignature,
			HashedKeyId:       hashedKeyId,
		},
	}

	return signMessageResponse, nil

}
