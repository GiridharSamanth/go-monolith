package data

import (
	"context"

	authordomain "go-monolith/internal/modules/author/domain"
	storydomain "go-monolith/internal/modules/story/domain"
)

// StoryDataProvider defines the interface for story data operations
type StoryDataProvider interface {
	GetStory(ctx context.Context, storyID string) (*storydomain.Story, error)
}

// AuthorDataProvider defines the interface for author data operations
type AuthorDataProvider interface {
	GetAuthor(ctx context.Context, authorID string) (*authordomain.Author, error)
}
