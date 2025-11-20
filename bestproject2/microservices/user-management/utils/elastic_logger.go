package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticLogger struct {
	client      *elasticsearch.Client
	indexPrefix string
	serviceName string
	hostname    string
}

type LogEntry struct {
	Timestamp   string                 `json:"@timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Service     string                 `json:"service"`
	Environment string                 `json:"environment"`
	Hostname    string                 `json:"hostname"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

func NewElasticLogger(esURL, serviceName string) (*ElasticLogger, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	// Test connection
	_, err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to Elasticsearch: %w", err)
	}

	hostname, _ := os.Hostname()

	return &ElasticLogger{
		client:      client,
		indexPrefix: "user-management-logs",
		serviceName: serviceName,
		hostname:    hostname,
	}, nil
}

func (el *ElasticLogger) Log(level, message string, fields map[string]interface{}) error {
	entry := LogEntry{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Level:       level,
		Message:     message,
		Service:     el.serviceName,
		Environment: "production",
		Hostname:    el.hostname,
		Fields:      fields,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("error marshaling log entry: %w", err)
	}

	// Create index name with date (e.g., user-management-logs-2025.11.18)
	indexName := fmt.Sprintf("%s-%s", el.indexPrefix, time.Now().Format("2006.01.02"))

	req := esapi.IndexRequest{
		Index: indexName,
		Body:  bytes.NewReader(data),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := req.Do(ctx, el.client)
	if err != nil {
		return fmt.Errorf("error indexing log: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	return nil
}

func (el *ElasticLogger) Info(message string, fields map[string]interface{}) {
	_ = el.Log("INFO", message, fields)
}

func (el *ElasticLogger) Warning(message string, fields map[string]interface{}) {
	_ = el.Log("WARNING", message, fields)
}

func (el *ElasticLogger) Error(message string, fields map[string]interface{}) {
	_ = el.Log("ERROR", message, fields)
}

func (el *ElasticLogger) Debug(message string, fields map[string]interface{}) {
	_ = el.Log("DEBUG", message, fields)
}
