package container

import (
	"go-monolith/internal/app/config"
	"go-monolith/internal/bff/data"
	"go-monolith/internal/bff/handler"
	"go-monolith/internal/bff/service"
	"go-monolith/internal/modules/author"
	"go-monolith/internal/modules/story"
	"go-monolith/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	Config       *config.Config
	Logger       *logger.Logger
	DB           *gorm.DB
	StoryModule  *story.Module
	AuthorModule *author.Module
	StoryRepo    *data.StoryProvider
	AuthorRepo   *data.AuthorProvider
	StoryService *service.StoryService
	Handlers     *handler.Handlers
}

// NewContainer creates a new dependency container
func NewContainer(db *gorm.DB) *Container {
	// Initialize config
	cfg, err := config.NewConfig()
	if err != nil {
		// Use default logger to log error and create default config
		logger.Default().Error(nil, "Failed to initialize config, using defaults", zap.Error(err))
		cfg = &config.Config{
			Environment: "development",
			Logger:      logger.Default(),
			Server: config.ServerConfig{
				Port:           "8080",
				EnableHTTPLogs: true,
			},
		}
	}

	// Initialize modules
	storyModule := story.NewModule(db)
	authorModule := author.NewModule(db)

	// Initialize repositories
	storyRepo := data.NewStoryProvider(storyModule.StoryService)
	authorRepo := data.NewAuthorProvider(authorModule.AuthorService)

	// Initialize BFF service
	storyService := service.NewStoryService(storyRepo, authorRepo, cfg.Logger)

	// Initialize handlers
	handlers := handler.NewHandlers(storyService)

	return &Container{
		Config:       cfg,
		Logger:       cfg.Logger,
		DB:           db,
		StoryModule:  storyModule,
		AuthorModule: authorModule,
		StoryRepo:    storyRepo,
		AuthorRepo:   authorRepo,
		StoryService: storyService,
		Handlers:     handlers,
	}
}
