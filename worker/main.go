package main

import (
	"worker/message"
)

func main() {

	// Get a RabbitMQ connection & channel
	conn := message.CreateConnection()
	ch := message.CreateChannel(conn)

	// Listen for new jobs
	message.Listen(ch, "work")
}
