package main

import (
	"FenixExecutionWorker/common_config"
	"fmt"
	"github.com/sirupsen/logrus"
)

func signMessageTest() {

	// Specify the service account to be used when signing
	var serviceAccountUsedWhenSigning string
	serviceAccountUsedWhenSigning = fmt.Sprintf("projects/-/serviceAccounts/%s",
		common_config.ServiceAccountUsedForSigningMessage)

	// Sign Message to prove Identity to BuilderServer
	var hashOfSignature string
	var hashedKeyId string
	var signedMessage []byte

	var err error
	if common_config.ExecutionLocationForWorker == common_config.GCP {
		// Worker is running in GCP
		signedMessage, hashOfSignature, hashedKeyId, err = signTest(
			"MyMessageToBeSigned",
			serviceAccountUsedWhenSigning,
			true)

	} else {
		// Worker is running locally
		signedMessage, hashOfSignature, hashedKeyId, err = signTest(
			"MyMessageToBeSigned",
			serviceAccountUsedWhenSigning,
			false)
	}

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":              "e988ae91-6354-487d-94c2-46e62d2e5814",
			"err":             err,
			"signedMessage":   signedMessage,
			"hashOfSignature": hashOfSignature,
			"hashedKeyId":     hashedKeyId,
		}).Fatal("Got some error when signing Message from BuilderServer")

	}

	common_config.Logger.WithFields(logrus.Fields{
		"id":              "9b6f09e8-429f-43ec-ac09-13d0349d0e74",
		"err":             err,
		"signedMessage":   signedMessage,
		"hashOfSignature": hashOfSignature,
		"hashedKeyId":     hashedKeyId,
	}).Info("Success when signing Message from BuilderServer")

}
