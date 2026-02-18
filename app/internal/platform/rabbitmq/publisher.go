package rabbitmq

import (
	"sync"

	"github.com/alfattd/category-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	amqpURL   string
	conn      *amqp.Connection
	channel   *amqp.Channel
	confirmCh chan amqp.Confirmation
	queue     string
	mu        sync.Mutex
}

type categoryEvent struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}

var _ domain.CategoryEventPublisher = (*Publisher)(nil)

func NewPublisher(amqpURL, queueName string) (*Publisher, error) {
	p := &Publisher{
		amqpURL: amqpURL,
		queue:   queueName,
	}

	if err := p.connect(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}
