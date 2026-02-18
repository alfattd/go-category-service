package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alfattd/category-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (p *Publisher) PublishCategoryCreated(
	ctx context.Context,
	c *domain.Category,
) error {
	event := categoryEvent{
		ID:   c.ID,
		Name: c.Name,
		Type: "category_created",
	}

	return p.publish(ctx, event)
}

func (p *Publisher) PublishCategoryUpdated(
	ctx context.Context,
	c *domain.Category,
) error {
	event := categoryEvent{
		ID:   c.ID,
		Name: c.Name,
		Type: "category_updated",
	}

	return p.publish(ctx, event)
}

func (p *Publisher) PublishCategoryDeleted(
	ctx context.Context,
	id string,
) error {
	event := categoryEvent{
		ID:   id,
		Type: "category_deleted",
	}

	return p.publish(ctx, event)
}

func (p *Publisher) publish(ctx context.Context, event categoryEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.channel.IsClosed() {
		if err := p.reconnect(); err != nil {
			return fmt.Errorf("failed to reconnect: %w", err)
		}
	}

	if err := p.channel.PublishWithContext(
		ctx,
		"",
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	select {
	case confirm, ok := <-p.confirmCh:
		if !ok {
			return fmt.Errorf("confirm channel closed")
		}
		if !confirm.Ack {
			return fmt.Errorf("message not acknowledged by broker")
		}
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while waiting for confirmation: %w", ctx.Err())
	}

	return nil
}
