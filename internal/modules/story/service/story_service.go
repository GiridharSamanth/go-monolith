package service

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"go-monolith/internal/modules/story/domain"
	"go-monolith/internal/modules/story/repository"
	"go-monolith/pkg/errors"
	"go-monolith/pkg/logger"
	"go-monolith/pkg/metrics"
)

type StoryService struct {
	repo    repository.StoryRepository
	logger  logger.Logger
	metrics *metrics.Client
}

func NewStoryService(repo repository.StoryRepository, logger logger.Logger, metrics *metrics.Client) *StoryService {
	return &StoryService{
		repo:    repo,
		logger:  logger,
		metrics: metrics,
	}
}

// Write Operations (Commands)
func (s *StoryService) Create(ctx context.Context, title, content, authorID string) (*domain.Story, error) {
	start := time.Now()
	s.logger.Info(ctx, "Creating new story", logger.String("author_id", authorID))

	// Record story creation attempt
	s.metrics.IncrementCounter("story.create.attempt", []string{
		"author_id:" + authorID,
	})

	story, err := domain.NewStory(title, content, authorID)
	if err != nil {
		s.logger.Error(ctx, "Failed to create story", logger.String("error", err.Error()))
		// Record validation error
		s.metrics.IncrementCounter("story.create.error", []string{
			"author_id:" + authorID,
			"error_type:validation",
		})
		return nil, err
	}

	if err := s.repo.Create(ctx, story); err != nil {
		s.logger.Error(ctx, "Failed to save story to repository",
			logger.String("error", err.Error()),
			logger.String("story_id", fmt.Sprintf("%d", story.ID)))
		// Record repository error
		s.metrics.IncrementCounter("story.create.error", []string{
			"author_id:" + authorID,
			"error_type:repository",
		})
		return nil, err
	}

	s.logger.Info(ctx, "Story created successfully",
		logger.String("story_id", fmt.Sprintf("%d", story.ID)),
		logger.String("author_id", authorID))

	// Record successful story creation
	s.metrics.IncrementCounter("story.create.success", []string{
		"story_id:" + fmt.Sprintf("%d", story.ID),
		"author_id:" + authorID,
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("story.create.duration", duration, []string{
		"story_id:" + fmt.Sprintf("%d", story.ID),
	})

	return story, nil
}

func (s *StoryService) Update(ctx context.Context, id, title, content string) (*domain.Story, error) {
	start := time.Now()
	s.logger.Info(ctx, "Updating story", logger.String("story_id", id))

	// Record story update attempt
	s.metrics.IncrementCounter("story.update.attempt", []string{
		"story_id:" + id,
	})

	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdateStory(ctx, id, title, content)
		}
		s.logger.Error(ctx, "Failed to get story for update",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		// Record fetch error
		s.metrics.IncrementCounter("story.update.error", []string{
			"story_id:" + id,
			"error_type:fetch",
		})
		return nil, err
	}

	if err := story.Update(title, content); err != nil {
		s.logger.Error(ctx, "Failed to update story domain object",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		// Record validation error
		s.metrics.IncrementCounter("story.update.error", []string{
			"story_id:" + id,
			"error_type:validation",
		})
		return nil, err
	}

	if err := s.repo.Update(ctx, story); err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdate(ctx, story)
		}
		s.logger.Error(ctx, "Failed to save story update",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		// Record repository error
		s.metrics.IncrementCounter("story.update.error", []string{
			"story_id:" + id,
			"error_type:repository",
		})
		return nil, err
	}

	s.logger.Info(ctx, "Story updated successfully", logger.String("story_id", id))

	// Record successful story update
	s.metrics.IncrementCounter("story.update.success", []string{
		"story_id:" + id,
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("story.update.duration", duration, []string{
		"story_id:" + id,
	})

	return story, nil
}

func (s *StoryService) Publish(ctx context.Context, id string) (*domain.Story, error) {
	start := time.Now()
	s.logger.Info(ctx, "Publishing story", logger.String("story_id", id))

	// Record story publish attempt
	s.metrics.IncrementCounter("story.publish.attempt", []string{
		"story_id:" + id,
	})

	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get story for publishing",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		// Record fetch error
		s.metrics.IncrementCounter("story.publish.error", []string{
			"story_id:" + id,
			"error_type:fetch",
		})
		return nil, err
	}

	story.Publish()
	if err := s.repo.Update(ctx, story); err != nil {
		s.logger.Error(ctx, "Failed to save story publication",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		// Record repository error
		s.metrics.IncrementCounter("story.publish.error", []string{
			"story_id:" + id,
			"error_type:repository",
		})
		return nil, err
	}

	s.logger.Info(ctx, "Story published successfully", logger.String("story_id", id))

	// Record successful story publication
	s.metrics.IncrementCounter("story.publish.success", []string{
		"story_id:" + id,
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("story.publish.duration", duration, []string{
		"story_id:" + id,
	})

	return story, nil
}

func (s *StoryService) Delete(ctx context.Context, id string) error {
	start := time.Now()
	s.logger.Info(ctx, "Deleting story", logger.String("story_id", id))

	// Record story deletion attempt
	s.metrics.IncrementCounter("story.delete.attempt", []string{
		"story_id:" + id,
	})

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to delete story",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		// Record deletion error
		s.metrics.IncrementCounter("story.delete.error", []string{
			"story_id:" + id,
			"error_type:repository",
		})
		return err
	}

	s.logger.Info(ctx, "Story deleted successfully", logger.String("story_id", id))

	// Record successful story deletion
	s.metrics.IncrementCounter("story.delete.success", []string{
		"story_id:" + id,
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("story.delete.duration", duration, []string{
		"story_id:" + id,
	})

	return nil
}

// Read Operations (Queries)
func (s *StoryService) GetByID(ctx context.Context, id string) (*domain.Story, error) {
	start := time.Now()
	s.logger.Debug(ctx, "Getting story by ID", logger.String("story_id", id))

	// Record story fetch attempt
	s.metrics.IncrementCounter("story.fetch.attempt", []string{
		"story_id:" + id,
		"type:id",
	})

	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) {
			switch baseErr.Kind {
			case errors.ErrKindTransient:
				return s.retryGet(ctx, id)
			case errors.ErrKindNotFound:
				s.logger.Warn(ctx, "Story not found", logger.String("story_id", id))
				// Record not found error
				s.metrics.IncrementCounter("story.fetch.error", []string{
					"story_id:" + id,
					"error_type:not_found",
					"type:id",
				})
				return nil, err
			}
		}
		s.logger.Error(ctx, "Failed to get story",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		// Record fetch error
		s.metrics.IncrementCounter("story.fetch.error", []string{
			"story_id:" + id,
			"error_type:repository",
			"type:id",
		})
		return nil, err
	}

	// Record successful story fetch
	s.metrics.IncrementCounter("story.fetch.success", []string{
		"story_id:" + id,
		"type:id",
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("story.fetch.duration", duration, []string{
		"story_id:" + id,
		"type:id",
	})

	return story, nil
}

func (s *StoryService) List(ctx context.Context, limit, offset int) ([]*domain.Story, error) {
	start := time.Now()
	s.logger.Debug(ctx, "Listing stories",
		logger.Int("limit", limit),
		logger.Int("offset", offset))

	// Record story list attempt
	s.metrics.IncrementCounter("story.list.attempt", []string{
		"limit:" + fmt.Sprintf("%d", limit),
		"offset:" + fmt.Sprintf("%d", offset),
	})

	stories, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error(ctx, "Failed to list stories",
			logger.String("error", err.Error()),
			logger.Int("limit", limit),
			logger.Int("offset", offset))
		// Record list error
		s.metrics.IncrementCounter("story.list.error", []string{
			"error_type:repository",
		})
		return nil, err
	}

	// Record successful story list
	s.metrics.IncrementCounter("story.list.success", []string{
		"count:" + fmt.Sprintf("%d", len(stories)),
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("story.list.duration", duration, []string{
		"count:" + fmt.Sprintf("%d", len(stories)),
	})

	return stories, nil
}

func (s *StoryService) ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*domain.Story, error) {
	start := time.Now()
	s.logger.Debug(ctx, "Listing stories by author",
		logger.String("author_id", authorID),
		logger.Int("limit", limit),
		logger.Int("offset", offset))

	// Record story list by author attempt
	s.metrics.IncrementCounter("story.list.attempt", []string{
		"author_id:" + authorID,
		"limit:" + fmt.Sprintf("%d", limit),
		"offset:" + fmt.Sprintf("%d", offset),
		"type:author",
	})

	stories, err := s.repo.ListByAuthor(ctx, authorID, limit, offset)
	if err != nil {
		s.logger.Error(ctx, "Failed to list stories by author",
			logger.String("error", err.Error()),
			logger.String("author_id", authorID),
			logger.Int("limit", limit),
			logger.Int("offset", offset))
		// Record list error
		s.metrics.IncrementCounter("story.list.error", []string{
			"author_id:" + authorID,
			"error_type:repository",
			"type:author",
		})
		return nil, err
	}

	// Record successful story list by author
	s.metrics.IncrementCounter("story.list.success", []string{
		"author_id:" + authorID,
		"count:" + fmt.Sprintf("%d", len(stories)),
		"type:author",
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("story.list.duration", duration, []string{
		"author_id:" + authorID,
		"count:" + fmt.Sprintf("%d", len(stories)),
		"type:author",
	})

	return stories, nil
}

// Helper method for retrying operations
func (s *StoryService) retryGet(ctx context.Context, id string) (*domain.Story, error) {
	for i := 0; i < 3; i++ {
		story, err := s.repo.GetByID(ctx, id)
		if err == nil {
			return story, nil
		}
		var baseErr *errors.BaseError
		if !stderrors.As(err, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
			return nil, err
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return nil, errors.NewUnexpectedError(fmt.Errorf("max retries exceeded"))
}

func (s *StoryService) retryUpdateStory(ctx context.Context, id, title, content string) (*domain.Story, error) {
	for i := 0; i < 3; i++ {
		story, err := s.repo.GetByID(ctx, id)
		if err == nil {
			story.Update(title, content)
			if err := s.repo.Update(ctx, story); err != nil {
				var baseErr *errors.BaseError
				if !stderrors.As(err, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
					return nil, err
				}
			}
			return story, nil
		}
		var baseErr *errors.BaseError
		if !stderrors.As(err, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
			return nil, err
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return nil, errors.NewUnexpectedError(fmt.Errorf("max retries exceeded"))
}

func (s *StoryService) retryUpdate(ctx context.Context, story *domain.Story) (*domain.Story, error) {
	for i := 0; i < 3; i++ {
		if err := s.repo.Update(ctx, story); err != nil {
			var baseErr *errors.BaseError
			if !stderrors.As(err, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
				return nil, err
			}
		}
		return story, nil
	}
	return nil, errors.NewUnexpectedError(fmt.Errorf("max retries exceeded"))
}

// TODO: Add domain events like story published, story updated, story deleted
