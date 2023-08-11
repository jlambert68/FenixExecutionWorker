package outgoingPubSubMessages

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/gcp"
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
	"io"
	"log"
	"os"
)

func Publish(w io.Writer, msg string) (bool, string, error) {
	projectID := common_config.GcpProject
	topicID := "testinstruction-execution" //"projects/mycloud-run-project/topics/testinstruction-execution"
	// msg := "Hello World"

	var pubSubClient *pubsub.Client
	var err error
	var opts []grpc.DialOption

	ctx := context.Background()

	// Add Access token
	var returnMessageAckNack bool
	var returnMessageString string

	ctx, returnMessageAckNack, returnMessageString = gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GenerateTokenForPubSub)
	if returnMessageAckNack == false {
		return returnMessageAckNack, returnMessageString, nil
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/jlambert/Downloads/mycloud-run-project-dbd73951e789.json")

	if true {
		ctx = context.Background()
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
			/*var ts oauth2.TokenSource
			ts, err = impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
				TargetPrincipal: "fenix-cloudrun-runner@mycloud-run-project.iam.gserviceaccount.com",
				Scopes:          []string{"https://www.googleapis.com/auth/pubsub"},
				//Lifetime:        time.Duration(time.Minute * 60),
				// Optionally supply delegates.
				//Delegates: []string{"bar@project-id.iam.gserviceaccount.com"},
			})
			if err != nil {
				log.Fatal(err)
			}

			*/

			var serviceAccountKeyJson = []byte(`{
  "type": "service_account",
  "project_id": "mycloud-run-project",
  "private_key_id": "dbd73951e78992dd612a4aace2430cfcba1861b2",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDOMYvugNY8BpqG\nFDBN34RZW1Ygl8GM+OxyfGAEV7fSz2wUkULom6ixEWTqCCkoM9GBhuKixDhLiC4B\nRNqdcOEOGHzAgu11Rxt4pcVfdPQ2XeDIsN7WaWmzCZhAJQ5eQarmDkdag7oYIoSs\nqWJWxtuYJjx+OBVB+qTiXGci316yVCfHC6F9Y5QU6LOTCPxO8VMnLsCXlwd80bmr\nTstHMWx/Wr5p/51Z6t526SRbp0aG2ZvWeWvdWn8navusCpJpNx2CKY9BVaxs+v2T\nuIAMSV/nyrx6O8rfO2RijQjL46JkNgeJmP8SRkYF2C4KIP3n3kmcpKp8PugwNSbA\nT0/egwDXAgMBAAECggEAHRHIKBpa2cCWlXOQMdJiztQ9KsAqZe0MLMHTKZmSTXK3\nypiGJdQYLsqEfygiYUwY69lv50GhrChpT/18krjAyeNy5xMuVhvtyA8a6e/LpESM\n9c5VxEW9RKQEJnry8R/x75gwwBaVLGTlbpA80H4dpHAzlBnlCVXEXDNpyPVT3PEW\nP1/0X7it91QhepxpYG57xXVu2lCMekrQnltOZHN4Ec5BlShXHWBzpSkLku6NFHIR\nAKfsmttWvMcAEVSq8xcXtsC3WnlNLOpe0A51AwKu6hx3zRbosB/CHb4RTmt0nMuF\n0PwwwRzn3UvTHn/uLiGAkGLaRNyezJhF6T+jKQaHqQKBgQD+zpXN5v8ijeWlpFS0\neiU3slxn4/ikViAS7XTLx/SCKUP+XtAkwem1FaZkjkgACFlI5hKjPWCxMQ5RthsE\nuss7v0vDh8SzW1i2AyUUXsJdJwcLvO3ax2VqjpkdsLwkc5W0V0mG+6DKXr1DEEUW\nJ4RRFQ4SgAFkYrZ4GFdNTF0f/wKBgQDPKLFZL7VQvFv8XAMJ10Rg1io9HMGRN/mf\ngVrpktSqr3aJa+Y4QHZHcGW8oWI2RY8e+2c7zztbY+G0N422OTuipg2LzJxUJmm5\nFTdcnrLb/xS8dT8Jn9Lr5CgIvTvT/ciNem6QT1g46h+aTDeRvys5NwvDSR3lhREA\n4F2Eg0cfKQKBgG/vZNO9JFuTpkyr8iIOfocHLZzeAv5+bzyhX+udfYYohpwaHqnn\nHbnVNvF5p5uMD1z85Tcc4Xs1p6qxqxDa1ij7EldlLz3zZPcgTouyTQQK/wdjCcJQ\nUfcLiawHVb9Vn3BH4B8SL0J3eAEcBp2C4peT/kiWPwZQbwQ2/TMR1t8VAoGAEs3x\n01O2VL7UZsL/b5w/0759FREQLRt0qFr4oq02aswEqZLG0iJf2jpEsevAW4bS6BAQ\nHejAzZnyegZ08a6eUDRclG0dX7Ig9LENVnX9bGTqP/UwpLICVnFehPSQgrzNwLH1\naVoaewgdmEcE4FEHHml1wuNXOGds1LSJKKc2BIkCgYEAzek1m/j1ScyqLaE7uXUT\nFf40kEHSJUUcIKsht0w3mtmNjjaney8FSwl7INjMQskQ5kIbg5eriUsBeqsMyBSp\nzvmwLuEe5cw840v4NmNzbZhdDJgvKuf3LHMunwx4KnOrgoYAfR4V8b+weTwpnXJa\n4XN+ZEAEyS5z5M70L6todFQ=\n-----END PRIVATE KEY-----\n",
  "client_email": "fenix-cloudrun-runner@mycloud-run-project.iam.gserviceaccount.com",
  "client_id": "103988388556721747241",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/fenix-cloudrun-runner%40mycloud-run-project.iam.gserviceaccount.com",
  "universe_domain": "googleapis.com"
}`)

			ctx = context.Background()
			tokenSource, err := idtoken.NewTokenSource(ctx,
				"https://www.googleapis.com/auth/pubsub",
				idtoken.WithCredentialsJSON(serviceAccountKeyJson))
			if err != nil {
				log.Fatal(err)
			}

			token, err := tokenSource.Token()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(token)
			ctx = context.Background()

			appendedCtx := grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)
			pubSubClient, err = pubsub.NewClient(appendedCtx, projectID) //, option.WithTokenSource(tokenSource))
			//pubSubClient, err = pubsub.NewClient(ctx, projectID, option.WithTokenSource(tokenSource))
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
		return false, "", fmt.Errorf("pubsub: pubSubResult.Get: %w", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return true, "", nil
}
