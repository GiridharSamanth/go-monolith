package metrics

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

// Config holds the configuration for metrics
type Config struct {
	Host     string
	Port     int
	Prefix   string
	Sampling float64
	Enabled  bool // Whether metrics are enabled
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     8125,
		Prefix:   "go-monolith",
		Sampling: 1.0,
		Enabled:  true,
	}
}

// Client wraps the statsd client
type Client struct {
	client *statsd.Client
	config *Config
}

// NoOpClient implements a no-op metrics client
type NoOpClient struct{}

// NewClient creates a new metrics client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if !config.Enabled {
		return &Client{
			client: nil,
			config: config,
		}, nil
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := statsd.New(addr,
		statsd.WithNamespace(config.Prefix),
		statsd.WithTags([]string{fmt.Sprintf("sampling:%f", config.Sampling)}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create statsd client: %w", err)
	}

	return &Client{
		client: client,
		config: config,
	}, nil
}

// Close closes the metrics client
func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

// Generic metrics methods
func (c *Client) IncrementCounter(name string, tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Incr(name, tags, c.config.Sampling)
}

func (c *Client) RecordTiming(name string, duration time.Duration, tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Timing(name, duration, tags, c.config.Sampling)
}

// HTTP metrics
func (c *Client) RecordHTTPRequestDuration(duration time.Duration, tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Timing("http.request_duration", duration, tags, c.config.Sampling)
}

func (c *Client) IncrementFailedRequests(tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Incr("http.failed_requests", tags, c.config.Sampling)
}

func (c *Client) IncrementErrors(tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Incr("http.errors", tags, c.config.Sampling)
}

// Database metrics
func (c *Client) RecordSQLQueryDuration(duration time.Duration, tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Timing("db.query_duration", duration, tags, c.config.Sampling)
}

func (c *Client) SetOpenConnections(count int64, tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Gauge("db.open_connections", float64(count), tags, c.config.Sampling)
}

func (c *Client) SetInUseConnections(count int64, tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Gauge("db.inuse_connections", float64(count), tags, c.config.Sampling)
}

// Cache metrics
func (c *Client) IncrementCacheHits(tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Incr("cache.hits", tags, c.config.Sampling)
}

func (c *Client) IncrementCacheMisses(tags []string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Incr("cache.misses", tags, c.config.Sampling)
}

// API metrics
func (c *Client) IncrementAPICalls(endpoint string, method string, tags []string) error {
	if c.client == nil {
		return nil
	}
	if tags == nil {
		tags = make([]string, 0)
	}
	tags = append(tags, fmt.Sprintf("endpoint:%s", endpoint), fmt.Sprintf("method:%s", method))
	return c.client.Incr("api.calls", tags, c.config.Sampling)
}
