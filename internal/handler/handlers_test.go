package handler

import (
	"1Devices_API/internal/resources"
	"bytes"
	"encoding/json"
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
