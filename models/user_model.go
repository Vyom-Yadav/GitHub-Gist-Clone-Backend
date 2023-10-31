package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Username        string `gorm:"type:varchar(255);primary_key"`
	FirstName       string `gorm:"type:varchar(255);not null"`
	Email           string `gorm:"uniqueIndex;not null"`
	Password        string `gorm:"not null"`
	Role            string `gorm:"type:varchar(255);not null"`
	Provider        string `gorm:"not null"`
	PasswordResetAt time.Time
	Verified        bool      `gorm:"not null"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`

	VerificationCode   *string `gorm:"default:null"`
	PasswordResetToken *string `gorm:"default:null"`
	LastName           *string `gorm:"type:varchar(255)"`

	// Required for uniquely identifying users using GitHub OAuth
	GithubUserId *string `gorm:"type:varchar(255);uniqueIndex;default:null"`

	// One-to-One Mapping
	UserMetadata UserMetadata `gorm:"foreignKey:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// One-to-Many Mapping
	Gists []Gist `gorm:"foreignKey:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type UserMetadata struct {
	Username       string `gorm:"type:varchar(255);primary_key"` // Foreign Key
	StatusIcon     string `gorm:"type:varchar(255)"`
	ProfilePicture string `gorm:"not null"`

	Location *string `gorm:"type:varchar(255);default:null"`
	Website  *string `gorm:"type:varchar(255);default:null"`
	Twitter  *string `gorm:"type:varchar(255);default:null"`
	Tagline  *string `gorm:"type:varchar(255);default:null"`

	StarredGistsCount int `gorm:"not null"`
	Followers         int `gorm:"not null"`
	Following         int `gorm:"not null"`
}

type Gist struct {
	Username string `gorm:"type:varchar(255)"` // Foreign Key

	StarCount int `gorm:"not null"`

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

type Follow struct {
	Username   string `gorm:"type:varchar(255);primary_key"`
	FollowedBy string `gorm:"type:varchar(255);primary_key"`
}

type Star struct {
	Username string    `gorm:"type:varchar(255);primary_key"`
	GistID   uuid.UUID `gorm:"type:uuid;primary_key"`
}
