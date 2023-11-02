package models

import (
	"time"
)

type SignUpInput struct {
	Username        string `json:"username" binding:"required"`
	FirstName       string `json:"firstName" binding:"required"`
	LastName        string `json:"lastName" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	Username     string       `json:"username,omitempty"`
	FirstName    string       `json:"firstName,omitempty"`
	LastName     string       `json:"lastName,omitempty"`
	Email        string       `json:"email,omitempty"`
	Role         string       `json:"role,omitempty"`
	Provider     string       `json:"provider"`
	UserMetadata UserMetadata `json:"userMetadata,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Gists        []Gist       `json:"gists,omitempty"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,min=8"`
}

type ResendVerificationEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type PublicUserProfileResponse struct {
	Username     string       `json:"username,omitempty"`
	FirstName    string       `json:"firstName,omitempty"`
	LastName     string       `json:"lastName,omitempty"`
	UserMetadata UserMetadata `json:"userMetadata,omitempty"`
	Verified     bool         `json:"verified"`
}

type CreateGistRequest struct {
	Private bool   `json:"private"`
	Content string `json:"content" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Title   string `json:"title" binding:"required"`
}

type CommentOnGistRequest struct {
	Content string `json:"content" binding:"required"`
	GistId  string `json:"gistId" binding:"required"`
}

type UpdateUserDetailsRequest struct {
	StatusIcon     string `json:"statusIcon"`
	ProfilePicture string `json:"profilePicture"`
	Location       string `json:"location"`
	Website        string `json:"website"`
	Twitter        string `json:"twitter"`
	Tagline        string `json:"tagline"`
}

type UpdateGistRequest struct {
	Private bool   `json:"private"`
	Content string `json:"content"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	GistId  string `json:"gistId" binding:"required"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type SuccessResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type GitHubClientIdResponse struct {
	ClientId string `json:"client_id"`
}

type AccessCodeResponse struct {
	AccessCode string `json:"access_code"`
}

type BooleanResponse struct {
	Result bool `json:"result"`
}