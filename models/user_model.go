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
	UserMetadata UserMetadata `gorm:"foreignKey:UserName;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// One-to-Many Mapping
	Gists []Gist `gorm:"foreignKey:UserName;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type UserMetadata struct {
	UserName       string   `gorm:"type:varchar(255)"` // Foreign Key
	StatusIcon     string   `gorm:"type:varchar(255)"`
	ProfilePicture string   `gorm:"not null"`
	Location       string   `gorm:"type:varchar(255)"`
	Website        string   `gorm:"type:varchar(255)"`
	Twitter        string   `gorm:"type:varchar(255)"`
	Tagline        string   `gorm:"type:varchar(255)"`
	Stars          []string `gorm:"type:varchar(255)[]"`
	Followers      []string `gorm:"type:varchar(255)[]"`
	Following      []string `gorm:"type:varchar(255)[]"`
}

type Gist struct {
	UserName  string    `gorm:"type:varchar(255)"` // Foreign Key
	Stars     []string  `gorm:"type:varchar(255)[]"`
	Forks     []string  `gorm:"type:varchar(255)[]"`
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Comments  []Comment `gorm:"foreignKey:GistID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Private   bool      `gorm:"not null"`
	Content   []byte    `gorm:"type:text;size:10485760;not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Link      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type Comment struct {
	GistID    uuid.UUID `gorm:"type:uuid; not null"` // Foreign Key
	Username  string    `gorm:"type:varchar(255); not null"`
	Content   []byte    `gorm:"type:text;size:10485760;not null"`
	CommentID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
