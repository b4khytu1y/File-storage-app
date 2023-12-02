package main

import (
	"database/sql"
	"log"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalln("Не удалось создать консьюмера:", err)
	}

	partitionConsumer, err := consumer.ConsumePartition("userConfirmation", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln("Не удалось создать partition consumer:", err)
	}

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Println("Получено сообщение:", string(msg.Value))

		case err := <-partitionConsumer.Errors():
			log.Println("Ошибка при получении сообщения:", err)
		}
	}
}
