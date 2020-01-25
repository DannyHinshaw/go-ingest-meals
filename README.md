## go-ingest-meals

Testing architectural patterns with Golang, RabbitMQ, & Docker Compose.
The idea is to use a Twitter stream to create a variable streaming input to test scaling patterns.

We will monitor the following keywords in Tweets:

- breakfast
- second breakfast
- elevenses
- luncheon
- lunch
- afternoon tea
- dinner
- supper

![Merry & Pippin](https://i.imgflip.com/7a9b3.jpg?a438885)

### Roadmap

- [ ] - Add ELK Stack

- [ ] - Scale with Kubernetes/Minikube


### Run

```bash
docker-compose up --build
```
