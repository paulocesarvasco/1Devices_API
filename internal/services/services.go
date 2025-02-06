package services

import (
	"1Devices_API/internal/constants"
	"1Devices_API/internal/database"
	"1Devices_API/internal/resources"
	"strconv"
	"time"
)

type Services interface {
	SaveDevice(device resources.Device) (resources.Device, error)
	SearchDeviceByID(id string) (resources.Device, error)
	ListAllDevices() ([]resources.Device, error)
	FilterDevicesBrand(brand string) ([]resources.Device, error)
	FilterDevicesState(state string) ([]resources.Device, error)
	RemoveDeviceByID(id string) error
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

func (s *service) SearchDeviceByID(id string) (resources.Device, error) {
	idValue, err := strconv.Atoi(id)
	if err != nil {
		return resources.Device{}, constants.ErrorInvalidIDFormat
	}
	device, err := s.db.SelectDevice(idValue)
	if err != nil {
		return resources.Device{}, err
	}
	return device, nil
}

func (s *service) ListAllDevices() ([]resources.Device, error) {
	return s.db.FetchAllDevices()
}

func (s *service) FilterDevicesBrand(brand string) ([]resources.Device, error) {
	devices, err := s.db.FetchDevicesByBrand(brand)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, constants.ErrorBrandNotFound
	}
	return devices, nil
}

func (s *service) FilterDevicesState(state string) ([]resources.Device, error) {
	devices, err := s.db.FetchDevicesByState(state)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, constants.ErrorDeviceNotFound
	}
	return devices, nil
}

func (s *service) RemoveDeviceByID(id string) error {
	idValue, err := strconv.Atoi(id)
	if err != nil {
		return constants.ErrorInvalidRequestParameter
	}
	return s.db.RemoveDevice(idValue)
}
