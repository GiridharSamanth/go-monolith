package handler

import (
	v1_2 "go-monolith/internal/bff/handler/v1_2"
	v2_0 "go-monolith/internal/bff/handler/v2_0"
	"go-monolith/internal/bff/service"
)

// Handlers struct to hold all handlers
type Handlers struct {
	V1_2StoryHandler *v1_2.StoryHandler
	V2_0StoryHandler *v2_0.StoryHandler
}

// NewHandlers initializes and returns all handlers
func NewHandlers(storyService *service.StoryService) *Handlers {
	return &Handlers{
		V1_2StoryHandler: v1_2.NewStoryHandler(storyService),
		V2_0StoryHandler: v2_0.NewStoryHandler(storyService),
	}
}
