package main

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/gcp"
	iam_credentials "cloud.google.com/go/iam/credentials/apiv1"
	"context"
	"crypto/tls"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	iam_credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

// SignMessageToProveIdentityToBuilderServer
// Sign Message to be sent to BuilderServer
func signTest(
	messageToBeSigned string,
	serviceAccountUsedForSigning string,
	signerIsRunningInGCP bool) (
	signedMessage []byte,
	hashOfSignature string,
	hashedKeyId string,
	err error) {

	ctx := context.Background()

	// Initialize the client
	var credsClient *iam_credentials.IamCredentialsClient
	if signerIsRunningInGCP == true {
		// Caller is running in GCP

		// Set up the custom TLS configuration
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}

		// Create a new gRPC client connection with the custom TLS settings
		conn, err := grpc.DialContext(ctx, "iamcredentials.googleapis.com:443", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		if err != nil {
			log.Fatalf("Failed to dial IAM Credentials API: %v", err)
		}
		defer conn.Close()

		credsClient, err = iam_credentials.NewIamCredentialsClient(ctx, option.WithGRPCConn(conn))

	} else {
		// Caller is running locally
		credsClient, err = iam_credentials.NewIamCredentialsClient(ctx)
	}
	if err != nil {
		return nil, "", "", err
	}

	defer credsClient.Close()

	// The data to be signed
	data := []byte(messageToBeSigned)

	// Request to sign a byte array with the service account's private key
	req := &iam_credentialspb.SignBlobRequest{
		Name:    serviceAccountUsedForSigning,
		Payload: data,
	}

	var (
		returnAckNack bool
		returnMessage string
	)
	ctx, returnAckNack, returnMessage = gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GenerateTokeForIamCredentials)
	if returnAckNack == false || len(returnMessage) > 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "1e46ea03-6a67-4ee1-853d-408d60b440d5",
			"err": err,
		}).Fatal("Problem getting the token")
	}

	// Call the API to sign the data
	var signResponse *iam_credentialspb.SignBlobResponse
	signResponse, err = credsClient.SignBlob(ctx, req)
	if err != nil {
		return nil, "", "", err
	}

	signedMessage = signResponse.SignedBlob

	// Hash the signature
	hashOfSignature = fenixSyncShared.HashSingleValue(string(signedMessage))

	// Extract KeyId used when signing
	var keyId string
	keyId = signResponse.GetKeyId()

	// Hash KeyId
	hashedKeyId = fenixSyncShared.HashSingleValue(keyId)

	// Return result
	return signedMessage, hashOfSignature, hashedKeyId, err
}
