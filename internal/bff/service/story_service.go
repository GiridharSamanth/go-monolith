package service

import (
	"context"
	"strconv"

	data "go-monolith/internal/bff/data"
	authordomain "go-monolith/internal/modules/author/domain"
	storydomain "go-monolith/internal/modules/story/domain"
	"go-monolith/pkg/logger"
)

type StoryService struct {
	storyProvider  data.StoryDataProvider
	authorProvider data.AuthorDataProvider
	Logger         logger.Logger
}

var storyService *StoryService

func NewStoryService(sp data.StoryDataProvider, ap data.AuthorDataProvider, log logger.Logger) *StoryService {
	if storyService == nil {
		storyService = &StoryService{
			storyProvider:  sp,
			authorProvider: ap,
			Logger:         log,
		}
	}
	return storyService
}

// GetStoryService returns the singleton instance of StoryService
func GetStoryService() *StoryService {
	return storyService
}

// GetStoryDisplayDetails retrieves a story and its author details
// Used by both v1.2 and v2.0, but v2.0 formats the response differently in its handler
func (s *StoryService) GetStoryDisplayDetails(ctx context.Context, storyID string) (*storydomain.Story, *authordomain.Author, error) {
	s.Logger.Info(ctx, "Fetching story details",
		logger.String("story_id", storyID),
	)

	story, err := s.storyProvider.GetStory(ctx, storyID)
	if err != nil {
		s.Logger.Error(ctx, "Failed to fetch story",
			logger.String("story_id", storyID),
			logger.String("error", err.Error()),
		)
		return nil, nil, err
	}

	authorID := strconv.FormatUint(uint64(story.AuthorID), 10)
	author, err := s.authorProvider.GetAuthor(ctx, authorID)
	if err != nil {
		s.Logger.Error(ctx, "Failed to fetch author",
			logger.String("story_id", storyID),
			logger.String("author_id", authorID),
			logger.String("error", err.Error()),
		)
		return story, nil, err
	}

	s.Logger.Info(ctx, "Successfully fetched story and author details",
		logger.String("story_id", storyID),
		logger.String("author_id", authorID),
	)

	return story, author, nil
}
