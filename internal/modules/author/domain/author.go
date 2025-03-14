package domain

import (
	"net/url"
	"time"

	"go-monolith/pkg/errors"

	"github.com/go-playground/validator/v10"
)

// Author represents the author domain entity
type Author struct {
	ID              uint
	FirstName       string `validate:"required,min=3"`
	LastName        string `validate:"required,min=1"`
	ProfileImageURL string `validate:"required,url"`
	Slug            string `validate:"required,min=8"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

var validate = validator.New()

func NewAuthor(firstName, lastName, profileImageURL, slug string) (*Author, error) {
	// Validate inputs before creating
	if err := validateInputs(firstName, lastName, profileImageURL, slug); err != nil {
		return nil, err
	}

	now := time.Now()
	author := &Author{
		FirstName:       firstName,
		LastName:        lastName,
		ProfileImageURL: profileImageURL,
		Slug:            slug,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Validate the struct
	if err := validate.Struct(author); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	return author, nil
}

func (a *Author) Update(firstName, lastName, profileImageURL string) error {
	// Validate inputs
	if err := validateInputs(firstName, lastName, profileImageURL, a.Slug); err != nil {
		return err
	}

	a.FirstName = firstName
	a.LastName = lastName
	a.ProfileImageURL = profileImageURL
	a.UpdatedAt = time.Now()

	// Validate the struct after update
	if err := validate.Struct(a); err != nil {
		return errors.NewValidationError(err.Error())
	}

	return nil
}

// validateInputs performs validation on raw input strings
func validateInputs(firstName, lastName, profileImageURL, slug string) error {
	if len(firstName) < 3 {
		return NewInvalidFirstNameError()
	}

	if lastName == "" {
		return NewInvalidLastNameError()
	}

	if _, err := url.ParseRequestURI(profileImageURL); err != nil {
		return NewInvalidProfileImageError()
	}

	if len(slug) < 8 {
		return NewInvalidSlugError()
	}

	return nil
}
