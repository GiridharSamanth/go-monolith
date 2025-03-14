package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

// GetStory handles GET /v2.0/stories/:id
func (h *StoryHandler) GetStory(c *gin.Context) {
	storyID := c.Param("id")
	if storyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story ID is required"})
		return
	}

	story, author, err := h.storyService.GetStoryDisplayDetails(c.Request.Context(), storyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"story":  story,
		"author": author,
	})
}
