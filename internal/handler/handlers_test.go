package handler

import (
	"1Devices_API/internal/database"
	"1Devices_API/internal/resources"
	"1Devices_API/internal/services"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertDevice(t *testing.T) {
	tt := []struct {
		name           string
		requestPayload any
		expectedCode   int
	}{
		{"Device Created", resources.Device{Name: "foo", Brand: "xPhone", State: "available"}, http.StatusCreated},
		{"Bad Payload", `{"invalid": "json"}`, http.StatusBadRequest},
	}

	for _, tc := range tt {
		h := NewHandler(services.NewService(database.NewSQLiteClient()))
		rawBody, _ := json.Marshal(tc.requestPayload)
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/devices", bytes.NewReader(rawBody))
		rr := httptest.NewRecorder()
		h.RegisterDevice(rr, req)
		assert.Equal(t, tc.expectedCode, rr.Code)
	}
}

func TestSearchSingleDevice(t *testing.T) {
	tt := []struct {
		name            string
		queryParameter  string
		queryValue      string
		expectedCode    int
		expectedPayload any
	}{
		{"Search ID", "id", "1", http.StatusOK, resources.Device{ID: 1, Brand: "xPhone", State: "available", CreationTime: time.Now().Format(time.RFC3339)}},
		{"Device not found", "id", "2", http.StatusNotFound, resources.Device{}},
		{"Invalid ID", "id", "a", http.StatusInternalServerError, resources.Device{}},
	}
	for _, tc := range tt {
		h := NewHandler(services.NewService(database.NewSQLiteClient()))
		rawBody, _ := json.Marshal(resources.Device{ID: 1, Brand: "xPhone", State: "available"})
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/devices", bytes.NewReader(rawBody))
		h.RegisterDevice(httptest.NewRecorder(), req)

		rr := httptest.NewRecorder()
		url := fmt.Sprintf("http://localhost:8080/api/v1/devices?%s=%s", tc.queryParameter, tc.queryValue)
		req, _ = http.NewRequest(http.MethodGet, url, nil)
		h.SearchDevice(rr, req)
		assert.Equal(t, tc.expectedCode, rr.Code)
		var response resources.Device
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, tc.expectedPayload, response, tc.name)
	}
}

func TestSearchMultDevices(t *testing.T) {
	tt := []struct {
		name            string
		queryParameter  string
		queryValue      string
		expectedCode    int
		expectedPayload any
	}{
		{"Retrieve all devices", "", "", http.StatusOK, []resources.Device{
			{ID: 1, Brand: "xPhone", State: "available", CreationTime: time.Now().Format(time.RFC3339)},
			{ID: 2, Brand: "Android", State: "available", CreationTime: time.Now().Format(time.RFC3339)}},
		},
		{"Fetch by brand", "brand", "xPhone", http.StatusOK, []resources.Device{
			{ID: 1, Brand: "xPhone", State: "available", CreationTime: time.Now().Format(time.RFC3339)}}},
		{"Fetch by state", "state", "available", http.StatusOK, []resources.Device{
			{ID: 1, Brand: "xPhone", State: "available", CreationTime: time.Now().Format(time.RFC3339)},
			{ID: 2, Brand: "Android", State: "available", CreationTime: time.Now().Format(time.RFC3339)}},
		},
		{"Device not found", "brand", "MotoX", http.StatusNotFound, []resources.Device{}},
	}
	devicesToRegister := []resources.Device{
		{Brand: "xPhone", State: "available"},
		{Brand: "Android", State: "available"},
	}
	for _, tc := range tt {
		h := NewHandler(services.NewService(database.NewSQLiteClient()))
		for _, dr := range devicesToRegister {
			rawBody, _ := json.Marshal(dr)
			req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/devices", bytes.NewReader(rawBody))
			h.RegisterDevice(httptest.NewRecorder(), req)
		}
		url := fmt.Sprintf("http://localhost:8080/api/v1/devices?%s=%s", tc.queryParameter, tc.queryValue)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		rr := httptest.NewRecorder()
		h.SearchDevice(rr, req)
		response := []resources.Device{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, tc.expectedCode, rr.Code)
		assert.Equal(t, tc.expectedPayload, response, tc.name)
	}
}

func TestDeleteDevice(t *testing.T) {
	tt := []struct {
		name         string
		deviceID     string
		deviceState  resources.State
		expectedCode int
	}{
		{"Delete successful", "1", "available", http.StatusNoContent},
		{"Device not found", "2", "available", http.StatusNotFound},
		{"Device in use", "1", "in-use", http.StatusUnauthorized},
		{"Invalid ID", "a", "available", http.StatusInternalServerError},
	}
	for _, tc := range tt {
		h := NewHandler(services.NewService(database.NewSQLiteClient()))
		rawBody, _ := json.Marshal(resources.Device{ID: 1, Brand: "xPhone", State: tc.deviceState})
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/devices", bytes.NewReader(rawBody))
		h.RegisterDevice(httptest.NewRecorder(), req)

		url := fmt.Sprintf("http://localhost:8080/api/v1/devices?id=%s", tc.deviceID)
		req, _ = http.NewRequest(http.MethodDelete, url, nil)
		rr := httptest.NewRecorder()
		h.DeleteDevice(rr, req)
		assert.Equal(t, tc.expectedCode, rr.Code, tc.name)
	}
}

func TestPutDevice(t *testing.T) {
	tt := []struct {
		name            string
		deviceID        string
		regiesterDevice resources.Device
		requestPayload  resources.Device
		expectedCode    int
		expectedDevice  resources.Device
	}{
		{"Update successful", "1", resources.Device{Name: "iPhone", State: "available"},
			resources.Device{Name: "samsung", State: "available"}, http.StatusOK,
			resources.Device{ID: 1, Name: "samsung", State: "available", CreationTime: time.Now().Format(time.RFC3339)}},
		{"Device not found", "2", resources.Device{Name: "iPhone", State: "available"},
			resources.Device{}, http.StatusNotFound,
			resources.Device{}},
		{"Device in-use", "1", resources.Device{Name: "iPhone", State: "in-use"},
			resources.Device{Name: "iPhone", State: "available"}, http.StatusUnauthorized,
			resources.Device{}},
		{"Invalid ID", "a", resources.Device{}, resources.Device{}, http.StatusInternalServerError, resources.Device{}},
	}
	for _, tc := range tt {
		h := NewHandler(services.NewService(database.NewSQLiteClient()))
		rawBody, _ := json.Marshal(tc.regiesterDevice)
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/devices", bytes.NewReader(rawBody))
		h.RegisterDevice(httptest.NewRecorder(), req)

		url := fmt.Sprintf("http://localhost:8080/api/v1/devices?id=%s", tc.deviceID)
		rawBody, _ = json.Marshal(tc.requestPayload)
		req, _ = http.NewRequest(http.MethodPut, url, bytes.NewBuffer(rawBody))
		rr := httptest.NewRecorder()

		h.UpdateDevice(rr, req)
		assert.Equal(t, tc.expectedCode, rr.Code, tc.name)

		if tc.expectedCode == http.StatusOK {
			h.SearchDevice(rr, req)
			var resDevice resources.Device
			json.NewDecoder(rr.Body).Decode(&resDevice)
			assert.Equal(t, tc.expectedDevice, resDevice, tc.name)
		}
	}
}

func TestPatchDevice(t *testing.T) {
	tt := []struct {
		name           string
		deviceID       string
		deviceToSave   resources.Device
		patchParameter []string
		patchValue     []string
		expectedCode   int
		expectedDevice resources.Device
	}{
		{
			"Update brand", "1", resources.Device{Brand: "Apple", State: "available"},
			[]string{"brand"}, []string{"Samsung"}, http.StatusOK,
			resources.Device{ID: 1, Brand: "Samsung", State: "available", CreationTime: time.Now().Format(time.RFC3339)},
		},
		{
			"Update blocked", "1", resources.Device{Name: "foo", Brand: "Apple", State: "in-use"},
			[]string{"brand"}, []string{"Samsung"}, http.StatusUnauthorized,
			resources.Device{ID: 1, Name: "foo", Brand: "Apple", State: "in-use", CreationTime: time.Now().Format(time.RFC3339)},
		},
	}

	for _, tc := range tt {
		h := NewHandler(services.NewService(database.NewSQLiteClient()))
		rawBody, _ := json.Marshal(tc.deviceToSave)
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/devices", bytes.NewReader(rawBody))
		h.RegisterDevice(httptest.NewRecorder(), req)

		url := fmt.Sprintf("http://localhost:8080/api/v1/devices?id=%s", tc.deviceID)
		for i := range tc.patchParameter {
			url += "&" + tc.patchParameter[i] + "=" + tc.patchValue[i]
		}
		req, _ = http.NewRequest(http.MethodPatch, url, nil)
		rr := httptest.NewRecorder()

		h.PatchDevice(rr, req)
		assert.Equal(t, tc.expectedCode, rr.Code, tc.name)

		rr = httptest.NewRecorder()
		h.SearchDevice(rr, req)
		var device resources.Device
		json.NewDecoder(rr.Body).Decode(&device)
		assert.Equal(t, tc.expectedDevice, device, tc.name)
	}
}
