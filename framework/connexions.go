package framework

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

var AmqpConn *amqp.Connection

func StartConnexions(){
	rabbitMq()
}

func rabbitMq(){

	var err error

	rabbitUrl := os.Getenv("RABBITMQ_URL")
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	rabbitUsername := os.Getenv("RABBITMQ_USERNAME")
	rabbitPassword := os.Getenv("RABBITMQ_PASSWORD")

	AmqpConn, err = amqp.Dial(rabbitUrl+rabbitUsername+":"+rabbitPassword+"@"+rabbitHost+":"+rabbitPort)

	if err != nil {
		log.Fatal(err)
	}
}

func CloseConnexions(){
	AmqpConn.Close()
}