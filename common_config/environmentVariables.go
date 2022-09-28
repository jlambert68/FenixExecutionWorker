package common_config

// ***********************************************************************************************************
// The following variables receives their values from environment variables

// Where is the Worker running
var ExecutionLocationForWorker ExecutionLocationTypeType

// Where is Fenix Execution Server running
var ExecutionLocationForFenixExecutionServer ExecutionLocationTypeType

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
	FenixExecutionServerAddress    string
	FenixExecutionWorkerServerPort int
)
