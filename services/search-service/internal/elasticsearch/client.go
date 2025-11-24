package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/nexusflow/nexusflow/pkg/logger"
)

type Client struct {
	es  *elasticsearch.Client
	log *logger.Logger
}

func NewClient(addresses []string, log *logger.Logger) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %w", err)
	}

	// Test connection
	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting elasticsearch info: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	log.Sugar().Infow("Connected to Elasticsearch")

	return &Client{es: es, log: log}, nil
}

// IndexDocument indexes a document
func (c *Client) IndexDocument(ctx context.Context, index, id string, document interface{}) error {
	data, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	res, err := c.es.Index(index, bytes.NewReader(data), c.es.Index.WithDocumentID(id), c.es.Index.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response: %s", res.String())
	}

	return nil
}

// Search performs a search query
func (c *Client) Search(ctx context.Context, index string, query map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("error performing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return result, nil
}

// CreateIndex creates an index with mappings
func (c *Client) CreateIndex(ctx context.Context, index string, mapping map[string]interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		return fmt.Errorf("error encoding mapping: %w", err)
	}

	res, err := c.es.Indices.Create(index, c.es.Indices.Create.WithBody(&buf), c.es.Indices.Create.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		// Ignore "already exists" error
		if strings.Contains(string(body), "resource_already_exists_exception") {
			c.log.Sugar().Infow("Index already exists", "index", index)
			return nil
		}
		return fmt.Errorf("error response: %s", string(body))
	}

	c.log.Sugar().Infow("Created index", "index", index)
	return nil
}

// DeleteDocument deletes a document
func (c *Client) DeleteDocument(ctx context.Context, index, id string) error {
	res, err := c.es.Delete(index, id, c.es.Delete.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response: %s", res.String())
	}

	return nil
}
