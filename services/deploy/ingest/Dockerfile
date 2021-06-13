FROM golang:1.16-alpine as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -a -installsuffix cgo -ldflags="-w -s" -o ingest cmd/ingest/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/ingest .

EXPOSE 3650

ENTRYPOINT ["./ingest"]
