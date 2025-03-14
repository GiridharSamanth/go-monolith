package dto

// CreateAuthorRequest represents the request body for creating an author
type CreateAuthorRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Bio   string `json:"bio" binding:"required"`
}

// UpdateAuthorRequest represents the request body for updating an author
type UpdateAuthorRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Bio   string `json:"bio" binding:"required"`
}
