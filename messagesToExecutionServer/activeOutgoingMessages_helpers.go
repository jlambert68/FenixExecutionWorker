package messagesToExecutionServer

import (
	"FenixExecutionWorker/common_config"
	"crypto/tls"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
	grpcMetadata "google.golang.org/grpc/metadata"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"golang.org/x/net/context"
)

// ********************************************************************************************************************

// SetConnectionToFenixTestExecutionServer - Set upp connection and Dial to FenixExecutionServer
func (fenixExecutionWorkerObject *MessagesToExecutionServerObjectStruct) SetConnectionToFenixTestExecutionServer() (err error) {

	// slice with sleep time, in milliseconds, between each attempt to Dial to Server
	var sleepTimeBetweenDialAttempts []int
	sleepTimeBetweenDialAttempts = []int{100, 100, 200, 200, 300, 300, 500, 500, 600, 1000} // Total: 3.6 seconds

	var opts []grpc.DialOption

	// Do multiple attempts to do connection to Execution Server
	var numberOfDialAttempts int
	var dialAttemptCounter int
	numberOfDialAttempts = len(sleepTimeBetweenDialAttempts)
	dialAttemptCounter = 0

	for {

		dialAttemptCounter = dialAttemptCounter + 1

		//When running on GCP then use credential otherwise not
		if common_config.ExecutionLocationForFenixExecutionServer == common_config.GCP {
			creds := credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})

			opts = []grpc.DialOption{
				grpc.WithTransportCredentials(creds),
			}
		}

		// Set up connection to Fenix Execution Server
		// When run on GCP, use credentials
		if common_config.ExecutionLocationForFenixExecutionServer == common_config.GCP {
			// Run on GCP
			remoteFenixExecutionServerConnection, err = grpc.Dial(common_config.FenixExecutionServerAddressToDial, opts...)
		} else {
			// Run Local
			remoteFenixExecutionServerConnection, err = grpc.Dial(common_config.FenixExecutionServerAddressToDial, grpc.WithInsecure())
		}

		if err != nil {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID": "50b59b1b-57ce-4c27-aa84-617f0cde3100",
				"common_config.FenixExecutionServerAddressToDial": common_config.FenixExecutionServerAddressToDial,
				"error message":      err,
				"dialAttemptCounter": dialAttemptCounter,
			}).Error("Did not connect to FenixExecutionServer via gRPC")

			// Only return the error after last attempt
			if dialAttemptCounter >= numberOfDialAttempts {
				return err
			}

		} else {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID": "0c650bbc-45d0-4029-bd25-4ced9925a059",
				"common_config.FenixExecutionServerAddressToDial": common_config.FenixExecutionServerAddressToDial,
			}).Info("gRPC connection OK to FenixExecutionServer")

			// Creates a new Clients
			fenixExecutionServerGrpcClient = fenixExecutionServerGrpcApi.NewFenixExecutionServerGrpcServicesClient(remoteFenixExecutionServerConnection)

			return err

		}

		// Sleep for some time before retrying to connect
		time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenDialAttempts[dialAttemptCounter-1]))

	}

}

// Generate Google access token. Used when running in GCP
func (fenixExecutionWorkerObject *MessagesToExecutionServerObjectStruct) generateGCPAccessToken(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Only create the token if there is none, or it has expired
	if fenixExecutionWorkerObject.gcpAccessToken == nil || fenixExecutionWorkerObject.gcpAccessToken.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.
		tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+common_config.FenixExecutionServerAddress)
		if err != nil {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID":  "8ba622d8-b4cd-46c7-9f81-d9ade2568eca",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		if err != nil {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID":  "0cf31da5-9e6b-41bc-96f1-6b78fb446194",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID":    "8b1ca089-0797-4ee6-bf9d-f9b06f606ae9",
				"token": token,
			}).Debug("Got Bearer Token")
		}

		fenixExecutionWorkerObject.gcpAccessToken = token

	}

	fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"ID": "cd124ca3-87bb-431b-9e7f-e044c52b4960",
		"FenixExecutionWorkerObject.gcpAccessToken": fenixExecutionWorkerObject.gcpAccessToken,
	}).Debug("Will use Bearer Token")

	// Add token to GrpcServer Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+fenixExecutionWorkerObject.gcpAccessToken.AccessToken)

	return appendedCtx, true, ""

}
