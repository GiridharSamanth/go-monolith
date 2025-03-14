package domain

import (
	"strconv"
	"time"

	"go-monolith/pkg/errors"

	"github.com/go-playground/validator/v10"
)

// Story represents the story domain entity
type Story struct {
	ID          uint
	Title       string `validate:"required,min=3,max=255"`
	Content     string `validate:"required,min=10"`
	AuthorID    uint   `validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
}

var validate = validator.New()

func NewStory(title, content, authorID string) (*Story, error) {
	// Validate inputs before creating
	if err := validateInputs(title, content, authorID); err != nil {
		return nil, err
	}

	now := time.Now()
	authorIDUint, _ := strconv.ParseUint(authorID, 10, 64)
	story := &Story{
		Title:     title,
		Content:   content,
		AuthorID:  uint(authorIDUint),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Validate the struct
	if err := validate.Struct(story); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	return story, nil
}

func (s *Story) Update(title, content string) error {
	// Validate inputs
	if err := validateInputs(title, content, strconv.FormatUint(uint64(s.AuthorID), 10)); err != nil {
		return err
	}

	s.Title = title
	s.Content = content
	s.UpdatedAt = time.Now()

	// Validate the struct after update
	if err := validate.Struct(s); err != nil {
		return errors.NewValidationError(err.Error())
	}

	return nil
}

func (s *Story) Publish() error {
	// Add any pre-publication validation rules
	if s.Title == "" || s.Content == "" {
		return errors.NewValidationError("cannot publish story without title or content")
	}

	now := time.Now()
	s.PublishedAt = &now
	s.UpdatedAt = now
	return nil
}

// validateInputs performs validation on raw input strings
func validateInputs(title, content, authorID string) error {
	if title == "" {
		return errors.NewValidationError("title cannot be empty")
	}
	if len(title) < 3 {
		return errors.NewValidationError("title must be at least 3 characters long")
	}
	if len(title) > 255 {
		return errors.NewValidationError("title cannot exceed 255 characters")
	}

	if content == "" {
		return errors.NewValidationError("content cannot be empty")
	}
	if len(content) < 10 {
		return errors.NewValidationError("content must be at least 10 characters long")
	}

	if authorID == "" {
		return errors.NewValidationError("author ID cannot be empty")
	}
	if _, err := strconv.ParseUint(authorID, 10, 64); err != nil {
		return errors.NewValidationError("invalid author ID format")
	}

	return nil
}
