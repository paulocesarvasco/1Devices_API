package database

import (
	"1Devices_API/internal/constants"
	"1Devices_API/internal/resources"
	"errors"
	"log"

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

func (c *sqliteClient) SelectDevice(id int) (resources.Device, error) {
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

func (c *sqliteClient) FetchAllDevices() ([]resources.Device, error) {
	var devices []resources.Device
	result := c.db.Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (c *sqliteClient) FetchDevicesByBrand(brand string) ([]resources.Device, error) {
	var devices []resources.Device
	result := c.db.Where("brand = ?", brand).Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (c *sqliteClient) FetchDevicesByState(state string) ([]resources.Device, error) {
	var devices []resources.Device
	result := c.db.Where("state = ?", state).Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (c *sqliteClient) RemoveDevice(id int) error {
	result := c.db.Delete(&resources.Device{}, id)
	if result.RowsAffected == 0 {
		return constants.ErrorDeviceNotFound
	}
	return result.Error
}
