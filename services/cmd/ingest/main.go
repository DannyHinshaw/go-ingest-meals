package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dannyhinshaw/go-ingest-meals/pkg/ingest"
	"github.com/dghubble/go-twitter/twitter"
)

// IngestApp the data ingestion point, streaming twitter data.
type IngestApp struct {
	Stream *twitter.Stream
	Demux  twitter.Demux
}

// newIngestApp constructor func creates a new IngestApp instance.
func newIngestApp() *IngestApp {

	// Create twitter stream client and define params
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessKey := os.Getenv("ACCESS_KEY")
	accessSecret := os.Getenv("ACCESS_SECRET")
	client := ingest.NewTwitterClient(consumerKey, consumerSecret, accessKey, accessSecret)
	params := &twitter.StreamFilterParams{
		StallWarnings: twitter.Bool(true),
		Track: []string{
			"breakfast",
			"second breakfast",
			"elevenses",
			"luncheon",
			"lunch",
			"afternoon tea",
			"dinner",
			"supper",
		},
	}

	// Init stream
	stream, err := client.Streams.Filter(params)
	if err != nil {
		log.Fatal(err)
	}

	// Stream data handling
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		task, err := json.Marshal(tweet)
		if err != nil {
			log.Fatal(err)
		}

		msg := string(task)
		log.Println("TODO: Something with msg::", msg)
	}

	return &IngestApp{
		Stream: stream,
		Demux:  demux,
	}
}

func main() {
	app := newIngestApp()
	app.Demux.HandleChan(app.Stream.Messages)
}
