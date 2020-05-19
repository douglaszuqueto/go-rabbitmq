package main

import (
	"flag"
	"fmt"

	"github.com/douglaszuqueto/go-rabbitmq/pkg/rabbit"
)

var (
	rabbitSetup = flag.Bool("setup", false, "Setup the rabbit exchange and queue")
)

func main() {
	fmt.Println("RabbitMQ Example")
	fmt.Println()

	flag.Parse()

	rabbitCli, err := rabbit.NewConn()
	if err != nil {
		panic(err)
	}

	// Criação da Exchange, Fila e bind
	if *rabbitSetup {
		setup(rabbitCli)
	}

	//
	// Envio
	//
	rabbitCli.SendMessage("test")

	//
	// Recebimento
	//
	messages, err := rabbitCli.ConsumeMessage()
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		fmt.Printf("Message: %s | size %v\n", msg.Body, len(msg.Body))

		msg.Ack(true)
	}
}

func setup(cli *rabbit.Client) {
	ch := cli.Channel()

	conf := map[string][]string{}

	conf["2fa.email"] = []string{
		"2fa.email.queue",
	}

	for exchange, queues := range conf {
		fmt.Println("Exchange:", exchange)

		// ExchangeDeclare
		err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
		if err != nil {
			fmt.Println("Error on ExchangeDeclare", err)
			continue
		}

		for _, q := range queues {
			fmt.Println("\tQueue conf", q)

			// QueueDeclare
			_, err := ch.QueueDeclare(q, true, false, false, false, nil)
			if err != nil {
				fmt.Println("Error on QueueDeclare", err)
				continue
			}

			// QueueBind
			err = ch.QueueBind(q, "*", exchange, true, nil)
			if err != nil {
				fmt.Println("Error on QueueBind", err)
				continue
			}
		}
	}
}
