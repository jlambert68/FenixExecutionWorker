package main

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/gcp"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

func (fenixExecutionWorkerObject *fenixExecutionWorkerObjectStruct) InitLogger(filename string) {
	fenixExecutionWorkerObject.logger = logrus.StandardLogger()

	switch common_config.LoggingLevel {

	case logrus.DebugLevel:
		log.Println("'common_config.LoggingLevel': ", common_config.LoggingLevel)

	case logrus.InfoLevel:
		log.Println("'common_config.LoggingLevel': ", common_config.LoggingLevel)

	case logrus.WarnLevel:
		log.Println("'common_config.LoggingLevel': ", common_config.LoggingLevel)

	default:
		log.Println("Not correct value for debugging-level, this was used: ", common_config.LoggingLevel)
		os.Exit(0)

	}

	logrus.SetLevel(common_config.LoggingLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		DisableSorting:  true,
	})

	//If no file then set standard out

	if filename == "" {
		fenixExecutionWorkerObject.logger.Out = os.Stdout

	} else {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			fenixExecutionWorkerObject.logger.Out = file
		} else {
			log.Println("Failed to log to file, using default stderr")
		}
	}

	// Should only be done from init functions
	//grpclog.SetLoggerV2(grpclog.NewLoggerV2(logger.Out, logger.Out, logger.Out))

	gcp.Gcp.SetLogger(fenixExecutionWorkerObject.logger)

}
