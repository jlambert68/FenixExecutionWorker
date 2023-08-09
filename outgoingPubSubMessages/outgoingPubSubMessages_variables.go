package outgoingPubSubMessages

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type PubSubMessagesToConnectorObjectStruct struct {
	Logger         *logrus.Logger
	gcpAccessToken *oauth2.Token
}

var PubSubMessagesToConnectorObject PubSubMessagesToConnectorObjectStruct
