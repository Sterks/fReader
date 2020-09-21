package amqp

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/logger"
	"github.com/streadway/amqp"
)

// ProducerMQ Структура отправителя
type ProducerMQ struct {
	Am     *amqp.Connection
	Config *config.Config
	Logger *logger.Logger
	AmqpMQ *amqp.Channel
}

// ProducerMQNew ...
func ProducerMQNew() *ProducerMQ {
	return &ProducerMQ{
		Am:     &amqp.Connection{},
		Config: &config.Config{},
		Logger: &logger.Logger{},
		AmqpMQ: &amqp.Channel{},
	}
}

// PublishSend добавление записи в очередь
func (pr *ProducerMQ) PublishSend(config *config.Config, info os.FileInfo, nameQueue string, in []byte, id int, region string, fullpath string, file string) {
	conn, err := amqp.Dial(config.Rabbit.ConnectRabbit)
	if err != nil {
		pr.Logger.ErrorLog("Failed to connect to RabbitMQ", err, "amqp.Dial()")
	}

	defer conn.Close()

	ch, err2 := conn.Channel()
	if err2 != nil {
		pr.Logger.ErrorLog("Failed to open a channel", err2, "conn.Channel()")
	}
	defer ch.Close()

	q, err3 := ch.QueueDeclare(
		nameQueue, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err3 != nil {
		pr.Logger.ErrorLog("Failed to declare a queue", err3, "ch.QueueDeclare(")
	}

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

	bodyJSON, err4 := json.Marshal(body)
	if err4 != nil {
		pr.Logger.ErrorLog("Не могу преобразовать в JSON", err4)
	}

	if err5 := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyJSON,
		}); err5 != nil {
		pr.Logger.ErrorLog("Failed to publish a message", err5)
	}
	pr.Logger.InfoLog("[x] Sent -", body.NameFile)
}
