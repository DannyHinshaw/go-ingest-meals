package stream

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"ingest/message"
	"log"
	"os"
)

type AddTweetTask struct {
	UserName  string
	CreatedAt string
	Text      string
}

func createTweetTask(tweet *twitter.Tweet) AddTweetTask {
	tweetText := tweet.Text
	if tweet.ExtendedTweet != nil {
		tweetText = tweet.ExtendedTweet.FullText
	}

	task := AddTweetTask{
		UserName:  tweet.User.Name,
		CreatedAt: tweet.CreatedAt,
		Text:      tweetText,
	}

	return task
}

func newClient() *twitter.Client {

	// Twitter keys
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumeSecret := os.Getenv("CONSUMER_SECRET")
	accessKey := os.Getenv("ACCESS_KEY")
	accessSecret := os.Getenv("ACCESS_SECRET")

	// Config auth
	config := oauth1.NewConfig(consumerKey, consumeSecret)
	token := oauth1.NewToken(accessKey, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Setup streaming client
	client := twitter.NewClient(httpClient)
	return client
}

func Init() {

	// Get a RabbitMQ connection & channel
	conn := message.CreateConnection()
	ch := message.CreateChannel(conn)

	// Create stream client and define params
	client := newClient()
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
		task := createTweetTask(tweet)
		body, err := json.Marshal(task)
		if err != nil {
			log.Fatal(err)
		}

		msg := string(body)
		fmt.Println("task::", task)
		message.SendMessage(ch, "work", msg)
	}

	// Start streaming
	demux.HandleChan(stream.Messages)
}
