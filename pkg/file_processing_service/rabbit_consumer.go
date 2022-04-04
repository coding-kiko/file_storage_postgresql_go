package file_processing_service

import (
	"errors"
	"log"

	"github.com/streadway/amqp"
)

func NewQueueConsumer(conn amqp.Connection, service FileProcessingService) error {
	ch, err := conn.Channel()
	if err != nil {
		return errors.New("failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"new-file-id-queue", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return errors.New("failed to declare a queue")
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return errors.New("failed to register a consumer")
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s\n", d.Body)
			if err = service.ResizeImage(string(d.Body)); err != nil {
				log.Fatalln(err)
			}
		}
	}()
	<-forever
	return nil
}
