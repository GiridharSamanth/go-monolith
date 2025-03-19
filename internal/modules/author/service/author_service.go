package service

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"go-monolith/internal/modules/author/domain"
	"go-monolith/internal/modules/author/repository"
	"go-monolith/pkg/errors"
	"go-monolith/pkg/logger"
)

type AuthorService struct {
	repo   repository.AuthorRepository
	logger logger.Logger
}

func NewAuthorService(repo repository.AuthorRepository, logger logger.Logger) *AuthorService {
	return &AuthorService{repo: repo, logger: logger}
}

// Write Operations (Commands)
func (s *AuthorService) Create(ctx context.Context, firstName, lastName, profileImageURL, slug string) (*domain.Author, error) {
	s.logger.Info(ctx, "Creating new author",
		logger.String("first_name", firstName),
		logger.String("last_name", lastName),
		logger.String("slug", slug))

	author, err := domain.NewAuthor(firstName, lastName, profileImageURL, slug)
	if err != nil {
		s.logger.Error(ctx, "Failed to create author",
			logger.String("error", err.Error()),
			logger.String("first_name", firstName),
			logger.String("last_name", lastName),
			logger.String("slug", slug))
		return nil, err
	}

	if err := s.repo.Create(ctx, author); err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryCreate(ctx, author)
		}
		s.logger.Error(ctx, "Failed to save author to repository",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", author.ID)))
		return nil, err
	}

	s.logger.Info(ctx, "Author created successfully",
		logger.String("author_id", fmt.Sprintf("%d", author.ID)),
		logger.String("slug", slug))
	return author, nil
}

func (s *AuthorService) Update(ctx context.Context, id uint, firstName, lastName, profileImageURL string) error {
	s.logger.Info(ctx, "Updating author",
		logger.String("author_id", fmt.Sprintf("%d", id)))

	author, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdate(ctx, id, firstName, lastName, profileImageURL)
		}
		s.logger.Error(ctx, "Failed to get author for update",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", id)))
		return err
	}

	if err := author.Update(firstName, lastName, profileImageURL); err != nil {
		s.logger.Error(ctx, "Failed to update author domain object",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", id)))
		return err
	}

	if err := s.repo.Update(ctx, author); err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdate(ctx, id, firstName, lastName, profileImageURL)
		}
		s.logger.Error(ctx, "Failed to save author update",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", id)))
		return err
	}

	s.logger.Info(ctx, "Author updated successfully",
		logger.String("author_id", fmt.Sprintf("%d", id)))
	return nil
}

func (s *AuthorService) Delete(ctx context.Context, id uint) error {
	s.logger.Info(ctx, "Deleting author",
		logger.String("author_id", fmt.Sprintf("%d", id)))

	if err := s.repo.Delete(ctx, id); err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryDelete(ctx, id)
		}
		s.logger.Error(ctx, "Failed to delete author",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", id)))
		return err
	}

	s.logger.Info(ctx, "Author deleted successfully",
		logger.String("author_id", fmt.Sprintf("%d", id)))
	return nil
}

// Read Operations (Queries)
func (s *AuthorService) GetByID(ctx context.Context, id uint) (*domain.Author, error) {
	s.logger.Debug(ctx, "Getting author by ID",
		logger.String("author_id", fmt.Sprintf("%d", id)))

	author, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) {
			switch baseErr.Kind {
			case errors.ErrKindTransient:
				return s.retryGetByID(ctx, id)
			case errors.ErrKindNotFound:
				s.logger.Warn(ctx, "Author not found",
					logger.String("author_id", fmt.Sprintf("%d", id)))
				return nil, err
			}
		}
		s.logger.Error(ctx, "Failed to get author",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", id)))
		return nil, err
	}
	return author, nil
}

func (s *AuthorService) GetBySlug(ctx context.Context, slug string) (*domain.Author, error) {
	s.logger.Debug(ctx, "Getting author by slug",
		logger.String("slug", slug))

	author, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) {
			switch baseErr.Kind {
			case errors.ErrKindTransient:
				return s.retryGetBySlug(ctx, slug)
			case errors.ErrKindNotFound:
				s.logger.Warn(ctx, "Author not found",
					logger.String("slug", slug))
				return nil, err
			}
		}
		s.logger.Error(ctx, "Failed to get author by slug",
			logger.String("error", err.Error()),
			logger.String("slug", slug))
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
