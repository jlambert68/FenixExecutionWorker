package outgoingPubSubMessages

import (
	"FenixExecutionWorker/common_config"
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
)

func Publish(msg string) (returnMessageAckNack bool, returnMessageString string, err error) {
	projectID := common_config.GcpProject
	topicID := "SubCustody-ProcessTestInstructionExecutionRequest" //"projects/mycloud-run-project/topics/testinstruction-execution"
	// msg := "Hello World"

	var pubSubClient *pubsub.Client
	var opts []grpc.DialOption

	ctx := context.Background()

	// Add Access token

	//ctx, returnMessageAckNack, returnMessageString = gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GenerateTokenForPubSub)
	returnMessageAckNack = true
	if returnMessageAckNack == false {
		return returnMessageAckNack, returnMessageString, nil
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/jlambert/Downloads/mycloud-run-project-a35e47ac3dc3.json")
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	if true {
		//ctx = context.Background()
		pubSubClient, err = pubsub.NewClient(ctx, projectID)

	} else {

		//When running on GCP then use credential otherwise not
		if common_config.ExecutionLocationForWorker == common_config.GCP {

			var creds credentials.TransportCredentials
			creds = credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})

			opts = []grpc.DialOption{
				grpc.WithTransportCredentials(creds),
			}

			pubSubClient, err = pubsub.NewClient(ctx, projectID, option.WithGRPCDialOption(opts[0]))

		} else {

			//ctx := context.Background()

			// Base credentials sourced from ADC or provided client options.
			var ts oauth2.TokenSource
			ts, err = impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
				TargetPrincipal: "fenix-cloudrun-runner@mycloud-run-project.iam.gserviceaccount.com",
				Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
				//Lifetime:        time.Duration(time.Minute * 60),
				// Optionally supply delegates.
				//Delegates: []string{"bar@project-id.iam.gserviceaccount.com"},
			})
			if err != nil {

				common_config.Logger.WithFields(logrus.Fields{
					"ID": "00dc111a-52e7-4ed9-af1e-b5fa2af6c669",
					"ts": ts,
				}).Error(fmt.Errorf("impersonate.CredentialsTokenSource(ctx: %w", err))

				return false, "", err

			}

			var serviceAccountKeyJson = []byte(`{XXX}`)

			//ctx = context.Background()
			tokenSource, err := idtoken.NewTokenSource(ctx,
				"https://www.googleapis.com/auth/cloud-platform",                                              //"https://www.googleapis.com/auth/pubsub",
				idtoken.WithCredentialsFile("/home/jlambert/Downloads/mycloud-run-project-a35e47ac3dc3.json")) // WithCredentialsJSON(serviceAccountKeyJson))
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"ID": "44930ec9-7083-4d42-b721-be6ed938360a",
				}).Error(fmt.Errorf("tokenSource, err := idtoken.NewTokenSource(ctx,: %w", err))

				return false, "", err

			}

			_, err = tokenSource.Token()
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"ID": "65d121ef-7f89-4a19-a4a2-c3d9757daccf",
				}).Error(fmt.Errorf("tokenSource.Token(): %w", err))

				return false, "", err
			}
			//fmt.Println(token)
			//ctx = context.Background()

			//appendedCtx := grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)
			//pubSubClient, err = pubsub.NewClient(appendedCtx, projectID) //, option.WithTokenSource(tokenSource))
			pubSubClient, err = pubsub.NewClient(ctx, projectID, option.WithCredentialsJSON(serviceAccountKeyJson)) //.WithTokenSource(tokenSource))
		}
	}

	if err != nil {
		return false, "", fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer pubSubClient.Close()

	var pubSubTopic *pubsub.Topic
	var pubSubResult *pubsub.PublishResult
	pubSubTopic = pubSubClient.Topic(topicID)
	pubSubResult = pubSubTopic.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the pubSubResult is returned and a server-generated
	// ID is returned for the published message.
	id, err := pubSubResult.Get(ctx)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID": "dc8bb67a-2caf-4a46-8a5c-598e253515c5",
		}).Error(fmt.Errorf("pubsub: pubSubResult.Get: %w", err))

	}

	common_config.Logger.WithFields(logrus.Fields{
		"ID": "8da81faa-a2a9-4130-83c8-e90b8fbbb955",
		//"token": token,
	}).Debug(fmt.Sprintf("Published a message; msg ID: %v", id))

	return true, "", err
}
