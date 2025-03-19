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
	"go-monolith/pkg/metrics"
)

type AuthorService struct {
	repo    repository.AuthorRepository
	logger  logger.Logger
	metrics *metrics.Client
}

func NewAuthorService(repo repository.AuthorRepository, logger logger.Logger, metrics *metrics.Client) *AuthorService {
	return &AuthorService{
		repo:    repo,
		logger:  logger,
		metrics: metrics,
	}
}

// Write Operations (Commands)
func (s *AuthorService) Create(ctx context.Context, firstName, lastName, profileImageURL, slug string) (*domain.Author, error) {
	start := time.Now()
	s.logger.Info(ctx, "Creating new author",
		logger.String("first_name", firstName),
		logger.String("last_name", lastName),
		logger.String("slug", slug))

	// Record author creation attempt
	s.metrics.IncrementCounter("author.create.attempt", []string{
		"slug:" + slug,
	})

	author, err := domain.NewAuthor(firstName, lastName, profileImageURL, slug)
	if err != nil {
		s.logger.Error(ctx, "Failed to create author",
			logger.String("error", err.Error()),
			logger.String("first_name", firstName),
			logger.String("last_name", lastName),
			logger.String("slug", slug))
		// Record author creation error
		s.metrics.IncrementCounter("author.create.error", []string{
			"slug:" + slug,
			"error_type:validation",
		})
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
		// Record repository error
		s.metrics.IncrementCounter("author.create.error", []string{
			"slug:" + slug,
			"error_type:repository",
		})
		return nil, err
	}

	s.logger.Info(ctx, "Author created successfully",
		logger.String("author_id", fmt.Sprintf("%d", author.ID)),
		logger.String("slug", slug))

	// Record successful author creation
	s.metrics.IncrementCounter("author.create.success", []string{
		"author_id:" + fmt.Sprintf("%d", author.ID),
		"slug:" + slug,
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("author.create.duration", duration, []string{
		"author_id:" + fmt.Sprintf("%d", author.ID),
	})

	return author, nil
}

func (s *AuthorService) Update(ctx context.Context, id uint, firstName, lastName, profileImageURL string) error {
	start := time.Now()
	s.logger.Info(ctx, "Updating author",
		logger.String("author_id", fmt.Sprintf("%d", id)))

	// Record author update attempt
	s.metrics.IncrementCounter("author.update.attempt", []string{
		"author_id:" + fmt.Sprintf("%d", id),
	})

	author, err := s.repo.GetByID(ctx, id)
	if err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryUpdate(ctx, id, firstName, lastName, profileImageURL)
		}
		s.logger.Error(ctx, "Failed to get author for update",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", id)))
		// Record author fetch error
		s.metrics.IncrementCounter("author.update.error", []string{
			"author_id:" + fmt.Sprintf("%d", id),
			"error_type:fetch",
		})
		return err
	}

	if err := author.Update(firstName, lastName, profileImageURL); err != nil {
		s.logger.Error(ctx, "Failed to update author domain object",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", id)))
		// Record validation error
		s.metrics.IncrementCounter("author.update.error", []string{
			"author_id:" + fmt.Sprintf("%d", id),
			"error_type:validation",
		})
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
		// Record repository error
		s.metrics.IncrementCounter("author.update.error", []string{
			"author_id:" + fmt.Sprintf("%d", id),
			"error_type:repository",
		})
		return err
	}

	s.logger.Info(ctx, "Author updated successfully",
		logger.String("author_id", fmt.Sprintf("%d", id)))

	// Record successful author update
	s.metrics.IncrementCounter("author.update.success", []string{
		"author_id:" + fmt.Sprintf("%d", id),
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("author.update.duration", duration, []string{
		"author_id:" + fmt.Sprintf("%d", id),
	})

	return nil
}

func (s *AuthorService) Delete(ctx context.Context, id uint) error {
	start := time.Now()
	s.logger.Info(ctx, "Deleting author",
		logger.String("author_id", fmt.Sprintf("%d", id)))

	// Record author deletion attempt
	s.metrics.IncrementCounter("author.delete.attempt", []string{
		"author_id:" + fmt.Sprintf("%d", id),
	})

	if err := s.repo.Delete(ctx, id); err != nil {
		var baseErr *errors.BaseError
		if stderrors.As(err, &baseErr) && baseErr.Kind == errors.ErrKindTransient {
			return s.retryDelete(ctx, id)
		}
		s.logger.Error(ctx, "Failed to delete author",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", id)))
		// Record deletion error
		s.metrics.IncrementCounter("author.delete.error", []string{
			"author_id:" + fmt.Sprintf("%d", id),
			"error_type:repository",
		})
		return err
	}

	s.logger.Info(ctx, "Author deleted successfully",
		logger.String("author_id", fmt.Sprintf("%d", id)))

	// Record successful author deletion
	s.metrics.IncrementCounter("author.delete.success", []string{
		"author_id:" + fmt.Sprintf("%d", id),
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("author.delete.duration", duration, []string{
		"author_id:" + fmt.Sprintf("%d", id),
	})

	return nil
}

// Read Operations (Queries)
func (s *AuthorService) GetByID(ctx context.Context, id uint) (*domain.Author, error) {
	start := time.Now()
	s.logger.Debug(ctx, "Getting author by ID",
		logger.String("author_id", fmt.Sprintf("%d", id)))

	// Record author fetch attempt
	s.metrics.IncrementCounter("author.fetch.attempt", []string{
		"author_id:" + fmt.Sprintf("%d", id),
		"type:id",
	})

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
				// Record not found error
				s.metrics.IncrementCounter("author.fetch.error", []string{
					"author_id:" + fmt.Sprintf("%d", id),
					"error_type:not_found",
					"type:id",
				})
				return nil, err
			}
		}
		s.logger.Error(ctx, "Failed to get author",
			logger.String("error", err.Error()),
			logger.String("author_id", fmt.Sprintf("%d", id)))
		// Record fetch error
		s.metrics.IncrementCounter("author.fetch.error", []string{
			"author_id:" + fmt.Sprintf("%d", id),
			"error_type:repository",
			"type:id",
		})
		return nil, err
	}

	// Record successful author fetch
	s.metrics.IncrementCounter("author.fetch.success", []string{
		"author_id:" + fmt.Sprintf("%d", id),
		"type:id",
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("author.fetch.duration", duration, []string{
		"author_id:" + fmt.Sprintf("%d", id),
		"type:id",
	})

	return author, nil
}

func (s *AuthorService) GetBySlug(ctx context.Context, slug string) (*domain.Author, error) {
	start := time.Now()
	s.logger.Debug(ctx, "Getting author by slug",
		logger.String("slug", slug))

	// Record author fetch attempt
	s.metrics.IncrementCounter("author.fetch.attempt", []string{
		"slug:" + slug,
		"type:slug",
	})

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
				// Record not found error
				s.metrics.IncrementCounter("author.fetch.error", []string{
					"slug:" + slug,
					"error_type:not_found",
					"type:slug",
				})
				return nil, err
			}
		}
		s.logger.Error(ctx, "Failed to get author by slug",
			logger.String("error", err.Error()),
			logger.String("slug", slug))
		// Record fetch error
		s.metrics.IncrementCounter("author.fetch.error", []string{
			"slug:" + slug,
			"error_type:repository",
			"type:slug",
		})
		return nil, err
	}

	// Record successful author fetch
	s.metrics.IncrementCounter("author.fetch.success", []string{
		"slug:" + slug,
		"type:slug",
	})

	// Record operation duration
	duration := time.Since(start)
	s.metrics.RecordTiming("author.fetch.duration", duration, []string{
		"slug:" + slug,
		"type:slug",
	})

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
