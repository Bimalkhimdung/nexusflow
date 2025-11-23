package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

// ConsumerConfig holds Kafka consumer configuration
type ConsumerConfig struct {
	Brokers       []string
	ConsumerGroup string
	Topics        []string
}

// MessageHandler is a function that processes a Kafka message
type MessageHandler func(ctx context.Context, message *sarama.ConsumerMessage) error

// Consumer wraps Kafka consumer
type Consumer struct {
	client  sarama.ConsumerGroup
	config  ConsumerConfig
	handler MessageHandler
	wg      sync.WaitGroup
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg ConsumerConfig, handler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V3_3_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	client, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		client:  client,
		config:  cfg,
		handler: handler,
	}, nil
}

// Start starts consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	c.wg.Add(1)
	defer c.wg.Done()

	handler := &consumerGroupHandler{
		handler: c.handler,
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := c.client.Consume(ctx, c.config.Topics, handler); err != nil {
				return fmt.Errorf("error from consumer: %w", err)
			}
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	c.wg.Wait()
	return c.client.Close()
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler MessageHandler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.handler(session.Context(), message); err != nil {
			// Log error but continue processing
			// In production, you might want to send to a dead letter queue
			continue
		}
		session.MarkMessage(message, "")
	}
	return nil
}

// EventConsumer provides helper methods for consuming events
type EventConsumer struct {
	*Consumer
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(cfg ConsumerConfig, handler func(ctx context.Context, event Event) error) (*EventConsumer, error) {
	messageHandler := func(ctx context.Context, message *sarama.ConsumerMessage) error {
		var event Event
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}
		return handler(ctx, event)
	}

	consumer, err := NewConsumer(cfg, messageHandler)
	if err != nil {
		return nil, err
	}

	return &EventConsumer{Consumer: consumer}, nil
}
