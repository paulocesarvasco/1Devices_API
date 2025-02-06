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
	HomePage(w http.ResponseWriter, r *http.Request)
	PatchDevice(w http.ResponseWriter, r *http.Request)
}

func NewHandler(s services.Services) Handler {
	return &handler{service: s}
}

type handler struct {
	service services.Services
}

func (h *handler) HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func (h *handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missed id parameter", http.StatusBadRequest)
		return
	}
	err := h.service.RemoveDeviceByID(id)
	if err != nil && errors.Is(err, constants.ErrorDeviceNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil && errors.Is(err, constants.ErrorDeviceInUse) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var requestPayload resources.Device
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	device, err := h.service.SaveDevice(requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

func (h *handler) SearchDevice(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	id := queryParams.Get("id")
	brand := queryParams.Get("brand")
	state := queryParams.Get("state")
	if id != "" {
		device, err := h.service.SearchDeviceByID(id)
		if err != nil {
			if errors.Is(err, constants.ErrorDeviceNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(device)
		return
	} else if brand != "" {
		devices, err := h.service.FilterDevicesBrand(brand)
		if err != nil {
			if errors.Is(err, constants.ErrorBrandNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(devices)
		return
	} else if state != "" {
		devices, err := h.service.FilterDevicesState(state)
		if err != nil {
			if errors.Is(err, constants.ErrorStateNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(devices)
		return
	}
	devices, err := h.service.ListAllDevices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(devices)
}

func (h *handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, constants.ErrorMissedRequestIDParameter.Error(), http.StatusBadRequest)
		return
	}
	var newDeviceValues resources.Device
	err := json.NewDecoder(r.Body).Decode(&newDeviceValues)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.UpdateDevice(id, newDeviceValues)
	if err != nil && errors.Is(err, constants.ErrorDeviceNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil && errors.Is(err, constants.ErrorDeviceInUse) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) PatchDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, constants.ErrorMissedRequestIDParameter.Error(), http.StatusBadRequest)
		return
	}
	queryParameters := r.URL.Query()
	name := queryParameters.Get("name")
	brand := queryParameters.Get("brand")
	state := queryParameters.Get("state")
	err := h.service.PatchDevice(id, name, brand, state)
	if err != nil && errors.Is(err, constants.ErrorDeviceNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil && errors.Is(err, constants.ErrorDeviceInUse) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
