package util

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var consumerName = ""

func Run(amqp_url string, deviceName string, exchange string, exchangeType string, routingKey string, queue string) {

	consumer, err := rabbitmq.NewConsumer(
		amqp_url, amqp.Config{},
		rabbitmq.WithConsumerOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}

	// wait for server to acknowledge the cancel
	noWait := false
	defer consumer.Disconnect()
	defer consumer.StopConsuming(consumerName, noWait)

	err = consumer.StartConsuming(
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))
			var message WireguardMessage
			jsonerr := json.Unmarshal(d.Body, &message)
			if jsonerr != nil {
				log.Printf(jsonerr.Error())
				return rabbitmq.NackRequeue
			}
			seterr := SetPeer(deviceName, message)
			if seterr != nil {
				log.Printf(jsonerr.Error())
				return rabbitmq.NackRequeue
			} else {
				return rabbitmq.Ack
			}
		},
		queue,
		[]string{routingKey},
		rabbitmq.WithConsumeOptionsConcurrency(1),
		rabbitmq.WithConsumeOptionsConsumerAutoAck(false),
		rabbitmq.WithConsumeOptionsBindingExchangeName(exchange),
		rabbitmq.WithConsumeOptionsBindingExchangeKind(exchangeType),
		rabbitmq.WithConsumeOptionsConsumerName(consumerName),
		rabbitmq.WithConsumeOptionsQOSGlobal,
		rabbitmq.WithConsumeOptionsQOSPrefetch(1),
	)
	if err != nil {
		log.Fatal(err)
	}

	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	log.Println("awaiting signal")
	<-done
	log.Println("stopping consumer")
}
