package models

import (
	"time"
)

// Database model for all providers

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Username    string `gorm:"uniqueIndex"`
	DisplayName string
	Bio         string     `gorm:"type:text"`
	Avatar      UserMedia  `gorm:"embedded;embeddedPrefix:avatar_"`
	Banner      UserMedia  `gorm:"embedded;embeddedPrefix:banner_"`
	Links       []UserLink `gorm:"foreignKey:UserID"`
	Posts       []UserPost `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserPost struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        uint   `gorm:"index"`
	Slug          string `gorm:"uniqueIndex"`
	Content       string `gorm:"type:text"`
	Likes         int
	IsPreview     bool
	PostCreatedAt time.Time
	PostUpdatedAt time.Time
	Media         []UserPostMedia `gorm:"foreignKey:UserPostID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserPostMedia struct {
	ID         uint `gorm:"primaryKey"`
	UserPostID uint
	UserMedia  `gorm:"embedded"`
	Type       int
	Width      int
	Height     int
	Duration   *int
}

type UserLink struct {
	ID               uint `gorm:"primaryKey"`
	UserID           uint `gorm:"index"`
	Username         string
	Website          string `gorm:"index"`
	URL              string
	UniqueConstraint string `gorm:"uniqueIndex:user_social_unique"`
}

type UserMedia struct {
	Filename string
	URL      string `gorm:"type:text"`
}
