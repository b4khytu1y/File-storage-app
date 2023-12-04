package main

import (
	"log"

	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9055"}, config)
	if err != nil {
		log.Fatalln("Не удалось создать продьюсера:", err)
	}

	confirmationCode := "123"

	msg := &sarama.ProducerMessage{
		Topic: "userConfirmation",
		Value: sarama.StringEncoder(confirmationCode),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Fatalln("Не удалось отправить сообщение:", err)
	}

	log.Println("Сообщение отправлено:", confirmationCode)

	producer.Close()
}
