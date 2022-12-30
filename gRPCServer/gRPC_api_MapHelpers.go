package gRPCServer

import "sync"

var TestInstructionExecutionReversedResponseMapMutex = &sync.RWMutex{}

// Load Subscription from the TestInstructionExecutionReversedResponseChannel-Map
func loadFromTestInstructionExecutionReversedResponseChannelMap(
	testInstructionExecutionUuid string) (
	processTestInstructionExecutionReversedResponse *processTestInstructionExecutionReversedResponseStruct,
	existInMap bool) {

	// Lock Map for Reading
	TestInstructionExecutionReversedResponseMapMutex.RLock()

	// Read Map
	processTestInstructionExecutionReversedResponse, existInMap = processTestInstructionExecutionReversedResponseChannelMap[testInstructionExecutionUuid]

	//UnLock Map
	TestInstructionExecutionReversedResponseMapMutex.RUnlock()

	return processTestInstructionExecutionReversedResponse, existInMap
}

// Save TestInstructionExecutionReversedResponse to the TestInstructionExecutionReversedResponseChannel-Map
func saveToTestInstructionExecutionReversedResponseChannelMap(
	testInstructionExecutionUuid string,
	processTestInstructionExecutionReversedResponse *processTestInstructionExecutionReversedResponseStruct) {

	// Lock Map for Writing
	TestInstructionExecutionReversedResponseMapMutex.Lock()

	// Save to TestInstructionExecutionReversedResponseChannel-Map
	processTestInstructionExecutionReversedResponseChannelMap[testInstructionExecutionUuid] = processTestInstructionExecutionReversedResponse

	//UnLock Map
	TestInstructionExecutionReversedResponseMapMutex.Unlock()

}

// De Subscription to the TestInstructionExecutionReversedResponseChannel-Map
func deleteFromTestInstructionExecutionReversedResponseChannelMap(
	testInstructionExecutionUuid string) {

	// Lock Map for Writing
	TestInstructionExecutionReversedResponseMapMutex.Lock()

	// Delete from TestInstructionExecutionReversedResponseChannel-Map
	delete(processTestInstructionExecutionReversedResponseChannelMap, testInstructionExecutionUuid)

	//UnLock Map
	TestInstructionExecutionReversedResponseMapMutex.Unlock()

}
