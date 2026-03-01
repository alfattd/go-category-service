package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/alfattd/category-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	publishTimeout = 5 * time.Second
	maxRetries     = 3
	baseDelay      = 200 * time.Millisecond
)

func (p *Publisher) PublishCategoryCreated(
	ctx context.Context,
	c *domain.Category,
) error {
	return p.publishWithRetry(ctx, categoryEvent{
		ID:   c.ID,
		Name: c.Name,
		Type: "category_created",
	})
}

func (p *Publisher) PublishCategoryUpdated(
	ctx context.Context,
	c *domain.Category,
) error {
	return p.publishWithRetry(ctx, categoryEvent{
		ID:   c.ID,
		Name: c.Name,
		Type: "category_updated",
	})
}

func (p *Publisher) PublishCategoryDeleted(
	ctx context.Context,
	id string,
) error {
	return p.publishWithRetry(ctx, categoryEvent{
		ID:   id,
		Type: "category_deleted",
	})
}

func (p *Publisher) publishWithRetry(ctx context.Context, event categoryEvent) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(math.Pow(2, float64(attempt-1))) * baseDelay
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry backoff: %w", ctx.Err())
			}
		}

		if err := p.publish(ctx, event); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to publish after %d attempts: %w", maxRetries+1, lastErr)
}

func (p *Publisher) publish(ctx context.Context, event categoryEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Use a dedicated timeout for publish, independent of the HTTP request context.
	// This prevents a client disconnect from aborting an in-flight publish.
	publishCtx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	p.mu.Lock()

	if !p.isConnected() {
		if err := p.reconnect(); err != nil {
			p.mu.Unlock()
			return fmt.Errorf("failed to reconnect: %w", err)
		}
	}

	err = p.channel.PublishWithContext(
		publishCtx,
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
	case <-publishCtx.Done():
		return fmt.Errorf("timed out waiting for broker confirmation: %w", publishCtx.Err())
	}
}
