package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

// ProducerConfig holds Kafka producer configuration
type ProducerConfig struct {
	Brokers []string
}

// Producer wraps Kafka producer
type Producer struct {
	producer sarama.SyncProducer
	config   ProducerConfig
}

// Event represents a generic event
type Event struct {
	ID             string                 `json:"id"`
	Type           string                 `json:"type"`
	OrganizationID string                 `json:"organization_id"`
	ProjectID      string                 `json:"project_id,omitempty"`
	UserID         string                 `json:"user_id"`
	Timestamp      time.Time              `json:"timestamp"`
	Payload        map[string]interface{} `json:"payload"`
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg ProducerConfig) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Producer{
		producer: producer,
		config:   cfg,
	}, nil
}

// PublishEvent publishes an event to a Kafka topic
func (p *Producer) PublishEvent(topic string, event Event) error {
	// Set event ID and timestamp if not set
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Marshal event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.ID),
		Value: sarama.ByteEncoder(data),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event_type"),
				Value: []byte(event.Type),
			},
			{
				Key:   []byte("organization_id"),
				Value: []byte(event.OrganizationID),
			},
		},
	}

	// Send message
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	_ = partition // Suppress unused variable warning
	_ = offset

	return nil
}

// PublishJSON publishes arbitrary JSON data to a topic
func (p *Producer) PublishJSON(topic string, key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(jsonData),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.producer.Close()
}

// Common event types
const (
	EventTypeIssueCreated   = "issue.created"
	EventTypeIssueUpdated   = "issue.updated"
	EventTypeIssueDeleted   = "issue.deleted"
	EventTypeIssueAssigned  = "issue.assigned"
	EventTypeIssueCompleted = "issue.completed"

	EventTypeProjectCreated = "project.created"
	EventTypeProjectUpdated = "project.updated"
	EventTypeProjectDeleted = "project.deleted"

	EventTypeUserCreated = "user.created"
	EventTypeUserUpdated = "user.updated"

	EventTypeOrgCreated = "org.created"
	EventTypeOrgUpdated = "org.updated"

	EventTypeCommentCreated = "comment.created"
	EventTypeCommentUpdated = "comment.updated"
	EventTypeCommentDeleted = "comment.deleted"

	EventTypeWorkflowTransition = "workflow.transition"
)

// Common topics
const (
	TopicIssueEvents    = "nexusflow.issues"
	TopicProjectEvents  = "nexusflow.projects"
	TopicUserEvents     = "nexusflow.users"
	TopicOrgEvents      = "nexusflow.organizations"
	TopicCommentEvents  = "nexusflow.comments"
	TopicWorkflowEvents = "nexusflow.workflows"
	TopicAuditLogs      = "nexusflow.audit"
	TopicNotifications  = "nexusflow.notifications"
)
