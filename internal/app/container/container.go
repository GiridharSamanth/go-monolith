package container

import (
	"go-monolith/internal/app/config"
	"go-monolith/internal/bff/data"
	"go-monolith/internal/bff/handler"
	"go-monolith/internal/bff/service"
	"go-monolith/internal/modules/author"
	"go-monolith/internal/modules/story"
	"go-monolith/pkg/logger"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	Config       *config.Config
	Logger       logger.Logger // Interface type for dependency inversion and easier testing/mocking
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

	// Initialize GORM
	db, err := gorm.Open(mysql.Open(cfg.DB.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize logger
	log, err := logger.Initialize(cfg.Logger)
	if err != nil {
		// Fallback to default logger
		log = logger.Default()
	}

	// Initialize modules
	storyModule := story.NewModule(db, log)
	authorModule := author.NewModule(db, log)

	// Initialize repositories
	storyRepo := data.NewStoryProvider(storyModule.StoryService)
	authorRepo := data.NewAuthorProvider(authorModule.AuthorService)

	// Initialize BFF service
	storyService := service.NewStoryService(storyRepo, authorRepo, log)

	// Initialize handlers
	handlers := handler.NewHandlers(storyService)

	return &Container{
		Config:       cfg,
		Logger:       log,
		DB:           db,
		StoryModule:  storyModule,
		AuthorModule: authorModule,
		StoryRepo:    storyRepo,
		AuthorRepo:   authorRepo,
		StoryService: storyService,
		Handlers:     handlers,
	}
}
