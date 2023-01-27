package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnErrorx(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@10.125.74.122:5672/")
	failOnErrorx(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrorx(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnErrorx(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	body := bodyFrom(os.Args, strconv.Itoa(r1.Intn(10000)))
	err = ch.PublishWithContext(ctx,
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnErrorx(err, "Failed to publish a message")

	log.Printf(" [x] Message to All Sub --> %s", body)
}

func bodyFrom(args []string, num string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = num
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
