package builder

import (
	authorDomain "go-monolith/internal/modules/author/domain"
	storyDomain "go-monolith/internal/modules/story/domain"
)

func BuildStoryResponse(story *storyDomain.Story, author *authorDomain.Author, structure ResponseStructure) StoryResponse {
	resp := StoryResponse{}

	if _, ok := structure["id"]; ok {
		resp.ID = &story.ID
	}
	if _, ok := structure["title"]; ok {
		resp.Title = &story.Title
	}
	if _, ok := structure["content"]; ok {
		resp.Content = &story.Content
	}

	// Handle Author
	if authorStruct, ok := structure["author"].(map[string]interface{}); ok {
		resp.Author = &AuthorResponse{}
		if _, include := authorStruct["name"]; include {
			fullName := author.FirstName + " " + author.LastName // Concatenate the strings directly
			resp.Author.Name = &fullName                         // Assign the address of the concatenated string
		}
		if _, include := authorStruct["profileImageUrl"]; include {
			resp.Author.ProfileImageURL = &author.ProfileImageURL
		}
	}

	// Handle Reviews
	// if reviewsStruct, ok := structure["reviews"].([]interface{}); ok && len(story.Reviews) > 0 {
	// 	for _, review := range story.Reviews {
	// 		reviewResp := ReviewResponse{}
	// 		if reviewFields, valid := reviewsStruct[0].(map[string]interface{}); valid {
	// 			if _, include := reviewFields["rating"]; include {
	// 				reviewResp.Rating = &review.Rating
	// 			}
	// 			if _, include := reviewFields["review"]; include {
	// 				reviewResp.Review = &review.Review
	// 			}
	// 			if userStruct, ok := reviewFields["user"].(map[string]interface{}); ok {
	// 				reviewResp.User = &UserResponse{}
	// 				if _, include := userStruct["name"]; include {
	// 					reviewResp.User.Name = &review.User.Name
	// 				}
	// 				if _, include := userStruct["profileImageUrl"]; include {
	// 					reviewResp.User.ProfileImageURL = &review.User.ProfileImageURL
	// 				}
	// 			}
	// 		}
	// 		resp.Reviews = append(resp.Reviews, reviewResp)
	// 	}
	// }

	return resp
}
