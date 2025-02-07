package database

import (
	"1Devices_API/internal/constants"
	"1Devices_API/internal/resources"
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Client interface {
	InsertDevice(device resources.Device) (resources.Device, error)
	SelectDevice(id int) (resources.Device, error)
	FetchAllDevices() ([]resources.Device, error)
	FetchDevicesByState(state string) ([]resources.Device, error)
	FetchDevicesByBrand(brand string) ([]resources.Device, error)
	RemoveDevice(id int) error
	UpdateDevice(currentValues resources.Device, newValues resources.Device) error
	RunMigrations()
}

type dbClient struct {
	db *gorm.DB
}

func (c *dbClient) RunMigrations() {
	err := c.db.AutoMigrate(&resources.Device{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration completed successfully!")
}

func NewPostgresClient() Client {
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return &dbClient{db: db}
}

func NewSQLiteClient() Client {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&resources.Device{})

	return &dbClient{db: db}
}

func (c *dbClient) InsertDevice(device resources.Device) (resources.Device, error) {
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return resources.Device{}, tx.Error
	}

	if err := tx.Create(&device).Error; err != nil {
		tx.Rollback()
		return resources.Device{}, err
	}
	err := tx.Commit().Error
	if err != nil {
		return resources.Device{}, err
	}
	return device, nil
}

func (c *dbClient) SelectDevice(id int) (resources.Device, error) {
	var device resources.Device
	result := c.db.First(&device, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return resources.Device{}, constants.ErrorDeviceNotFound
	}
	if result.Error != nil {
		return resources.Device{}, result.Error
	}
	return device, nil
}

func (c *dbClient) FetchAllDevices() ([]resources.Device, error) {
	var devices []resources.Device
	result := c.db.Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (c *dbClient) FetchDevicesByBrand(brand string) ([]resources.Device, error) {
	var devices []resources.Device
	result := c.db.Where("brand = ?", brand).Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (c *dbClient) FetchDevicesByState(state string) ([]resources.Device, error) {
	var devices []resources.Device
	result := c.db.Where("state = ?", state).Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (c *dbClient) RemoveDevice(id int) error {
	var device resources.Device
	result := c.db.Where("id = ?", id).First(&device)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return constants.ErrorDeviceNotFound
	}
	if device.State == "in-use" {
		return constants.ErrorDeviceInUse
	}
	deleteResult := c.db.Delete(&device)
	if deleteResult.Error != nil {
		return deleteResult.Error
	}
	return nil
}

func (c *dbClient) UpdateDevice(currentValues resources.Device, newValues resources.Device) error {
	if err := c.db.Model(&currentValues).Updates(newValues).Error; err != nil {
		return err
	}
	return nil
}
