package common_config

// Create the PubSub-topic from Domain-Uuid
func GeneratePubSubTopicNameForTestInstructionExecution(
	thisDomainUuid string,
	thisExecutionDomainUuid string) (
	statusExecutionTopic string) {

	var pubSubTopicBase string
	pubSubTopicBase = TestInstructionExecutionPubSubTopicBase

	// Get the first 8 characters from DomainsUuid
	var shortedThisDomainsUuid string
	shortedThisDomainsUuid = thisDomainUuid[0:8] //ThisDomainsUuid[0:8]

	// Get the first 8 characters from 'thisExecutionDomainUuid'
	var shortedThisExecutionDomainUuid string
	shortedThisExecutionDomainUuid = thisExecutionDomainUuid[0:8]

	// Build PubSub-topic
	statusExecutionTopic = pubSubTopicBase + "-" + shortedThisDomainsUuid + "-" + shortedThisExecutionDomainUuid

	return statusExecutionTopic
}
