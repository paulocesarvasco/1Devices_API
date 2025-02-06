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
			{ID: 1, Brand: "xPhone", State: "available"}}},
		{"Fetch by state", "state", "available", http.StatusOK, []resources.Device{
			{ID: 1, Brand: "xPhone", State: "available"},
			{ID: 2, Brand: "Android", State: "available"}},
		},
		{"Device not found", "brand", "Android", http.StatusNotFound, []resources.Device{}},
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
		var response []resources.Device
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, tc.expectedCode, rr.Code)
		assert.Equal(t, tc.expectedPayload, response, tc.name)
	}
}

func TestDeleteDevice(t *testing.T) {
	tt := []struct {
		name           string
		queryParameter string
		queryValue     string
		expectedCode   int
	}{
		{"Delete successful", "id", "123", http.StatusNoContent},
		{"Device not found", "id", "124", http.StatusNotFound},
	}
	for _, tc := range tt {
		h := NewHandler(services.NewService(database.NewSQLiteClient()))
		rawBody, _ := json.Marshal(resources.Device{ID: 1, Brand: "xPhone", State: "available"})
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/devices", bytes.NewReader(rawBody))
		h.RegisterDevice(httptest.NewRecorder(), req)

		url := fmt.Sprintf("http://localhost:8080/api/v1/devices?%s=%s", tc.queryParameter, tc.queryValue)
		req, _ = http.NewRequest(http.MethodDelete, url, nil)
		rr := httptest.NewRecorder()
		h.DeleteDevice(rr, req)
		assert.Equal(t, tc.expectedCode, rr.Code)
	}
}
