package common_config

import "github.com/sirupsen/logrus"

// Used for keeping track of the proto file versions for ExecutionServer and this Worker
var highestFenixExecutionServerProtoFileVersion int32 = -1
var highestExecutionWorkerProtoFileVersion int32 = -1
var highestBuilderServerProtoFileVersion int32 = -1

var Logger *logrus.Logger

// TopicAndSubscriptionsExists
// When a check that Topic, DeadLettering-Topic and their Subscriptions exist then this variable is set to true
var TopicAndSubscriptionsExists bool

// TopicAndSubscriptionsExistsMap
// When a check that Topic, DeadLettering-Topic and their Subscriptions exist then this Map is checked
var TopicAndSubscriptionsExistsMap map[string]bool //map[ExecutionDomainUuid]true
