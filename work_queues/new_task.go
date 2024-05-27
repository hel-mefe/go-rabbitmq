package main

import (
  "context"
  "log"
  "os"
  "time"
  "strings"

  amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", msg, err)
  }
}

func bodyFrom(args []string) string {
  var s string
  if (len(args) < 2 || os.Args[1] == "") {
    s = "hello"
  } else {
    s = strings.Join(args[1:], " ")
  }
  return s
}

func main() {
  // creating the connection to rabbitmq open socket
  conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  failOnError(err, "Failed to connect to RabbitMQ")
  defer conn.Close()
  
  // creating a channel therefore we can abstract away the socket connection
  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  defer ch.Close()

  // to send, I have to declare a queue for me to send messages to
  // declaring a queue is idempotent (it will be created only if doesn't exist already)
  q, err := ch.QueueDeclare(
    "hello", // name
    false, // durable
    false, // delete when unused
    false, // exclusive
    false, // no-wait
    nil, // arguments
  )

  failOnError(err, "Failed to declare a queue")

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  
  // here we'll publish a message with the body of "Hello World!"
  body := bodyFrom(os.Args)
  err = ch.PublishWithContext(ctx,
  "", // exchange
  q.Name, // routing key
  false, // mandatory
  false, // immediate
  amqp.Publishing{
    DeliveryMode: amqp.Persistent,
    ContentType: "text/plain",
    Body: []byte(body),
  })

  failOnError(err, "Failed to publish a message")
  log.Printf("[x] Sent %s\n", body)

}
