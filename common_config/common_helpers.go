package common_config

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sort"
	"strconv"
	"time"
)

func HashValues(valuesToHash []string, isNonHashValue bool) string {

	hash_string := ""
	sha256_hash := ""

	// Concatenate array position to its content if it is a 'NonHashValue'
	if isNonHashValue == true {
		for valuePosition, value := range valuesToHash {
			valuesToHash[valuePosition] = value + strconv.Itoa(valuePosition)
		}
	}

	// Always sort values before hash them
	sort.Strings(valuesToHash)

	//Hash all values
	for _, valueToHash := range valuesToHash {
		hash_string = hash_string + valueToHash

		hash := sha256.New()
		hash.Write([]byte(hash_string))
		sha256_hash = hex.EncodeToString(hash.Sum(nil))
		hash_string = sha256_hash

	}

	return sha256_hash

}

// HashSingleValue HashSingleValue Hash a single value
func HashSingleValue(valueToHash string) (hashValue string) {

	hash := sha256.New()
	hash.Write([]byte(valueToHash))
	hashValue = hex.EncodeToString(hash.Sum(nil))

	return hashValue

}

// GenerateDatetimeTimeStampForDB
// Generate DataBaseTimeStamp, eg '2022-02-08 17:35:04.000000'
func GenerateDatetimeTimeStampForDB() (currentTimeStampAsString string) {

	timeStampLayOut := "2006-01-02 15:04:05.000000 -0700" //milliseconds
	currentTimeStamp := time.Now()
	currentTimeStampAsString = currentTimeStamp.Format(timeStampLayOut)

	return currentTimeStampAsString
}

// GenerateDatetimeFromTimeInputForDB
// Generate DataBaseTimeStamp, eg '2022-02-08 17:35:04.000000'
func GenerateDatetimeFromTimeInputForDB(currentTime time.Time) (currentTimeStampAsString string) {

	timeStampLayOut := "2006-01-02 15:04:05.000000 -0700" //milliseconds
	currentTimeStampAsString = currentTime.Format(timeStampLayOut)

	return currentTimeStampAsString
}

// ConvertGrpcTimeStampToStringForDB
// Convert a gRPCServer-timestamp into a string that can be used to store in the database
func ConvertGrpcTimeStampToStringForDB(grpcTimeStamp *timestamppb.Timestamp) (grpcTimeStampAsTimeStampAsString string) {
	grpcTimeStampAsTimeStamp := grpcTimeStamp.AsTime()

	timeStampLayOut := "2006-01-02 15:04:05.000000" //milliseconds

	grpcTimeStampAsTimeStampAsString = grpcTimeStampAsTimeStamp.Format(timeStampLayOut)

	return grpcTimeStampAsTimeStampAsString
}

/*
// ********************************************************************************************************************
// Get the highest FenixProtoFileVersionEnumeration
func GetHighestFenixTestDataProtoFileVersion() int32 {

	// Check if there already is a 'highestFenixProtoFileVersion' saved, if so use that one
	if highestFenixProtoFileVersion != -1 {
		return highestFenixProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixExecutionWorkerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestFenixProtoFileVersion = maxValue

	return highestFenixProtoFileVersion
}

*/

// InitiateLogger - Initiate local logger object
func InitiateLogger(logger *logrus.Logger) {

	Logger = logger
}
