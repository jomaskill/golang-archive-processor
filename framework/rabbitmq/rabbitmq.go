package rabbitmq

import (
	"app/framework"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

func Consume(queue string)  ([]map[string]string, error){

	conn := framework.AmqpConn

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	var messages []map[string]string

	go func() {
		for d := range msgs {
			data := make(map[string]string)
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				fmt.Println(err)
			}
			messages = append(messages, data)
		}
	}()
	ch.Cancel("", false)

	/*
	select {
		case <-time.After(time.Second * 2):
			return messages, nil
	}*/
	return messages, nil
}

// É apenas uma função usada para test
func Publish(queue string, msg string)  error{

	conn := framework.AmqpConn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(msg),
		})

	return nil
}