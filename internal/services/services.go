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
	FilterDevicesState(state resources.State) ([]resources.Device, error)
	RemoveDeviceByID(id string) error
	UpdateDevice(id string, newValues resources.Device) error
	PatchDevice(id string, name string, brand string, state resources.State) error
}

type service struct {
	db database.Client
}

func NewService(db database.Client) Services {
	return &service{db: db}
}

func (s *service) SaveDevice(device resources.Device) (resources.Device, error) {
	if device.State != constants.AVAILABLE && device.State != constants.INACTIVE && device.State != constants.IN_USE {
		return resources.Device{}, constants.ErrorInvalidDeviceState
	}
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

func (s *service) FilterDevicesState(state resources.State) ([]resources.Device, error) {
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
	device, err := s.SearchDeviceByID(id)
	if err != nil {
		return err
	}
	if device.State == constants.IN_USE {
		return constants.ErrorDeviceInUse
	}
	return s.db.RemoveDevice(device.ID)
}

func (s *service) UpdateDevice(id string, newValues resources.Device) error {
	idValue, err := strconv.Atoi(id)
	if err != nil {
		return constants.ErrorInvalidRequestParameter
	}
	currentValues, err := s.db.SelectDevice(idValue)
	if err != nil {
		return err
	}
	if currentValues.State == constants.IN_USE {
		return constants.ErrorDeviceInUse
	}
	newValues.CreationTime = currentValues.CreationTime
	return s.db.UpdateDevice(currentValues, newValues)
}

func (s *service) PatchDevice(id string, name string, brand string, state resources.State) error {
	idValue, err := strconv.Atoi(id)
	if err != nil {
		return constants.ErrorInvalidRequestParameter
	}
	currentValues, err := s.db.SelectDevice(idValue)
	if err != nil {
		return err
	}
	if (name != "" || brand != "") && currentValues.State == constants.IN_USE {
		return constants.ErrorDeviceInUse
	}
	newValues := currentValues
	if name != "" {
		newValues.Name = name
	}
	if brand != "" {
		newValues.Brand = brand
	}
	if state != "" {
		newValues.State = state
	}
	return s.db.UpdateDevice(currentValues, newValues)
}
