package handler

import (
	"1Devices_API/internal/constants"
	"1Devices_API/internal/resources"
	"1Devices_API/internal/services"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler interface {
	DeleteDevice(w http.ResponseWriter, r *http.Request)
	RegisterDevice(w http.ResponseWriter, r *http.Request)
	SearchDevice(w http.ResponseWriter, r *http.Request)
	UpdateDevice(w http.ResponseWriter, r *http.Request)
}

func NewHandler(s services.Services) Handler {
	return &handler{service: s}
}

type handler struct {
	service services.Services
}

func (h *handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {}

func (h *handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var requestPayload resources.Device
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	device, err := h.service.SaveDevice(requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

func (h *handler) SearchDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		devices, err := h.service.ListAllDevices()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(devices)
		return
	}
	device, err := h.service.SearchDeviceByID(id)
	if err != nil {
		if errors.Is(err, constants.ErrorDeviceNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(device)
}

func (h *handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {}
