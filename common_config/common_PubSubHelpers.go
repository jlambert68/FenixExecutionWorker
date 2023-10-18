package common_config

// Create the PubSub-topic from Domain-Uuid
func GeneratePubSubTopicForTestInstructionExecutions() (statusExecutionTopic string) {

	var pubSubTopicBase string
	pubSubTopicBase = TestInstructionExecutionPubSubTopicBase

	var testerGuiApplicationUuid string
	testerGuiApplicationUuid = ThisDomainsUuid

	// Get the first 8 characters from TesterGui-ApplicationUuid
	var shortedAppUuid string
	shortedAppUuid = testerGuiApplicationUuid[0:8]

	// Build PubSub-topic
	statusExecutionTopic = pubSubTopicBase + "-" + shortedAppUuid

	return statusExecutionTopic
}
