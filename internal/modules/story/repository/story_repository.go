package repository

import (
	"context"
	stderrors "errors"
	"strconv"
	"time"

	"gorm.io/gorm"

	"go-monolith/internal/modules/story/domain"
	"go-monolith/pkg/errors"

	"github.com/go-sql-driver/mysql"
)

// storyModel represents the database model
type storyModel struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Content     string    `gorm:"type:text;not null"`
	AuthorID    uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	PublishedAt *time.Time
	Views       int64 `gorm:"not null;default:0"`
	Likes       int64 `gorm:"not null;default:0"`
	Comments    int64 `gorm:"not null;default:0"`
}

// TableName sets the insert table name for this struct type
func (storyModel) TableName() string {
	return "stories" // Use the correct table name
}

// StoryRepository interface defines the contract for story repository operations
// In this context, only benefit of using interface is to allow for mocking in tests, otherwise not needed
type StoryRepository interface {
	Create(ctx context.Context, story *domain.Story) error
	Update(ctx context.Context, story *domain.Story) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.Story, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Story, error)
	ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*domain.Story, error)
}

type storyRepository struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) StoryRepository {
	return &storyRepository{
		db: db,
	}
}

func (r *storyRepository) Create(ctx context.Context, story *domain.Story) error {
	model := toModel(story)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if isDuplicateKeyError(err) {
			return errors.NewValidationError("story already exists")
		}
		if isTransientError(err) {
			return errors.NewTransientError(err)
		}
		return errors.NewUnexpectedError(err)
	}
	story.ID = model.ID
	return nil
}

func (r *storyRepository) Update(ctx context.Context, story *domain.Story) error {
	model := toModel(story)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		if isDuplicateKeyError(err) {
			return errors.NewValidationError("story already exists")
		}
		if isTransientError(err) {
			return errors.NewTransientError(err)
		}
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func (r *storyRepository) Delete(ctx context.Context, id string) error {
	idUint, _ := strconv.ParseUint(id, 10, 64)
	if err := r.db.WithContext(ctx).Delete(&storyModel{}, uint(idUint)).Error; err != nil {
		if isTransientError(err) {
			return errors.NewTransientError(err)
		}
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func (r *storyRepository) GetByID(ctx context.Context, id string) (*domain.Story, error) {
	var model storyModel
	idUint, _ := strconv.ParseUint(id, 10, 64)
	if err := r.db.WithContext(ctx).First(&model, uint(idUint)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("story", id)
		}
		if err == gorm.ErrInvalidTransaction || err == gorm.ErrRegistered {
			return nil, errors.NewTransientError(err)
		}
		return nil, errors.NewUnexpectedError(err)
	}
	return toDomain(&model), nil
}

func (r *storyRepository) List(ctx context.Context, limit, offset int) ([]*domain.Story, error) {
	var models []*storyModel
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		if isTransientError(err) {
			return nil, errors.NewTransientError(err)
		}
		return nil, errors.NewUnexpectedError(err)
	}

	stories := make([]*domain.Story, len(models))
	for i, model := range models {
		stories[i] = toDomain(model)
	}
	return stories, nil
}

func (r *storyRepository) ListByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*domain.Story, error) {
	var models []*storyModel
	authorIDUint, _ := strconv.ParseUint(authorID, 10, 64)
	err := r.db.WithContext(ctx).
		Where("author_id = ?", uint(authorIDUint)).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		if isTransientError(err) {
			return nil, errors.NewTransientError(err)
		}
		return nil, errors.NewUnexpectedError(err)
	}

	stories := make([]*domain.Story, len(models))
	for i, model := range models {
		stories[i] = toDomain(model)
	}
	return stories, nil
}

// toModel converts domain story to database model
func toModel(story *domain.Story) *storyModel {
	return &storyModel{
		ID:          story.ID,
		Title:       story.Title,
		Content:     story.Content,
		AuthorID:    story.AuthorID,
		CreatedAt:   story.CreatedAt,
		UpdatedAt:   story.UpdatedAt,
		PublishedAt: story.PublishedAt,
	}
}

// toDomain converts database model to domain story
func toDomain(model *storyModel) *domain.Story {
	return &domain.Story{
		ID:          model.ID,
		Title:       model.Title,
		Content:     model.Content,
		AuthorID:    model.AuthorID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		PublishedAt: model.PublishedAt,
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
