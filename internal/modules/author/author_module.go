package author

import (
	"gorm.io/gorm"

	"go-monolith/internal/modules/author/repository"
	"go-monolith/internal/modules/author/service"
)

type Module struct {
	AuthorService *service.AuthorService
}

func NewModule(db *gorm.DB) *Module {
	repo := repository.NewAuthorRepository(db)

	return &Module{
		AuthorService: service.NewAuthorService(repo),
	}
}
