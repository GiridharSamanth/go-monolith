package author

import (
	"gorm.io/gorm"

	"go-monolith/internal/modules/author/repository"
	"go-monolith/internal/modules/author/service"
	"go-monolith/pkg/logger"
	"go-monolith/pkg/metrics"
)

type Module struct {
	AuthorService *service.AuthorService
}

func NewModule(db *gorm.DB, logger logger.Logger, metrics *metrics.Client) *Module {
	repo := repository.NewAuthorRepository(db)

	return &Module{
		AuthorService: service.NewAuthorService(repo, logger, metrics),
	}
}
