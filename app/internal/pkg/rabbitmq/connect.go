package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (p *Publisher) connect() error {
	conn, err := amqp.Dial(p.amqpURL)
	if err != nil {
		return fmt.Errorf("failed to dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to set confirm mode: %w", err)
	}

	_, err = ch.QueueDeclare(
		p.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	p.conn = conn
	p.channel = ch

	p.confirmCh = ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	return nil
}

func (p *Publisher) reconnect() error {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil && !p.conn.IsClosed() {
		_ = p.conn.Close()
	}

	return p.connect()
}

func (p *Publisher) isConnected() bool {
	return p.conn != nil && !p.conn.IsClosed() &&
		p.channel != nil && !p.channel.IsClosed()
}
