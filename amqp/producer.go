package amqp

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/logger"
	"github.com/streadway/amqp"
)

type ProducerMQ struct {
	am     *amqp.Connection
	config *config.Config
	logger *logger.Logger
	amqpMQ *amqp.Channel
}

func NewProducerMQ(config *config.Config) *ProducerMQ {
	return &ProducerMQ{
		am:     &amqp.Connection{},
		config: config,
		logger: &logger.Logger{},
	}
}

func (pr *ProducerMQ) Connect() (*amqp.Connection, error) {
	connectMQ, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	pr.am = connectMQ
	return connectMQ, nil
}

func (pr *ProducerMQ) ChannelMQ(connect *amqp.Connection, nameQueue string) (amqp.Queue, *amqp.Channel) {
	channel, err := connect.Channel()
	if err != nil {
		pr.logger.ErrorLog("Не могу создать канал - ", err)
	}
	defer channel.Close()

	q, err := channel.QueueDeclare(
		nameQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	pr.amqpMQ = channel
	return q, channel
}

func (pr *ProducerMQ) PublishMQ(amq *amqp.Queue, ch *amqp.Channel) {
	err := ch.Publish(
		"",
		amq.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Test"),
		},
	)
	if err != nil {
		log.Println(err)
	}
	log.Printf("[x] Sent %s", "TEST")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// PublishSender Отправитель
func (pr *ProducerMQ) PublishSender() {
	connect, err := pr.Connect()
	if err != nil {
		pr.logger.ErrorLog("Не могу соединиться с Rabbit", err)
	}
	ch, queue := pr.ChannelMQ(connect, "Add")
	pr.PublishMQ(&ch, queue)
}

// PublishSend добавление записи в очередь
func (pr *ProducerMQ) PublishSend(config *config.Config, nameQueue string, in interface{}) {
	conn, err := amqp.Dial(config.Rabbit.ConnectRabbit)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		nameQueue, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err2 := enc.Encode(in)
	if err != nil {
		log.Println(err2)
	}
	kk := buf.Bytes()

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(kk),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
