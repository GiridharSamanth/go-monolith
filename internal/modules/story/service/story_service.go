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
)

type StoryService struct {
	repo   repository.StoryRepository
	logger logger.Logger
}

func NewStoryService(repo repository.StoryRepository, logger logger.Logger) *StoryService {
	return &StoryService{
		repo:   repo,
		logger: logger,
	}
}

// Write Operations (Commands)
func (s *StoryService) Create(ctx context.Context, title, content, authorID string) (*domain.Story, error) {
	s.logger.Info(ctx, "Creating new story", logger.String("author_id", authorID))

	story, err := domain.NewStory(title, content, authorID)
	if err != nil {
		s.logger.Error(ctx, "Failed to create story", logger.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Create(ctx, story); err != nil {
		s.logger.Error(ctx, "Failed to save story to repository",
			logger.String("error", err.Error()),
			logger.String("story_id", fmt.Sprintf("%d", story.ID)))
		return nil, err
	}

	s.logger.Info(ctx, "Story created successfully",
		logger.String("story_id", fmt.Sprintf("%d", story.ID)),
		logger.String("author_id", authorID))
	return story, nil
}

func (s *StoryService) Update(ctx context.Context, id, title, content string) (*domain.Story, error) {
	s.logger.Info(ctx, "Updating story", logger.String("story_id", id))

	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdateStory(ctx, id, title, content)
		}
		s.logger.Error(ctx, "Failed to get story for update",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		return nil, err
	}

	if err := story.Update(title, content); err != nil {
		s.logger.Error(ctx, "Failed to update story domain object",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
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
		return nil, err
	}

	s.logger.Info(ctx, "Story updated successfully", logger.String("story_id", id))
	return story, nil
}

func (s *StoryService) Publish(ctx context.Context, id string) (*domain.Story, error) {
	s.logger.Info(ctx, "Publishing story", logger.String("story_id", id))

	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get story for publishing",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		return nil, err
	}

	story.Publish()
	if err := s.repo.Update(ctx, story); err != nil {
		s.logger.Error(ctx, "Failed to save story publication",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		return nil, err
	}

	s.logger.Info(ctx, "Story published successfully", logger.String("story_id", id))
	return story, nil
}

func (s *StoryService) Delete(ctx context.Context, id string) error {
	s.logger.Info(ctx, "Deleting story", logger.String("story_id", id))

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to delete story",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		return err
	}

	s.logger.Info(ctx, "Story deleted successfully", logger.String("story_id", id))
	return nil
}

// Read Operations (Queries)
func (s *StoryService) GetByID(ctx context.Context, id string) (*domain.Story, error) {
	s.logger.Debug(ctx, "Getting story by ID", logger.String("story_id", id))

	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) {
			switch baseErr.Kind {
			case errors.ErrKindTransient:
				return s.retryGet(ctx, id)
			case errors.ErrKindNotFound:
				s.logger.Warn(ctx, "Story not found", logger.String("story_id", id))
				return nil, err
			}
		}
		s.logger.Error(ctx, "Failed to get story",
			logger.String("error", err.Error()),
			logger.String("story_id", id))
		return nil, err
	}
	return story, nil
}

func (s *StoryService) List(ctx context.Context, limit, offset int) ([]*domain.Story, error) {
	s.logger.Debug(ctx, "Listing stories",
		logger.Int("limit", limit),
		logger.Int("offset", offset))

	stories, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error(ctx, "Failed to list stories",
			logger.String("error", err.Error()),
			logger.Int("limit", limit),
			logger.Int("offset", offset))
		return nil, err
	}
	return stories, nil
}

func (s *StoryService) ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*domain.Story, error) {
	s.logger.Debug(ctx, "Listing stories by author",
		logger.String("author_id", authorID),
		logger.Int("limit", limit),
		logger.Int("offset", offset))

	stories, err := s.repo.ListByAuthor(ctx, authorID, limit, offset)
	if err != nil {
		s.logger.Error(ctx, "Failed to list stories by author",
			logger.String("error", err.Error()),
			logger.String("author_id", authorID),
			logger.Int("limit", limit),
			logger.Int("offset", offset))
		return nil, err
	}
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
