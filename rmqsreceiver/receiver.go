package rmqsreceiver

import (
        "encoding/json"
        "fmt"
        "os"
        "os/signal"
        "syscall"
        "strconv"
        "net/http"

        "cust-mobile-eu/database/model"
        "cust-mobile-eu/handler"

        amqp "github.com/rabbitmq/amqp091-go"
        log "github.com/sirupsen/logrus"
        "github.com/prometheus/client_golang/prometheus"
)

type Receiver struct {
	//connection *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

/*
type Product struct {
        Id          int     `json:"id"`
        Name        string  `json:"name"`
        Description string  `json:"description"`
        Price       float32 `json:"price"` //UnitPrice
        Quantity    int     `json:"quantity"`
        Discount    int     `json:"discount"` //MaxDiscountPercent
        Country     string  `json:"country"`
        Region      string  `json:"region"`
}
*/


var (
        IdAccessCounter = prometheus.NewCounterVec(
                prometheus.CounterOpts{
                        Name: "product_id_access_total_eu",
                        Help: "Total number of times product ids are accessed",
                },
                []string{"id"},
        )
)

func init() {
        //prometheus.MustRegister(IdAccessCounter)
}



func Connect(queueName string) *Receiver {
	//connection, err := amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
	connection, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672")
	if err != nil {
		log.Panicf("failed to connect RabbitMQ : %+v ", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Panicf("failed to get channel : %+v ", err)
	}

	//queue, err := channel.QueueDeclare("Produced", false, false, false, false, nil)
	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Panicf("failed to get queue : %+v ", err)
	}

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		<-interrupt

		connection.Close()
		channel.Close()
	}()

	return &Receiver{channel: channel, queue: &queue}
	//return &Receiver{connection: connection, channel: channel, queue: &queue}
}

func (r *Receiver) ReceiveMessage() error {
	messages, err := r.channel.Consume(r.queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	var productJson model.Product
	
	for message := range messages {
		fmt.Println(string(message.Body))
		err := json.Unmarshal([]byte(message.Body), &productJson)
                if err != nil {
                        log.Printf("failed to parse json consumer message %+v ", err)
                }
                fmt.Println("Processed ", productJson)
                IdAccessCounter.WithLabelValues(strconv.Itoa(productJson.Id)).Inc()
                
                p := handler.Product{}
                p.Build(nil)

                status := p.UpsertProduct(&productJson)
                if status == http.StatusOK {
                        log.Printf("upsert worked fine")
                } else {
                        log.Printf("issue with upsert")
                }
	}
	return nil
}

/*
func (r *Receiver) Close() {
	r.channel.Close()
	r.connection.Close()
}
*/
