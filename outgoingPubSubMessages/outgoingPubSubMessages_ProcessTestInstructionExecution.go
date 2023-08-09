package outgoingPubSubMessages

import (
	"FenixExecutionWorker/common_config"
	"FenixExecutionWorker/gcp"
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
)

func Publish(w io.Writer, msg string) (bool, string, error) {
	projectID := common_config.GcpProject
	topicID := "testinstruction-execution" //"projects/mycloud-run-project/topics/testinstruction-execution"
	// msg := "Hello World"

	ctx := context.Background()

	// Add Access token
	var returnMessageAckNack bool
	var returnMessageString string

	ctx, returnMessageAckNack, returnMessageString = gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GenerateTokenForPuSub)
	if returnMessageAckNack == false {
		return returnMessageAckNack, returnMessageString, nil
	}

	//When running on GCP then use credential otherwise not
	var opts []grpc.DialOption

	//var httpClient *http.Client
	if common_config.ExecutionLocationForFenixExecutionServer == common_config.GCP {
		//var creds credentials.TransportCredentials
		//creds = credentials.NewTLS(&tls.Config{
		//	InsecureSkipVerify: true,
		//})

		//opts = []grpc.DialOption{
		//	grpc.WithTransportCredentials(creds),
		//}

		/*
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true, // Insecure: skip certificate verification
			}

			transport := &http.Transport{
				TLSClientConfig: tlsConfig,
			}

			httpClient = &http.Client{
				Transport: transport,
				//Timeout:   10 * time.Second,
			}
		*/

		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})

		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}
	}

	client, err := pubsub.NewClient(ctx, projectID, option.WithGRPCDialOption(opts[0]))
	if err != nil {
		return false, "", fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return false, "", fmt.Errorf("pubsub: result.Get: %w", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return true, "", nil
}
