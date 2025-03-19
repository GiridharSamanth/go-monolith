package service

import (
	"context"
	"strconv"
	"time"

	data "go-monolith/internal/bff/data"
	authordomain "go-monolith/internal/modules/author/domain"
	storydomain "go-monolith/internal/modules/story/domain"
	"go-monolith/pkg/logger"
	"go-monolith/pkg/metrics"
)

type StoryService struct {
	storyProvider  data.StoryDataProvider
	authorProvider data.AuthorDataProvider
	Logger         logger.Logger
	Metrics        *metrics.Client
}

var storyService *StoryService

func NewStoryService(sp data.StoryDataProvider, ap data.AuthorDataProvider, log logger.Logger, metrics *metrics.Client) *StoryService {
	if storyService == nil {
		storyService = &StoryService{
			storyProvider:  sp,
			authorProvider: ap,
			Logger:         log,
			Metrics:        metrics,
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
	start := time.Now()
	s.Logger.Info(ctx, "Fetching story details",
		logger.String("story_id", storyID),
	)

	// Record story fetch attempt
	s.Metrics.IncrementCounter("story.fetch.attempt", []string{
		"story_id:" + storyID,
	})

	story, err := s.storyProvider.GetStory(ctx, storyID)
	if err != nil {
		s.Logger.Error(ctx, "Failed to fetch story",
			logger.String("story_id", storyID),
			logger.String("error", err.Error()),
		)
		// Record story fetch error
		s.Metrics.IncrementCounter("story.fetch.error", []string{
			"story_id:" + storyID,
			"error_type:story_fetch",
		})
		return nil, nil, err
	}

	authorID := strconv.FormatUint(uint64(story.AuthorID), 10)
	// Record author fetch attempt
	s.Metrics.IncrementCounter("author.fetch.attempt", []string{
		"story_id:" + storyID,
		"author_id:" + authorID,
	})

	author, err := s.authorProvider.GetAuthor(ctx, authorID)
	if err != nil {
		s.Logger.Error(ctx, "Failed to fetch author",
			logger.String("story_id", storyID),
			logger.String("author_id", authorID),
			logger.String("error", err.Error()),
		)
		// Record author fetch error
		s.Metrics.IncrementCounter("author.fetch.error", []string{
			"story_id:" + storyID,
			"author_id:" + authorID,
			"error_type:author_fetch",
		})
		return story, nil, err
	}

	s.Logger.Info(ctx, "Successfully fetched story and author details",
		logger.String("story_id", storyID),
		logger.String("author_id", authorID),
	)

	// Record successful story and author fetch
	s.Metrics.IncrementCounter("story.fetch.success", []string{
		"story_id:" + storyID,
	})
	s.Metrics.IncrementCounter("author.fetch.success", []string{
		"story_id:" + storyID,
		"author_id:" + authorID,
	})

	// Record total operation duration
	duration := time.Since(start)
	s.Metrics.RecordTiming("story.display.duration", duration, []string{
		"story_id:" + storyID,
	})

	return story, author, nil
}
