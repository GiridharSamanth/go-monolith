package metrics

import (
	"time"

	"gorm.io/gorm"
)

// GormMetricsPlugin implements gorm.Plugin interface for metrics tracking
type GormMetricsPlugin struct {
	client *Client
}

// NewGormMetricsPlugin creates a new GORM metrics plugin
func NewGormMetricsPlugin(client *Client) *GormMetricsPlugin {
	return &GormMetricsPlugin{
		client: client,
	}
}

// Name returns the plugin name
func (p *GormMetricsPlugin) Name() string {
	return "gorm:metrics"
}

// Initialize implements gorm.Plugin interface
func (p *GormMetricsPlugin) Initialize(db *gorm.DB) error {
	// Track query execution time
	db.Callback().Query().Before("gorm:query").Register("metrics:before_query", func(db *gorm.DB) {
		db.InstanceSet("metrics:start_time", time.Now())
	})

	db.Callback().Query().After("gorm:query").Register("metrics:after_query", func(db *gorm.DB) {
		startTime, ok := db.InstanceGet("metrics:start_time")
		if !ok {
			return
		}

		duration := time.Since(startTime.(time.Time))
		tags := []string{
			"query:" + db.Statement.SQL.String(),
		}

		if err := db.Error; err != nil {
			p.client.IncrementErrors(tags)
		} else {
			p.client.RecordSQLQueryDuration(duration, tags)
		}
	})

	// Track connection pool stats
	db.Callback().Create().Before("gorm:create").Register("metrics:before_create", func(db *gorm.DB) {
		sqlDB, err := db.DB()
		if err != nil {
			return
		}

		stats := sqlDB.Stats()
		p.client.SetOpenConnections(int64(stats.OpenConnections), nil)
		p.client.SetInUseConnections(int64(stats.InUse), nil)
	})

	return nil
}
