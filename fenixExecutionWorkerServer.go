package main

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/gRPCServer"
	"FenixExecutionWorker/gcp"
	"fmt"
	uuidGenerator "github.com/google/uuid"
	"github.com/jlambert68/FenixSyncShared/pubSubHelpers"
	"github.com/sirupsen/logrus"
	"os"
)

// Used for only process cleanup once
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// Cleanup before close down application
		FenixExecutionWorkerObject.logger.WithFields(logrus.Fields{}).Info("Clean up and shut down servers")

		// Stop Backend GrpcServer Server
		FenixExecutionWorkerObject.GrpcServer.StopGrpcServer()

	}
}

func fenixExecutionWorkerMain() {

	// Create Unique Uuid for run time instance used as identification when communication with ExecutionWorker
	common_config.ApplicationRunTimeUuid = uuidGenerator.New().String()
	fmt.Println("sharedCode.ApplicationRunTimeUuid: " + common_config.ApplicationRunTimeUuid)

	// Set up BackendObject
	FenixExecutionWorkerObject = &fenixExecutionWorkerObjectStruct{
		logger:     nil,
		GrpcServer: &gRPCServer.FenixExecutionWorkerGrpcObjectStruct{},
	}

	// Initiate gcp_helper object
	gcp.Gcp = gcp.GcpObjectStruct{}

	// Init logger
	FenixExecutionWorkerObject.InitLogger("")

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	// Initiate Logger for gRPC-server
	FenixExecutionWorkerObject.GrpcServer.InitiateLogger(FenixExecutionWorkerObject.logger)

	// Initiate shared Logger
	common_config.InitiateLogger(FenixExecutionWorkerObject.logger)

	// Only check if Topics and Subscriptions exists of that hasn't previously been done
	if common_config.TopicAndSubscriptionsExists == false {

		// Initiate PubSub-code
		pubSubHelpers.InitiatePubSubFunctionality(common_config.GcpProject, common_config.Logger)

		// Create PubSub-Topic
		var pubSubTopicToLookFor string
		pubSubTopicToLookFor = common_config.GeneratePubSubTopicForTestInstructionExecutions()

		// Secure that PubSub Topic, DeadLetteringTopic and their Subscriptions exist
		var err error
		err = pubSubHelpers.CreateTopicDeadLettingAndSubscriptionIfNotExists(
			pubSubTopicToLookFor, common_config.TestInstructionExecutionPubSubTopicSchema)
		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":                   "dbc8cc8e-d83d-42f0-a757-7f30cf3b62eb",
				"Error":                err,
				"pubSubTopicToLookFor": pubSubTopicToLookFor,
			}).Error("Something went wrong when Creating 'PubSub-Topics and Subscriptions")

			os.Exit(0)

		}

		// Secure that checking and initiating only will be done once
		common_config.TopicAndSubscriptionsExists = true
	}

	/*
		msg := "{\n \"ProtoFileVersionUsedByClient\": \"VERSION_0_3\",\n \"TestInstruction\": {\n \"TestInstructionExecutionUuid\": \"e1865111-88a8-4db5-b408-65be20d85a1f\",\n \"TestInstructionUuid\": \"26d38886-c112-48ef-a20f-4da8fb9a5ccb\",\n \"TestInstructionName\": \"TestCaseSetUp\",\n \"MinorVersionNumber\": 1,\n \"TestInstructionAttributes\": [\n {\n \"TestInstructionAttributeUuid\": \"f4682904-8f60-447c-b851-e713f2b4a03d\",\n \"TestInstructionAttributeName\": \"ExpectedToBePassed\",\n \"AttributeValueAsString\": \"true\",\n \"AttributeValueUuid\": \"f4682904-8f60-447c-b851-e713f2b4a03d\"\n }\n ]\n },\n \"TestData\": {\n \"TestDataSetUuid\": \"8e9671bd-5ded-485a-a2e6-cf8a44a63109\",\n \"ManualOverrideForTestData\": [\n {\n \"TestDataSetAttributeUuid\": \"f4682904-8f60-447c-b851-e713f2b4a03d\",\n \"TestDataSetAttributeName\": \"ExpectedToBePassed\",\n \"TestDataSetAttributeValue\": \"f4682904-8f60-447c-b851-e713f2b4a03d\"\n }\n ]\n }\n}"
		result, returnMessageString, err := outgoingPubSubMessages.Publish(msg)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		if result {
			fmt.Printf("Message published successfully: %s\n", returnMessageString)
		} else {
			fmt.Printf("Message publish failed: %s\n", returnMessageString)
		}
	*/

	// Start Backend GrpcServer-server
	FenixExecutionWorkerObject.GrpcServer.InitGrpcServer(FenixExecutionWorkerObject.logger)

}
