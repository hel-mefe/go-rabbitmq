package main

import (
  "log"

  amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", msg, err)
  }
}

func main() {

// creating the connection to rabbitMQ
conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
failOnError(err, "Failed to connect to RabbitMQ")
defer conn.Close()

// creating a channel therefore I can talk to rabbitMQ
ch, err := conn.Channel()
failOnError(err, "Failed to open a channel")
defer ch.Close()

// declaring a queue named "hello" if it doesn't exist already
q, err := ch.QueueDeclare(
  "hello", // name
  false,   // durable
  false,   // delete when unused
  false,   // exclusive
  false,   // no-wait
  nil,     // arguments
)
failOnError(err, "Failed to declare a queue")

// consuming messages from the queue
msgs, err := ch.Consume(
  q.Name, // queue
  "",     // consumer
  true,   // auto-ack
  false,  // exclusive
  false,  // no-local
  false,  // no-wait
  nil,    // args
)
failOnError(err, "Failed to register a consumer")


// a forever channel struct that will hold our data received from the queue
var forever chan struct{}

// go routing to receive messages
go func() {
  for d := range msgs {
    log.Printf("Received a message: %s", d.Body)
  }
}()

// line that will be printed before reading any messages (only one time)
log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
<-forever


}
