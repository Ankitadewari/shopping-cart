package kafka

import (
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type service struct {
	kafkaClient *kafka.Conn
}

func New() *service {
	host := os.Getenv("KAFKA_HOST")
	port := os.Getenv("KAFKA_PORT")

	kafkaEndpoint := fmt.Sprintf("%s:%s", host, port)

	conn, err := kafka.Dial("tcp", kafkaEndpoint)
	if err != nil {
		log.Println("kafka not able to connect", "error", err.Error())

		return nil
	}

	log.Println("successfully connected to kafka")

	return &service{
		kafkaClient: conn,
	}
}
