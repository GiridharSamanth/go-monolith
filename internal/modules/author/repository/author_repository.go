package repository

import (
	"context"
	stderrors "errors"

	"gorm.io/gorm"

	"go-monolith/internal/modules/author/domain"
	"go-monolith/pkg/errors"

	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// authorModel represents the database model
type authorModel struct {
	ID              uint         `gorm:"primaryKey;autoIncrement"`
	FirstName       string       `gorm:"type:varchar(255);not null"`
	LastName        string       `gorm:"type:varchar(255);not null"`
	ProfileImageURL string       `gorm:"type:varchar(255);not null"`
	Slug            string       `gorm:"type:varchar(255);not null;uniqueIndex"`
	CreatedAt       sql.NullTime `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       sql.NullTime `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName sets the insert table name for this struct type
func (authorModel) TableName() string {
	return "authors" // Use the correct table name
}

// In this context, Only benefit of using interface is to allow for mocking in tests, otherwise not needed
type AuthorRepository interface {
	Create(ctx context.Context, author *domain.Author) error
	GetByID(ctx context.Context, id uint) (*domain.Author, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Author, error)
	Update(ctx context.Context, author *domain.Author) error
	Delete(ctx context.Context, id uint) error
}

type authorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{db: db}
}

func (r *authorRepository) Create(ctx context.Context, author *domain.Author) error {
	model := toModel(author)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if isDuplicateKeyError(err) {
			return errors.NewValidationError("author with this slug already exists")
		}
		if isTransientError(err) {
			return errors.NewTransientError(err)
		}
		return errors.NewUnexpectedError(err)
	}
	author.ID = model.ID
	return nil
}

func (r *authorRepository) GetByID(ctx context.Context, id uint) (*domain.Author, error) {
	var model authorModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("author", string(id))
		}
		if err == gorm.ErrInvalidTransaction || err == gorm.ErrRegistered {
			return nil, errors.NewTransientError(err)
		}
		return nil, errors.NewUnexpectedError(err)
	}
	return toDomain(&model), nil
}

func (r *authorRepository) GetBySlug(ctx context.Context, slug string) (*domain.Author, error) {
	var model authorModel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("author", slug)
		}
		if err == gorm.ErrInvalidTransaction || err == gorm.ErrRegistered {
			return nil, errors.NewTransientError(err)
		}
		return nil, errors.NewUnexpectedError(err)
	}
	return toDomain(&model), nil
}

func (r *authorRepository) Update(ctx context.Context, author *domain.Author) error {
	model := toModel(author)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		if isDuplicateKeyError(err) {
			return errors.NewValidationError("author with this slug already exists")
		}
		if isTransientError(err) {
			return errors.NewTransientError(err)
		}
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func (r *authorRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&authorModel{}, id).Error; err != nil {
		if isTransientError(err) {
			return errors.NewTransientError(err)
		}
		return errors.NewUnexpectedError(err)
	}
	return nil
}

// toModel converts domain author to database model
func toModel(author *domain.Author) *authorModel {
	return &authorModel{
		ID:              author.ID,
		FirstName:       author.FirstName,
		LastName:        author.LastName,
		ProfileImageURL: author.ProfileImageURL,
		Slug:            author.Slug,
		CreatedAt:       sql.NullTime{Time: author.CreatedAt, Valid: !author.CreatedAt.IsZero()},
		UpdatedAt:       sql.NullTime{Time: author.UpdatedAt, Valid: !author.UpdatedAt.IsZero()},
	}
}

// toDomain converts database model to domain author
func toDomain(model *authorModel) *domain.Author {
	return &domain.Author{
		ID:              model.ID,
		FirstName:       model.FirstName,
		LastName:        model.LastName,
		ProfileImageURL: model.ProfileImageURL,
		Slug:            model.Slug,
		CreatedAt:       model.CreatedAt.Time,
		UpdatedAt:       model.UpdatedAt.Time,
	}
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return stderrors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func isTransientError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if !stderrors.As(err, &mysqlErr) {
		return false
	}
	// Common MySQL transient error codes
	switch mysqlErr.Number {
	case 1213, // Deadlock
		1205, // Lock wait timeout
		2006, // MySQL server has gone away
		2013: // Lost connection to MySQL server
		return true
	}
	return false
}
