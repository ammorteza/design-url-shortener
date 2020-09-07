package rabbit_mq

import (
	job_queue "github.com/ammorteza/clean_architecture/job-queue"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type rabbitMq struct {
	conn 		*amqp.Connection
	ch 			*amqp.Channel
	eventName 	string
}

func New(eventName string) job_queue.JobQueue{
	var conn *amqp.Connection
	var err error
	for {
		log.Println("trying to dial on rabbitMq ...")
		conn, err = amqp.Dial("amqp://admin:admin@172.28.1.110:5672")
		if err == nil{
			break
		}
		time.Sleep(time.Second * 5)
	}

	ch, err := conn.Channel()
	if err != nil{
		log.Fatal(err)
	}

	err  = ch.ExchangeDeclare(eventName , "topic", true, false, false, false, nil)
	if err != nil{
		log.Fatal(err)
	}

	return &rabbitMq{
		conn: conn,
		ch: ch,
		eventName: eventName,
	}
}

func (rmq *rabbitMq)CreateQueue(name string) error{
	_, err := rmq.ch.QueueDeclare(name, false, false, false, false, nil)
	if err != nil{
		return err
	}

	if err = rmq.ch.QueueBind(name, "#", rmq.eventName, false, nil); err != nil{
		return err
	}
	return nil
}

func (rmq *rabbitMq)Publish(identifier string, body []byte) error{
	msg := amqp.Publishing{
		MessageId: identifier,
		Body: body,
	}
	if err := rmq.ch.Publish(rmq.eventName, "random-key", false, false, msg); err != nil{
		return err
	}

	return nil
}

func (rmq *rabbitMq)Consume(qName string, consumer func(id string, body []byte)) error{
	msgs, err := rmq.ch.Consume(qName, "", false, false, false, false, nil)
	if err != nil{
		return err
	}
	go func() {
		for msg := range msgs{
			consumer(msg.MessageId, msg.Body)
			msg.Ack(false)
		}
	}()
	return nil
}