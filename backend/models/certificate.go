package models

import "gorm.io/gorm"

// Certificate represents the certificate model
type Certificate struct {
	gorm.Model
	Code      string `gorm:"unique"`
	Type      string // Certificate type
    Name      string // User name
	Email     string // User email
	Tags      string // Comma sperarated tags
    Date      string // Assignment date: "11. September 2024"
    Generated bool
}
