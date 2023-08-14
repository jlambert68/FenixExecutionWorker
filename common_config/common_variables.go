package common_config

import "github.com/sirupsen/logrus"

// Used for keeping track of the proto file versions for ExecutionServer and this Worker
var highestFenixExecutionServerProtoFileVersion int32 = -1
var highestExecutionWorkerProtoFileVersion int32 = -1

var Logger *logrus.Logger
