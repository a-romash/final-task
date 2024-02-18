package orchestrator

import (
	"calculator/internal/model"
	"calculator/pkg/config"
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SolveExpression(expr *model.Expression) (solvedExpression model.Expression, err error) {
	conn, err := amqp.Dial(config.Config.RabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

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

	corrId := expr.IdExpression

	byteExpression, err := json.Marshal(expr)
	failOnError(err, "Failed to marshal JSON")

	_, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          byteExpression,
		})
	failOnError(err, "Failed to publish a message")

	log.Println("Expression sent")

	for d := range msgs {
		if corrId == d.CorrelationId {
			err = json.Unmarshal(d.Body, &solvedExpression)
			failOnError(err, "Failed to unmarshal JSON")
			log.Println("Got solved expression")
			log.Println(solvedExpression.Result)
			break
		}
	}
	return
}
