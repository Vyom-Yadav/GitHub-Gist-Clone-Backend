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
