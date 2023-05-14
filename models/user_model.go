package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Username           string `gorm:"type:varchar(255);primary_key"`
	FirstName          string `gorm:"type:varchar(255);not null"`
	LastName           string `gorm:"type:varchar(255);not null"`
	Email              string `gorm:"uniqueIndex;not null"`
	Password           string `gorm:"not null"`
	Role               string `gorm:"type:varchar(255);not null"`
	Provider           string `gorm:"not null"`
	PasswordResetToken string
	PasswordResetAt    time.Time
	VerificationCode   string
	Verified           bool      `gorm:"not null"`
	CreatedAt          time.Time `gorm:"not null"`
	UpdatedAt          time.Time `gorm:"not null"`

	// One-to-One Mapping
	UserMetadata UserMetadata `gorm:"foreignKey:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// One-to-Many Mapping
	Gists []Gist `gorm:"foreignKey:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type UserMetadata struct {
	Username       string `gorm:"type:varchar(255);primary_key"` // Foreign Key
	StatusIcon     string `gorm:"type:varchar(255)"`
	ProfilePicture string `gorm:"not null"`
	Location       string `gorm:"type:varchar(255)"`
	Website        string `gorm:"type:varchar(255)"`
	Twitter        string `gorm:"type:varchar(255)"`
	Tagline        string `gorm:"type:varchar(255)"`

	// ID's of the gists user has starred
	StarredGists []string `gorm:"type:varchar(255)[]"`

	// Username of people who have followed the user
	Followers []string `gorm:"type:varchar(255)[]"`

	// Username of people who the user has followed
	Following []string `gorm:"type:varchar(255)[]"`
}

type Gist struct {
	Username string `gorm:"type:varchar(255)"` // Foreign Key

	// Username of people who have starred the gist
	Stars []string `gorm:"type:varchar(255)[]"`

	// Username of people who have forked the gist
	Forks []string `gorm:"type:varchar(255)[]"`

	ID          uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Comments    []Comment   `gorm:"foreignKey:GistID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Private     bool        `gorm:"not null"`
	GistContent GistContent `gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// We are hard-coding in logic to make sure name is unique across all gists of a user
	Name string `gorm:"type:varchar(255);not null"`

	Title     string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type GistContent struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"` // Foreign Key
	Content string    `gorm:"type:text;size:10485760;not null"`
}

type Comment struct {
	GistID    uuid.UUID `gorm:"type:uuid; not null"` // Foreign Key
	Username  string    `gorm:"type:varchar(255); not null"`
	Content   string    `gorm:"type:text;size:10485760;not null"`
	CommentID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
