package service

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"go-monolith/internal/modules/author/domain"
	"go-monolith/internal/modules/author/repository"
	"go-monolith/pkg/errors"
)

type AuthorService struct {
	repo repository.AuthorRepository
}

func NewAuthorService(repo repository.AuthorRepository) *AuthorService {
	return &AuthorService{repo: repo}
}

// Write Operations (Commands)
func (s *AuthorService) Create(ctx context.Context, firstName, lastName, profileImageURL, slug string) (*domain.Author, error) {
	author, err := domain.NewAuthor(firstName, lastName, profileImageURL, slug)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, author); err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryCreate(ctx, author)
		}
		return nil, err
	}

	return author, nil
}

func (s *AuthorService) Update(ctx context.Context, id uint, firstName, lastName, profileImageURL string) error {
	author, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdate(ctx, id, firstName, lastName, profileImageURL)
		}
		return err
	}

	if err := author.Update(firstName, lastName, profileImageURL); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, author); err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdate(ctx, id, firstName, lastName, profileImageURL)
		}
		return err
	}
	return nil
}

func (s *AuthorService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryDelete(ctx, id)
		}
		return err
	}
	return nil
}

// Read Operations (Queries)
func (s *AuthorService) GetByID(ctx context.Context, id uint) (*domain.Author, error) {
	author, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) {
			switch baseErr.Kind {
			case errors.ErrKindTransient:
				return s.retryGetByID(ctx, id)
			case errors.ErrKindNotFound:
				return nil, err
			}
		}
		return nil, err
	}
	return author, nil
}

func (s *AuthorService) GetBySlug(ctx context.Context, slug string) (*domain.Author, error) {
	author, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) {
			switch baseErr.Kind {
			case errors.ErrKindTransient:
				return s.retryGetBySlug(ctx, slug)
			case errors.ErrKindNotFound:
				return nil, err
			}
		}
		return nil, err
	}
	return author, nil
}

// Retry Operations
func (s *AuthorService) retryCreate(ctx context.Context, author *domain.Author) (*domain.Author, error) {
	for i := 0; i < 3; i++ {
		createErr := s.repo.Create(ctx, author)
		if createErr == nil {
			return author, nil
		}
		var baseErr *errors.BaseError
		if !stderrors.As(createErr, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
			return nil, createErr
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return nil, errors.NewUnexpectedError(fmt.Errorf("max retries exceeded"))
}

func (s *AuthorService) retryUpdate(ctx context.Context, id uint, firstName, lastName, profileImageURL string) error {
	for i := 0; i < 3; i++ {
		author, err := s.repo.GetByID(ctx, id)
		if err == nil {
			if err := author.Update(firstName, lastName, profileImageURL); err != nil {
				return err
			}
			updateErr := s.repo.Update(ctx, author)
			if updateErr == nil {
				return nil
			}
			var baseErr *errors.BaseError
			if !stderrors.As(updateErr, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
				return updateErr
			}
		} else {
			var baseErr *errors.BaseError
			if !stderrors.As(err, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
				return err
			}
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return errors.NewUnexpectedError(fmt.Errorf("max retries exceeded"))
}

func (s *AuthorService) retryDelete(ctx context.Context, id uint) error {
	for i := 0; i < 3; i++ {
		err := s.repo.Delete(ctx, id)
		if err == nil {
			return nil
		}
		var baseErr *errors.BaseError
		if !stderrors.As(err, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
			return err
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return errors.NewUnexpectedError(fmt.Errorf("max retries exceeded"))
}

func (s *AuthorService) retryGetByID(ctx context.Context, id uint) (*domain.Author, error) {
	for i := 0; i < 3; i++ {
		author, findErr := s.repo.GetByID(ctx, id)
		if findErr == nil {
			return author, nil
		}
		var baseErr *errors.BaseError
		if !stderrors.As(findErr, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
			return nil, findErr
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return nil, errors.NewUnexpectedError(fmt.Errorf("max retries exceeded"))
}

func (s *AuthorService) retryGetBySlug(ctx context.Context, slug string) (*domain.Author, error) {
	for i := 0; i < 3; i++ {
		author, findErr := s.repo.GetBySlug(ctx, slug)
		if findErr == nil {
			return author, nil
		}
		var baseErr *errors.BaseError
		if !stderrors.As(findErr, &baseErr) || baseErr.Kind != errors.ErrKindTransient {
			return nil, findErr
		}
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return nil, errors.NewUnexpectedError(fmt.Errorf("max retries exceeded"))
}
