package services

import (
	"1Devices_API/internal/database"
	"1Devices_API/internal/resources"
	"time"
)

type Services interface {
	SaveDevice(device resources.Device) (resources.Device, error)
}

type service struct {
	db database.Client
}

func NewService(db database.Client) Services {
	return &service{db: db}
}

func (s *service) SaveDevice(device resources.Device) (resources.Device, error) {
	device.CreationTime = time.Now().Format(time.RFC3339)
	device, err := s.db.InsertDevice(device)
	if err != nil {
		return resources.Device{}, err
	}
	return device, nil
}
