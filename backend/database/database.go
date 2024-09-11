package database

import (
	"am-certificates/models"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(file string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
	DB.AutoMigrate(&models.Certificate{})
}

func CleanDatabase() error {
	err := DB.Where("type = ? AND generated = ?", "", false).Delete(&models.Certificate{}).Error
	return err
}

// GetTotalEntries returns the total number of entries in the Certificate table
func GetTotalEntries() int64 {
	var count int64
	if err := DB.Model(&models.Certificate{}).Count(&count).Error; err != nil {
		log.Println("Could not get total number of database entries", err)
		return 0
	}
	return count
}

// GetUnusedEntries returns the number of certificates that have not been used (Generated = false)
func GetUnusedEntries() int64 {
	var count int64
	if err := DB.Model(&models.Certificate{}).Where("generated = ?", false).Count(&count).Error; err != nil {
		return 0
	}
	return count
}

// GetEntriesByType returns a map of certificate types and their counts
func GetEntriesByType() map[string]int64 {
	var results []struct {
		Type  string
		Count int64
	}
	typeStats := make(map[string]int64)

	// Query to count entries grouped by type
	if err := DB.Model(&models.Certificate{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&results).Error; err != nil {
		return nil
	}

	// Convert query result to a map
	for _, result := range results {
		typeStats[result.Type] = result.Count
	}

	return typeStats
}

func AssignEntry(certificate models.Certificate) error {
	certificate.Generated = true
	certificate.Date = time.Now().Format("01.02.2006")
	return DB.Save(&certificate).Error
}

func GetEntry(code string, used bool) (models.Certificate, error) {

	var certificate models.Certificate
	err := DB.Where("code = ? AND generated = ?", code, used).First(&certificate).Error
	return certificate, err
}

func DeleteEntry(code string) error {
	// Find the certificate by code and delete it
	err := DB.Where("code = ?", code).Delete(&models.Certificate{}).Error
	return err
}
