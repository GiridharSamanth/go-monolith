package data

import (
	"context"

	storydomain "go-monolith/internal/modules/story/domain"
	storyModuleService "go-monolith/internal/modules/story/service"
)

type StoryProvider struct {
	storyService *storyModuleService.StoryService
}

func NewStoryProvider(sq *storyModuleService.StoryService) *StoryProvider {
	return &StoryProvider{
		storyService: sq,
	}
}

func (p *StoryProvider) GetStory(ctx context.Context, id string) (*storydomain.Story, error) {
	return p.storyService.GetByID(ctx, id)
}
