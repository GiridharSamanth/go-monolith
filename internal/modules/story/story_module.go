package story

import (
	"gorm.io/gorm"

	"go-monolith/internal/modules/story/repository"
	"go-monolith/internal/modules/story/service"
	"go-monolith/pkg/logger"
)

type Module struct {
	StoryService *service.StoryService
}

func NewModule(db *gorm.DB, logger logger.Logger) *Module {
	repo := repository.NewStoryRepository(db)

	return &Module{
		StoryService: service.NewStoryService(repo, logger),
	}
}
