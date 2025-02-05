package handler

import (
	"1Devices_API/internal/resources"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertDevice(t *testing.T) {
	tt := []struct {
		name           string
		requestPayload any
		expectedCode   int
	}{
		{"Device Created", resources.Device{}, http.StatusCreated},
		{"Bad Payload", `{"invalid": "json"}`, http.StatusBadRequest},
	}

	for _, tc := range tt {
		h := NewHandler()
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
		{"Search ID", "id", "123", http.StatusOK, resources.Device{}},
		{"Device not found", "id", "124", http.StatusNotFound, resources.Device{}},
	}
	for _, tc := range tt {
		h := NewHandler()
		url := fmt.Sprintf("http://localhost:8080/api/v1/devices?%s=%s", tc.queryParameter, tc.queryValue)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		rr := httptest.NewRecorder()
		h.SearchDevice(rr, req)
		assert.Equal(t, tc.expectedCode, rr.Code)
		var response resources.Device
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, tc.expectedPayload, response)
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
		{"Retrieve all devices", "", "", http.StatusOK, []resources.Device{}},
		{"Fetch by brand", "brand", "xPhone", http.StatusOK, []resources.Device{}},
		{"Fetch by state", "state", "available", http.StatusOK, []resources.Device{}},
		{"Device not found", "brand", "Android", http.StatusNotFound, nil},
	}
	for _, tc := range tt {
		h := NewHandler()
		url := fmt.Sprintf("http://localhost:8080/api/v1/devices?%s=%s", tc.queryParameter, tc.queryValue)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		rr := httptest.NewRecorder()
		h.SearchDevice(rr, req)
		assert.Equal(t, tc.expectedCode, rr.Code)
		var response resources.Device
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, tc.expectedPayload, response)
	}
}
