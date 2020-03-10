package amqp

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/logger"
	"github.com/streadway/amqp"
)

type ProducerMQ struct {
	am     *amqp.Connection
	config *config.Config
	logger *logger.Logger
	amqpMQ *amqp.Channel
}

func (pr *ProducerMQ) Connect() (*amqp.Connection, error) {
	connectMQ, err := amqp.Dial("rabbit://guest:guest@127.0.0.1:5672/")
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
		log.Printf("Не могу создать канал - ", err)
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// PublishSend добавление записи в очередь
func (pr *ProducerMQ) PublishSend(config *config.Config, info os.FileInfo, nameQueue string, in []byte, id int, region string, fullpath string, file string) {
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

	// var buf bytes.Buffer
	// enc := gob.NewEncoder(&buf)
	// err2 := enc.Encode(in)
	// if err != nil {
	// 	log.Println(err2)
	// }
	// kk := buf.Bytes()

	// body := "Hello World!"

	type InformationFile struct {
		FileID   int
		NameFile string
		SizeFile int64
		DateMode time.Time
		Fullpath string
		Region   string
		FileZip  []byte
		TypeFile string
	}

	body := &InformationFile{
		FileID:   id,
		DateMode: info.ModTime(),
		NameFile: info.Name(),
		FileZip:  in,
		SizeFile: info.Size(),
		Fullpath: fullpath,
		Region:   region,
		TypeFile: file,
	}

	bodyJSON, err3 := json.Marshal(body)
	if err3 != nil {
		log.Println(err3)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyJSON,
		})
	log.Printf(" [x] Sent %s", body.NameFile)
	failOnError(err, "Failed to publish a message")
}
