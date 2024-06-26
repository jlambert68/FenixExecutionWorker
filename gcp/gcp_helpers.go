package gcp

import (
	"FenixExecutionWorker/common_config"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/sirupsen/logrus"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	grpcMetadata "google.golang.org/grpc/metadata"
	"html/template"
	"net/http"
	"time"
)

// GenerateTokenTargetType
// Type used to define
type GenerateTokenTargetType int

// GenerateTokenForExecutionServer
// Constants used to define what Token should be used for
const (
	GenerateTokenForGrpcTowardsExecutionServer GenerateTokenTargetType = iota
	GenerateTokenForPubSub
	GenerateTokeForIamCredentials
)

func (gcp *GcpObjectStruct) GenerateGCPAccessToken(ctx context.Context, tokenTarget GenerateTokenTargetType) (
	appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Chose correct method for authentication
	switch tokenTarget { // common_config.UseServiceAccount == true {

	case GenerateTokenForGrpcTowardsExecutionServer:

		// Only use Authorized used when running locally and ExecutionServer is in GCP
		if common_config.ExecutionLocationForWorker == common_config.LocalhostNoDocker &&
			common_config.ExecutionLocationForFenixExecutionServer == common_config.GCP {

			// Use Authorized user when targeting GCP from local
			appendedCtx, returnAckNack, returnMessage = gcp.GenerateGCPAccessTokenForAuthorizedUser(ctx)

		} else {
			// Use Authorized user
			appendedCtx, returnAckNack, returnMessage = gcp.generateGCPAccessToken(ctx)
		}

	case GenerateTokenForPubSub:
		// Extracts a token to be used for external PubSub-request
		if common_config.ExecutionLocationForWorker == common_config.LocalhostNoDocker {

			// Use Authorized user when targeting GCP from local
			appendedCtx, returnAckNack, returnMessage = gcp.GenerateGCPAccessTokenForAuthorizedUserPubSub(ctx)

		} else {
			// Use Authorized user
			appendedCtx, returnAckNack, returnMessage = gcp.generateGCPAccessTokenPubSub(ctx)
		}

	case GenerateTokeForIamCredentials:
		// Extracts a token to be used for external IAM Credentials
		if common_config.ExecutionLocationForWorker == common_config.LocalhostNoDocker {

			// Use Authorized user when targeting GCP from local
			appendedCtx, returnAckNack, returnMessage = gcp.GenerateGCPAccessTokenForAuthorizedUserPubSub(ctx)

		} else {
			// Use Authorized user
			appendedCtx, returnAckNack, returnMessage = gcp.generateGCPAccessTokenIamCredentials(ctx)
		}
	}

	return appendedCtx, returnAckNack, returnMessage

}

// Generate Google access token. Used when running in GCP
func (gcp *GcpObjectStruct) generateGCPAccessToken(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Only create the token if there is none, or it has expired
	if gcp.gcpAccessTokenForServiceAccounts == nil || gcp.gcpAccessTokenForServiceAccounts.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.

		tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+common_config.FenixExecutionServerAddress)
		if err != nil {
			gcp.logger.WithFields(logrus.Fields{
				"ID":  "11b41921-92fa-48ed-914f-0dde41282609",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		if err != nil {
			gcp.logger.WithFields(logrus.Fields{
				"ID":  "c1870620-d615-45e8-aaae-a1329d2ff4af",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			gcp.logger.WithFields(logrus.Fields{
				"ID":    "fee61402-aefa-4d4a-87ff-04b02c055366",
				"token": token,
			}).Debug("Got Bearer Token")
		}

		gcp.gcpAccessTokenForServiceAccounts = token

	}

	gcp.logger.WithFields(logrus.Fields{
		"ID":                                  "9bfd3d3a-7155-4f72-9cbc-e051f4544135",
		"gcp.gcpAccessTokenForServiceAccount": gcp.gcpAccessTokenForServiceAccounts,
	}).Debug("Will use Bearer Token")

	// Add token to GrpcServer Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.gcpAccessTokenForServiceAccounts.AccessToken)

	return appendedCtx, true, ""

}

// Generate Google access token for Pub Sub
func (gcp *GcpObjectStruct) generateGCPAccessTokenPubSub(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Initiate variable if 'nil'
	if gcp.GcpAccessTokenForServiceAccountsPubSub == nil {
		gcp.gcpAccessTokenForServiceAccounts = &oauth2.Token{}
	}

	// Only create the token if there is none, or it has expired
	if gcp.GcpAccessTokenForServiceAccountsPubSub == nil || gcp.GcpAccessTokenForServiceAccountsPubSub.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.

		tokenSource, err := idtoken.NewTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform") //"https://www.googleapis.com/auth/pubsub")

		if err != nil {
			gcp.logger.WithFields(logrus.Fields{
				"ID":  "ffb7cdcc-00f1-4560-9fd6-a45d2423230d",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		if err != nil {
			gcp.logger.WithFields(logrus.Fields{
				"ID":  "6f335c25-b020-4748-85ab-eda80e53b9a0",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			gcp.logger.WithFields(logrus.Fields{
				"ID":    "a17e40dc-e7fc-4d7e-afbc-072a4c21850b",
				"token": token,
			}).Debug("Got Bearer Token")
		}

		gcp.GcpAccessTokenForExternalPubSubRequests = token

	}

	gcp.logger.WithFields(logrus.Fields{
		"ID": "42427b1e-af8d-4153-9963-85c36a0f58cf",
		"gcp.GcpAccessTokenForExternalPubSubRequests":             gcp.GcpAccessTokenForExternalPubSubRequests,
		"gcp.GcpAccessTokenForExternalPubSubRequests.AccessToken": gcp.GcpAccessTokenForExternalPubSubRequests.AccessToken,
	}).Debug("Will use Bearer Token")

	// Add token to GrpcServer Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.GcpAccessTokenForExternalPubSubRequests.AccessToken)

	return appendedCtx, true, ""

}

// Generate Google access token for IAM Credentials
func (gcp *GcpObjectStruct) generateGCPAccessTokenIamCredentials(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Initiate variable if 'nil'
	if gcp.GcpAccessTokenForServiceAccountsIamCredentials == nil {
		gcp.GcpAccessTokenForServiceAccountsIamCredentials = &oauth2.Token{}
	}

	// Only create the token if there is none, or it has expired
	if gcp.GcpAccessTokenForServiceAccountsIamCredentials == nil || gcp.GcpAccessTokenForServiceAccountsIamCredentials.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.

		tokenSource, err := idtoken.NewTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform") //"https://www.googleapis.com/auth/pubsub")

		if err != nil {
			gcp.logger.WithFields(logrus.Fields{
				"ID":  "21394809-6175-40bf-b45e-7bce4a77cbdf",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		if err != nil {
			gcp.logger.WithFields(logrus.Fields{
				"ID":  "85841308-751d-4d16-bc27-c6448da6c399",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			gcp.logger.WithFields(logrus.Fields{
				"ID":    "63c06d0d-ce49-482c-bbb6-73e5413c717c",
				"token": token,
			}).Debug("Got Bearer Token")
		}

		gcp.GcpAccessTokenForServiceAccountsIamCredentials = token

	}

	gcp.logger.WithFields(logrus.Fields{
		"ID": "ef9487d9-69e5-4690-acdd-d6f18ff3f207",
		"gcp.GcpAccessTokenForServiceAccountsIamCredentials":             gcp.GcpAccessTokenForServiceAccountsIamCredentials,
		"gcp.GcpAccessTokenForServiceAccountsIamCredentials.AccessToken": gcp.GcpAccessTokenForServiceAccountsIamCredentials.AccessToken,
	}).Debug("Will use Bearer Token")

	// Add token to GrpcServer Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.GcpAccessTokenForServiceAccountsIamCredentials.AccessToken)

	return appendedCtx, true, ""

}

// DoneChannel - channel used for to close down local web server
var DoneChannel chan bool

func (gcp *GcpObjectStruct) GenerateGCPAccessTokenForAuthorizedUser(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Secure that User is initiated
	gcp.initiateUserObject()

	// Only create the token if there is none, or it has expired (or 5 minutes before expiration
	timeToCompareTo := time.Now().Add(-time.Minute * 5)
	if !(gcp.gcpAccessTokenForAuthorizedAccounts.IDToken == "" || gcp.gcpAccessTokenForAuthorizedAccounts.ExpiresAt.Before(timeToCompareTo)) {
		// We already have a ID-token that can be used, so return that
		appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.gcpAccessTokenForAuthorizedAccounts.IDToken)

		return appendedCtx, true, ""
	}

	// Need to create a new ID-token

	key := common_config.ApplicationRunTimeUuid // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30                        // 30 days
	isProd := false                             // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		// Use 'Fenix End User Authentication'
		google.New(
			common_config.AuthClientId,
			common_config.AuthClientSecret,
			"http://localhost:3000/auth/google/callback",
			"email", "profile"),
	)

	router := pat.New()

	router.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {

			fmt.Fprintln(res, err)

			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)

		// Save ID-token
		gcp.gcpAccessTokenForAuthorizedAccounts = user

		// Trigger Close of Web Server, and 'true' means that a ID-to
		DoneChannel <- true

	})

	router.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	router.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})

	// Initiate channel used to stop server
	DoneChannel = make(chan bool, 1)

	// Initiate http server
	localWebServer := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	// Start Local Web Server as go routine
	go gcp.startLocalWebServer(localWebServer)

	gcp.logger.WithFields(logrus.Fields{
		"ID": "689d42de-3cc0-4237-b1e9-3a6c769f65ea",
	}).Debug("Local webServer Started")

	// Wait for message in channel to stop local web server
	gotIdTokenResult := <-DoneChannel

	// Shutdown local web server
	gcp.stopLocalWebServer(context.Background(), localWebServer)

	// Depending on the outcome of getting a token return different results
	if gotIdTokenResult == true {
		// Success in getting an ID-token
		appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.gcpAccessTokenForAuthorizedAccounts.IDToken)

		return appendedCtx, true, ""
	} else {
		// Didn't get any ID-token
		return nil, false, "Couldn't generate access token"
	}

}

func (gcp *GcpObjectStruct) GenerateGCPAccessTokenForAuthorizedUserPubSub(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Secure that User is initiated
	gcp.initiateUserObjectPubSub()

	// Only create the token if there is none, or it has expired (or 5 minutes before expiration
	timeToCompareTo := time.Now().Add(-time.Minute * 5)
	if !(gcp.GcpAccessTokenForExternalPubSubUserRequests.IDToken == "" || gcp.GcpAccessTokenForExternalPubSubUserRequests.ExpiresAt.Before(timeToCompareTo)) {
		// We already have a ID-token that can be used, so return that
		appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.GcpAccessTokenForExternalPubSubUserRequests.IDToken)

		return appendedCtx, true, ""
	}

	// Need to create a new ID-token

	key := common_config.ApplicationRunTimeUuid // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30                        // 30 days
	isProd := false                             // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		// Use 'Fenix End User Authentication'
		google.New(
			common_config.AuthClientId,
			common_config.AuthClientSecret,
			"http://localhost:3000/auth/google/callback",
			"email",
			"https://www.googleapis.com/auth/cloud-platform"), //"https://www.googleapis.com/auth/pubsub"),
	)

	router := pat.New()

	router.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {

			fmt.Fprintln(res, err)

			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)

		// Save ID-token
		gcp.GcpAccessTokenForExternalPubSubUserRequests = user

		// Trigger Close of Web Server, and 'true' means that a ID-to
		DoneChannel <- true

	})

	router.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	router.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})

	// Initiate channel used to stop server
	DoneChannel = make(chan bool, 1)

	// Initiate http server
	localWebServer := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	// Start Local Web Server as go routine
	go gcp.startLocalWebServer(localWebServer)

	gcp.logger.WithFields(logrus.Fields{
		"ID": "689d42de-3cc0-4237-b1e9-3a6c769f65ea",
	}).Debug("Local webServer Started")

	// Wait for message in channel to stop local web server
	gotIdTokenResult := <-DoneChannel

	// Shutdown local web server
	gcp.stopLocalWebServer(context.Background(), localWebServer)

	// Depending on the outcome of getting a token return different results
	if gotIdTokenResult == true {
		// Success in getting an ID-token
		appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.GcpAccessTokenForExternalPubSubUserRequests.IDToken)

		return appendedCtx, true, ""
	} else {
		// Didn't get any ID-token
		return nil, false, "Couldn't generate access token"
	}

}

// Start and run Local Web Server
func (gcp *GcpObjectStruct) startLocalWebServer(webServer *http.Server) {

	go func() {
		time.Sleep(1 * time.Second)
		err := webbrowser.Open("http://localhost:3000")

		if err != nil {
			gcp.logger.WithFields(logrus.Fields{
				"ID":  "17bc0305-4594-48e1-bb8d-c642579e5e56",
				"err": err,
			}).Fatalf("Couldn't open the web browser")
		}
	}()
	err := webServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		gcp.logger.WithFields(logrus.Fields{
			"ID": "8226cf74-0cdc-4e29-a441-116504b4b333",
		}).Fatalf("Local Web Server failed to listen: %s\n", err)

	}
}

// Close down Local Web Server
func (gcp *GcpObjectStruct) stopLocalWebServer(ctx context.Context, webServer *http.Server) {

	gcp.logger.WithFields(logrus.Fields{
		"ID": "1f4e0354-2a09-4a1d-be61-67ecda781142",
	}).Debug("Trying to stop local web server")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := webServer.Shutdown(ctx)
	if err != nil {
		if err != nil {
			gcp.logger.WithFields(logrus.Fields{
				"ID": "ea06dfab-39b9-4df6-b3ca-7f5f56b3cb91",
			}).Fatalf("Local Web Server Shutdown Failed:%+v", err)

		} else {
			gcp.logger.WithFields(logrus.Fields{
				"gcp": "ea06dfab-39b9-4df6-b3ca-7f5f56b3cb91",
			}).Debug("Local Web Server Exited Properly")
		}

	}

}

// SetLogger
// Set to use the same Logger reference as is used by central part of system
func (gcp *GcpObjectStruct) SetLogger(logger *logrus.Logger) {

	gcp.logger = logger

	return

}

// initiateUserObject
// Set to use the same Logger reference as is used by central part of system
func (gcp *GcpObjectStruct) initiateUserObject() {

	// Only do initiation if it's not done before

	if gcp.gcpAccessTokenForAuthorizedAccounts.UserID == "" {
		gcp.gcpAccessTokenForAuthorizedAccounts = goth.User{}
	}

	return

}

// initiateUserObjectPubSub
// Set to use the same Logger reference as is used by central part of system
func (gcp *GcpObjectStruct) initiateUserObjectPubSub() {

	// Only do initiation if it's not done before

	if gcp.GcpAccessTokenForExternalPubSubUserRequests.UserID == "" {
		gcp.GcpAccessTokenForExternalPubSubUserRequests = goth.User{}
	}

	return

}
