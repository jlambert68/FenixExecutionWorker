package messagesToGuiBuilderServer

import (
	"FenixExecutionWorker/common_config"
	"crypto/tls"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
	grpcMetadata "google.golang.org/grpc/metadata"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"golang.org/x/net/context"
)

// ********************************************************************************************************************

// SetConnectionToFenixGuiBuilderServer - Set upp connection and Dial to FenixExecutionServer
func (fenixExecutionWorkerObject *MessagesToGuiBuilderServerObjectStruct) SetConnectionToFenixGuiBuilderServer() (err error) {

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
			remoteFenixGuiBuilderServerConnection, err = grpc.Dial(common_config.FenixGuiBuilderServerAddressToDial, opts...)
		} else {
			// Run Local
			remoteFenixGuiBuilderServerConnection, err = grpc.Dial(common_config.FenixGuiBuilderServerAddressToDial, grpc.WithInsecure())
		}

		if err != nil {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID": "1cb282c4-4864-42b6-b943-89d4ed5b5300",
				"common_config.FenixGuiBuilderServerAddressToDial": common_config.FenixGuiBuilderServerAddressToDial,
				"error message":      err,
				"dialAttemptCounter": dialAttemptCounter,
			}).Error("Did not connect to FenixGuiBuilderServer via gRPC")

			// Only return the error after last attempt
			if dialAttemptCounter >= numberOfDialAttempts {
				return err
			}

		} else {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID": "bd3ee1ee-9849-4455-97a7-b7fbdaad1705",
				"common_config.FenixGuiBuilderServerAddressToDial": common_config.FenixGuiBuilderServerAddressToDial,
			}).Debug("gRPC connection OK to FenixGuiBuilderServer")

			// Creates a new Clients
			fenixTestCaseBuilderServerGrpcWorkerServicesClient = fenixTestCaseBuilderServerGrpcApi.
				NewFenixTestCaseBuilderServerGrpcWorkerServicesClient(remoteFenixGuiBuilderServerConnection)

			return err

		}

		// Sleep for some time before retrying to connect
		time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenDialAttempts[dialAttemptCounter-1]))

	}

}

// Generate Google access token. Used when running in GCP
func (fenixExecutionWorkerObject *MessagesToGuiBuilderServerObjectStruct) generateGCPAccessToken(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Only create the token if there is none, or it has expired
	if fenixExecutionWorkerObject.gcpAccessToken == nil || fenixExecutionWorkerObject.gcpAccessToken.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.
		tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+common_config.FenixExecutionServerAddress)
		if err != nil {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID":  "9b993a21-5019-4d87-b2df-7963d7963b2c",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		if err != nil {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID":  "619535af-1388-4f3c-af99-16f8df9da86b",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID": "8b1ca089-0797-4ee6-bf9d-f9b06f606ae9",
				//"token": token,
			}).Debug("Got Bearer Token")
		}

		fenixExecutionWorkerObject.gcpAccessToken = token

	}

	fenixExecutionWorkerObject.Logger.WithFields(logrus.Fields{
		"ID": "d4de2ade-8dcb-4d02-b511-e255cc8e00d9",
		//"FenixExecutionWorkerObject.gcpAccessToken": fenixExecutionWorkerObject.gcpAccessToken,
	}).Debug("Will use Bearer Token")

	// Add token to GrpcServer Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+fenixExecutionWorkerObject.gcpAccessToken.AccessToken)

	return appendedCtx, true, ""

}
