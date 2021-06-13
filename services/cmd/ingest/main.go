package main

import (
	"encoding/json"
	"os"

	"github.com/dannyhinshaw/go-ingest-meals/pkg/ingest"
	"github.com/dannyhinshaw/go-ingest-meals/pkg/ingest/config"
	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"
)

const serviceName = "ingest"

// IngestApp the data ingestion point, streaming twitter data.
type IngestApp struct {
	Stream *twitter.Stream
	Demux  twitter.Demux
}

// newIngestApp constructor func creates a new IngestApp instance.
func newIngestApp() (*IngestApp, error) {
	configPath := ""
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Create twitter stream client and define params
	client := ingest.NewTwitterClient(cfg.ConsumerKey, cfg.ConsumerSecret, cfg.AccessKey, cfg.AccessSecret)
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
	}, nil
}

func main() {
	log.Infof(`new "%s" instance is deployed`, serviceName)
	app, err := newIngestApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Demux.HandleChan(app.Stream.Messages)
}
