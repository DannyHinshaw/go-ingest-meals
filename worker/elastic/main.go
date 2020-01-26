package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	es "github.com/olivere/elastic/v7"
	"log"
	"strings"
	"time"
)

type TweetElastic struct {
	User     string           `json:"user"`
	Meals    []string         `json:"meals"`
	Message  string           `json:"message"`
	Created  time.Time        `json:"created,omitempty"`
	Location string           `json:"location,omitempty"`
	Suggest  *es.SuggestField `json:"suggest_field,omitempty"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"properties":{
			"user":{
				"type":"keyword"
			},
			"meals":{
				"type":"keyword"
			},
			"message":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"created":{
				"type":"date"
			},
			"location":{
				"type":"keyword"
			},
			"suggest_field":{
				"type":"completion"
			}
		}
	}
}`

var meals = []string{
	"breakfast",
	"second breakfast",
	"elevenses",
	"luncheon",
	"lunch",
	"afternoon tea",
	"dinner",
	"supper",
}

// elasticAdd - Add a new tweet to elastic search.
func elasticAdd(client *es.Client, meals []string, tweet twitter.Tweet) {

	// Use the IndexExists service to check if a specified meals exists.
	ctx := context.Background()
	index := "tweet"
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	// If doesn't exist yet, create it
	if !exists {
		createIndex, err := client.CreateIndex(index).BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			fmt.Println("createIndex.Acknowledged returned false")
		}
	}

	// Convert tweet time format to golang builtin
	layout := "Wed Oct 10 20:19:24 +0000 2018"
	tweetTime, err := time.Parse(layout, tweet.CreatedAt)
	if err != nil {
		// Handle error
		panic(err)
	}

	// Index a tweet (using JSON serialization)
	newTweet := TweetElastic{
		Message:  tweet.Text,
		User:     tweet.User.Name,
		Location: tweet.User.Location,
		Created:  tweetTime,
		Meals:    meals,
	}
	put, err := client.Index().
		Index(index).
		BodyJson(newTweet).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to meals %s, type %s\n", put.Id, put.Index, put.Type)

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index(index).Do(ctx)
	if err != nil {
		panic(err)
	}
}

// deserializeTweet - Util function to deserialize tweet into object form.
func deserializeTweet(body []byte) (twitter.Tweet, error) {
	tweet := twitter.Tweet{}
	err := json.Unmarshal(body, &tweet)

	return tweet, err
}

// isReTweet - Tests whether current tweet is original or a re-tweet.
func isReTweet(text string) bool {
	return strings.HasPrefix(text, "RT @")
}

// detectKeywords - Checks a tweet message body for one of out keywords.
func detectKeywords(msg string) []string {

	// Filter out re-tweets
	var indexes []string
	if isReTweet(msg) {
		return indexes
	}

	fmt.Println("TEXT::", msg)
	msgLower := strings.ToLower(msg)
	for _, meal := range meals {
		if strings.Contains(msgLower, meal) {
			indexes = append(indexes, meal)
		}
	}

	return indexes
}

// ProcessTweet - Start data processing on newly received tweet from queue.
func ProcessTweet(client *es.Client, body []byte) {
	tweet, err := deserializeTweet(body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if tweet.ExtendedTweet != nil {
		tweet.Text = tweet.ExtendedTweet.FullText
	}

	indexes := detectKeywords(tweet.Text)
	if len(indexes) > 0 {
		fmt.Println("INDEXES::", indexes)
		elasticAdd(client, indexes, tweet)
	}
}

// NewClient - Creates and returns a new elastic-search client with current context.
func NewClient() *es.Client {

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default ElasticSearch installation
	url := "http://elasticsearch:9200"
	urlFn := es.SetURL(url)
	authFn := es.SetBasicAuth("elastic", "changeme")
	client, err := es.NewClient(urlFn, authFn)
	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the ElasticSearch server to get e.g. the version number
	info, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(url)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	return client
}
