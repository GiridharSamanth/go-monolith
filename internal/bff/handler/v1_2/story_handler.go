package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-monolith/internal/bff/handler/builder"
	"go-monolith/internal/bff/service"
)

type StoryHandler struct {
	storyService *service.StoryService
}

var storyHandler *StoryHandler

func NewStoryHandler(ss *service.StoryService) *StoryHandler {
	if storyHandler == nil {
		storyHandler = &StoryHandler{
			storyService: ss,
		}
	}
	return storyHandler
}

// GetStory handles GET /v1.2/stories?id=123
func (h *StoryHandler) GetStory(c *gin.Context) {
	storyID := c.Query("id")
	fmt.Println("GetStory", "storyID", storyID)
	if storyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story ID is required"})
		return
	}
	story, author, err := h.storyService.GetStoryDisplayDetails(c.Request.Context(), storyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseStructure := builder.ResponseStructure{
		"id":    true,
		"title": true,
		"author": map[string]interface{}{
			"name":            true,
			"profileImageUrl": true,
		},
	}
	storyResponse := builder.BuildStoryResponse(story, author, responseStructure)

	c.JSON(http.StatusOK, storyResponse)
}
