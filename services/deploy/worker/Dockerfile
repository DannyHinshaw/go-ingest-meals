FROM golang:1.16-alpine as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -a -installsuffix cgo -ldflags="-w -s" -o worker cmd/worker/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/worker .

EXPOSE 3666

ENTRYPOINT ["./worker"]
