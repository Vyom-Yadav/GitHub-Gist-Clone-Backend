package models

import (
	"time"

	"github.com/google/uuid"
)

type ErrorResponseWrapper struct {
	Error ErrorResponse `json:"error"`
}

type SuccessResponseWrapper struct {
	Success SuccessResponse `json:"success"`
}

type GitHubClientIdResponseWrapper struct {
	GitHubClientId GitHubClientIdResponse `json:"data"`
}

type AccessCodeResponseWrapper struct {
	AccessCode AccessCodeResponse `json:"data"`
}

// GistWithoutComments : Manually sync with Gist
type GistWithoutComments struct {
	Username string

	StarCount int

	ID          uuid.UUID
	Private     bool
	GistContent GistContent

	// We are hard-coding in logic to make sure name is unique across all gists of a user
	Name string

	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GistWithoutCommentsWrapper struct {
	Gist GistWithoutComments `json:"data"`
}

type StringArrayWrapper struct {
	StringArray []string `json:"data"`
}

type CommentArrayWrapper struct {
	Comments []Comment `json:"data"`
}

type UserResponseWrapper struct {
	UserResponse UserResponse `json:"data"`
}

type PublicUserProfileResponseWrapper struct {
	PublicUserProfileResponse PublicUserProfileResponse `json:"data"`
}

type GistWithoutCommentsArrayWrapper struct {
	Gists []GistWithoutComments `json:"data"`
}

type UUIDArrayWrapper struct {
	UUIDArray []uuid.UUID `json:"data"`
}

type CommentWrapper struct {
	Comment Comment `json:"data"`
}

type UserMetadataWrapper struct {
	UserMetadata UserMetadata `json:"data"`
}

type BooleanResponseWrapper struct {
	BooleanResponse BooleanResponse `json:"data"`
}