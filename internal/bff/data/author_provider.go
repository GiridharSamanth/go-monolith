package data

import (
	"context"
	"strconv"

	authordomain "go-monolith/internal/modules/author/domain"
	authorModuleService "go-monolith/internal/modules/author/service"
)

type AuthorProvider struct {
	authorService *authorModuleService.AuthorService
}

func NewAuthorProvider(as *authorModuleService.AuthorService) *AuthorProvider {
	return &AuthorProvider{
		authorService: as,
	}
}

func (p *AuthorProvider) GetAuthor(ctx context.Context, id string) (*authordomain.Author, error) {
	authorID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}
	return p.authorService.GetByID(ctx, uint(authorID))
}
