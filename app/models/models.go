package models

import "gorm.io/gorm"

// Represents a model for video data in the database
type Video struct {
	gorm.Model         // Embeds GORM's default model fields (ID, CreatedAt, UpdatedAt, DeletedAt)
	ID          uint   // Custom field for video ID
	FilePath    string `gorm:"default:null"` // File path for the video. Default value is null
	Title       string // Title of the video
	Description string // Description of the video
}
