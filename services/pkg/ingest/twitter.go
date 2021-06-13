package ingest

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func NewTwitterClient(consumerKey, consumerSecret, accessKey, accessSecret string) *twitter.Client {

	// Config auth
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessKey, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Setup streaming client
	client := twitter.NewClient(httpClient)
	return client
}
