package common_config

// ***********************************************************************************************************
// The following variables receives their values from environment variables

// Where is the Worker running
var ExecutionLocationForWorker ExecutionLocationTypeType

// ExecutionLocationForFenixExecutionServer
// Where is Fenix Execution Server running
var ExecutionLocationForFenixExecutionServer ExecutionLocationTypeType

// ExecutionLocationForFenixGuiBuilderServer
// Where is Fenix GuiBuilder Server running
var ExecutionLocationForFenixGuiBuilderServer ExecutionLocationTypeType

// Definitions for where client and Fenix Server is running
type ExecutionLocationTypeType int

// Constants used for where stuff is running
const (
	LocalhostNoDocker ExecutionLocationTypeType = iota
	LocalhostDocker
	GCP
)

// Address to Fenix Execution Server & Execution Worker, will have their values from Environment variables at startup
var (
	ApplicationRunTimeUuid string
	// Worker
	ExecutionWorkerServerPort int

	// Fenix Execution Server
	FenixExecutionServerAddress       string
	FenixExecutionServerPort          int
	FenixExecutionServerAddressToDial string

	// Fenix GuiBuilder Server
	FenixGuiBuilderServerAddress       string
	FenixGuiBuilderServerPort          int
	FenixGuiBuilderServerAddressToDial string

	GcpProject string

	AuthClientId     string
	AuthClientSecret string

	UsePubSubWhenSendingTestInstructionExecutions bool
	LocalServiceAccountPath                       string

	ThisDomainsUuid                         string
	TestInstructionExecutionPubSubTopicBase string

	// TestInstructionExecutionPubSubTopicSchema
	// Topic-schema name to be used when sending 'TestInstructionExecutions' to Connector
	TestInstructionExecutionPubSubTopicSchema string
)
