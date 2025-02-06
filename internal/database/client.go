package database

import (
	"1Devices_API/internal/resources"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Client interface {
	InsertDevice(device resources.Device) (resources.Device, error)
}

type sqliteClient struct {
	db *gorm.DB
}

func NewSQLiteClient() Client {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&resources.Device{})

	return &sqliteClient{db: db}
}

func (c *sqliteClient) InsertDevice(device resources.Device) (resources.Device, error) {
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return resources.Device{}, tx.Error
	}
	device.CreationTime = time.Now().Format(time.RFC3339)

	if err := tx.Create(device).Error; err != nil {
		tx.Rollback()
		return resources.Device{}, err
	}
	err := tx.Commit().Error
	if err != nil {
		return resources.Device{}, err
	}
	return device, nil
}
