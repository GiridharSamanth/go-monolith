package container

import (
	"go-monolith/internal/app/config"
	"go-monolith/internal/bff/data"
	"go-monolith/internal/bff/handler"
	"go-monolith/internal/bff/service"
	"go-monolith/internal/modules/author"
	"go-monolith/internal/modules/story"
	"go-monolith/pkg/logger"
	"go-monolith/pkg/metrics"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	Config       *config.Config
	Logger       logger.Logger
	Metrics      *metrics.Client
	DB           *gorm.DB
	StoryModule  *story.Module
	AuthorModule *author.Module
	StoryRepo    *data.StoryProvider
	AuthorRepo   *data.AuthorProvider
	StoryService *service.StoryService
	Handlers     *handler.Handlers
}

// NewContainer creates a new dependency container
func NewContainer() *Container {
	// Initialize config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize logger
	logger, err := logger.Initialize(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize metrics
	metricsClient, err := metrics.NewClient(&cfg.Metrics)
	if err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}

	// Initialize GORM with metrics plugin if enabled
	db, err := gorm.Open(mysql.Open(cfg.DB.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Add metrics plugin if enabled
	if cfg.Metrics.Enabled {
		if err := db.Use(metrics.NewGormMetricsPlugin(metricsClient)); err != nil {
			log.Fatalf("Failed to add metrics plugin: %v", err)
		}
	}

	// Initialize modules
	storyModule := story.NewModule(db, logger, metricsClient)
	authorModule := author.NewModule(db, logger, metricsClient)

	// Initialize repositories
	storyRepo := data.NewStoryProvider(storyModule.StoryService)
	authorRepo := data.NewAuthorProvider(authorModule.AuthorService)

	// Initialize BFF service
	storyService := service.NewStoryService(storyRepo, authorRepo, logger, metricsClient)

	// Initialize handlers
	handlers := handler.NewHandlers(storyService)

	return &Container{
		Config:       cfg,
		Logger:       logger,
		Metrics:      metricsClient,
		DB:           db,
		StoryModule:  storyModule,
		AuthorModule: authorModule,
		StoryRepo:    storyRepo,
		AuthorRepo:   authorRepo,
		StoryService: storyService,
		Handlers:     handlers,
	}
}
