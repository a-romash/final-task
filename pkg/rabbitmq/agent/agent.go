package agent

import (
	"calculator/internal/agent"
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

func Init(agent *agent.Agent) {
	conn, err := amqp.Dial(config.Config.RabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		100,   // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		_, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		for d := range msgs {
			var expr model.Expression
			err = json.Unmarshal(d.Body, &expr)
			failOnError(err, "Failed to unmarshal JSON")

			log.Println("Expression start solving")

			agent.SolveExpression(&expr)

			log.Printf("Expressions solved, result - %v\n", expr.Result)

			byteExpression, err := json.Marshal(expr)
			failOnError(err, "Failed to marshal JSON")
			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          byteExpression,
				})
			failOnError(err, "Failed to publish a message")
			if err == nil {
				log.Println("Expression sent")
				log.Println(ch.IsClosed())
			}

			d.Ack(false)
		}
		log.Println("end of goroutine")
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
	log.Println("end of func")
}
