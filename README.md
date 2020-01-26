# go-ingest-meals

Testing architectural patterns with Golang, RabbitMQ, & ELK with Docker Compose/Kubernetes.

The idea is to use a Twitter stream to create a variable streaming input to test scaling patterns.


## Stream
We will monitor for the following keywords in Tweets:

- breakfast
- second breakfast
- elevenses
- luncheon
- lunch
- afternoon tea
- dinner
- supper


![Merry & Pippin](https://i.imgflip.com/7a9b3.jpg?a438885)
&nbsp;


## Roadmap

- [x] Twitter meals stream ingest

- [x] Twitter meals workers

- [x] ELK Stack

- [ ] Scale with Kubernetes


## Setup

1. Get a [Twitter developer account](https://developer.twitter.com/).

2. You will need to place a `.env` file in the root directory containing:

```bash
# Or whichever version you prefer
ELK_VERSION=7.5.1

# Twitter Keys
CONSUMER_KEY=<TWITTER_CONSUMER_KEY>
CONSUMER_SECRET=<TWITTER_CONSUMER_SECRET>
ACCESS_KEY=<TWITTER_ACCESS_KEY>
ACCESS_SECRET=<TWITTER_ACCESS_SECRET>
```

## Run

```bash
docker-compose up --build
```


## Kibana

Visit `http://localhost:5601/` in your browser. Login with:
- User: `elastic`
- Pass: `changeme`
