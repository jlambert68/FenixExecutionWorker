package outgoingPubSubMessages

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
	grpcMetadata "google.golang.org/grpc/metadata"
	"time"
)

// Generate Google access token. Used when running in GCP
func (pubSubMessagesToConnectorObject *PubSubMessagesToConnectorObjectStruct) generateGCPAccessToken(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Only create the token if there is none, or it has expired
	if pubSubMessagesToConnectorObject.gcpAccessToken == nil || pubSubMessagesToConnectorObject.gcpAccessToken.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.
		tokenSource, err := idtoken.NewTokenSource(ctx, "https://www.googleapis.com/auth/pubsub") //"https://"+common_config.FenixExecutionServerAddress)
		if err != nil {
			pubSubMessagesToConnectorObject.Logger.WithFields(logrus.Fields{
				"ID":  "8ba622d8-b4cd-46c7-9f81-d9ade2568eca",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		if err != nil {
			pubSubMessagesToConnectorObject.Logger.WithFields(logrus.Fields{
				"ID":  "346c6b45-bec5-425c-8ac4-dd8f5e92961f",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			pubSubMessagesToConnectorObject.Logger.WithFields(logrus.Fields{
				"ID": "0937b872-31dd-43de-b15f-d2b52594ebc3",
				//"token": token,
			}).Debug("Got Bearer Token")
		}

		pubSubMessagesToConnectorObject.gcpAccessToken = token

	}

	pubSubMessagesToConnectorObject.Logger.WithFields(logrus.Fields{
		"ID": "f76007b9-7748-4aed-963f-ee388068ac19",
		//"FenixExecutionWorkerObject.gcpAccessToken": pubSubMessagesToConnectorObject.gcpAccessToken,
	}).Debug("Will use Bearer Token")

	// Add token to GrpcServer Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+pubSubMessagesToConnectorObject.gcpAccessToken.AccessToken)

	return appendedCtx, true, ""

}
