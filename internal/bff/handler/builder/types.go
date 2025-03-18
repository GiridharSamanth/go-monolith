package builder

type StoryResponse struct {
	ID      *uint            `json:"id,omitempty"`
	Title   *string          `json:"title,omitempty"`
	Summary *string          `json:"summary,omitempty"`
	Content *string          `json:"content,omitempty"`
	Tags    []string         `json:"tags,omitempty"`
	Author  *AuthorResponse  `json:"author,omitempty"`
	Reviews []ReviewResponse `json:"reviews,omitempty"`
	Likes   *int             `json:"likes,omitempty"`
}

type AuthorResponse struct {
	ID              *uint   `json:"id,omitempty"`
	Name            *string `json:"name,omitempty"`
	ProfileImageURL *string `json:"profileImageUrl,omitempty"`
	ProfilePageURL  *string `json:"profilePageUrl,omitempty"`
	UserID          *int    `json:"userId,omitempty"`
}

type ReviewResponse struct {
	Rating *int            `json:"rating,omitempty"`
	Review *string         `json:"review,omitempty"`
	User   *AuthorResponse `json:"user,omitempty"`
}

type ResponseStructure map[string]interface{}
