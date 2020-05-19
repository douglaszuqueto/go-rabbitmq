package rabbit

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// Config config
type Config struct {
	Username    string
	Password    string
	IP          string
	Port        string
	VirtualHost string
}

// Client client
type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewConn NewConn
func NewConn() (*Client, error) {
	config := Config{
		IP:          os.Getenv("RABBITMQ_IP"),
		Port:        os.Getenv("RABBITMQ_PORT"),
		Username:    os.Getenv("RABBITMQ_USERNAME"),
		Password:    os.Getenv("RABBITMQ_PASSWORD"),
		VirtualHost: os.Getenv("RABBITMQ_VIRTUALHOST"),
	}

	return New(config)
}

// New new
func New(cfg Config) (*Client, error) {
	url := makeURL(cfg)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn: conn,
		ch:   ch,
	}

	return client, nil
}

// Channel Channel
func (s *Client) Channel() *amqp.Channel {
	return s.ch
}

// SendMessage SendMessage
func (s *Client) SendMessage(data string) error {
	return s.publish("2fa.email", "*", []byte(data))
}

// ConsumeMessage ConsumeMessage
func (s *Client) ConsumeMessage() (<-chan amqp.Delivery, error) {
	return s.consume("2fa.email.queue")
}

func (s *Client) consume(queue string) (<-chan amqp.Delivery, error) {
	return s.ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func (s *Client) publish(exchange, routingKey string, body []byte) error {
	return s.ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Body: body,
		},
	)
}

// Stop Stop
func (s *Client) Stop() {
	log.Println("Closing rabbit channel")
	err := s.ch.Close()

	if err != nil {
		return
	}

	log.Println("Closing rabbit connection")

	err = s.conn.Close()

	if err != nil {
		return
	}
}

func makeURL(cfg Config) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.IP,
		cfg.Port,
		cfg.VirtualHost,
	)
}
