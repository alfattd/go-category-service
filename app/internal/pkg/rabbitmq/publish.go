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
	return p.publish(ctx, categoryEvent{
		ID:   c.ID,
		Name: c.Name,
		Type: "category_created",
	})
}

func (p *Publisher) PublishCategoryUpdated(
	ctx context.Context,
	c *domain.Category,
) error {
	return p.publish(ctx, categoryEvent{
		ID:   c.ID,
		Name: c.Name,
		Type: "category_updated",
	})
}

func (p *Publisher) PublishCategoryDeleted(
	ctx context.Context,
	id string,
) error {
	return p.publish(ctx, categoryEvent{
		ID:   id,
		Type: "category_deleted",
	})
}

func (p *Publisher) publish(ctx context.Context, event categoryEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	p.mu.Lock()

	if !p.isConnected() {
		if err := p.reconnect(); err != nil {
			p.mu.Unlock()
			return fmt.Errorf("failed to reconnect: %w", err)
		}
	}

	err = p.channel.PublishWithContext(
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
	)

	confirmCh := p.confirmCh

	p.mu.Unlock()

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	select {
	case confirm, ok := <-confirmCh:
		if !ok {
			return fmt.Errorf("confirm channel closed unexpectedly")
		}
		if !confirm.Ack {
			return fmt.Errorf("message not acknowledged by broker (nack received)")
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while waiting for broker confirmation: %w", ctx.Err())
	}
}
