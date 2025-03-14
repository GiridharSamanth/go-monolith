package service

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"go-monolith/internal/modules/story/domain"
	"go-monolith/internal/modules/story/repository"
	"go-monolith/pkg/errors"
)

type StoryService struct {
	repo repository.StoryRepository
}

func NewStoryService(repo repository.StoryRepository) *StoryService {
	return &StoryService{
		repo: repo,
	}
}

// Write Operations (Commands)
func (s *StoryService) Create(ctx context.Context, title, content, authorID string) (*domain.Story, error) {
	story, err := domain.NewStory(title, content, authorID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, story); err != nil {
		return nil, err
	}
	return story, nil
}

func (s *StoryService) Update(ctx context.Context, id, title, content string) (*domain.Story, error) {
	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdateStory(ctx, id, title, content)
		}
		return nil, err
	}

	if err := story.Update(title, content); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, story); err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdate(ctx, story)
		}
		return nil, err
	}
	return story, nil
}

func (s *StoryService) Publish(ctx context.Context, id string) (*domain.Story, error) {
	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	story.Publish()
	if err := s.repo.Update(ctx, story); err != nil {
		return nil, err
	}
	return story, nil
}

func (s *StoryService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// Read Operations (Queries)
func (s *StoryService) GetByID(ctx context.Context, id string) (*domain.Story, error) {
	story, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) {
			switch baseErr.Kind {
			case errors.ErrKindTransient:
				return s.retryGet(ctx, id)
			case errors.ErrKindNotFound:
				return nil, err
			}
		}
		return nil, err
	}
	return story, nil
}

func (s *StoryService) List(ctx context.Context, limit, offset int) ([]*domain.Story, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *StoryService) ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*domain.Story, error) {
	return s.repo.ListByAuthor(ctx, authorID, limit, offset)
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
