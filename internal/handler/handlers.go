package handler

import "net/http"

type Handler interface {
	DeleteDevice(w http.ResponseWriter, r *http.Request)
	RegisterDevice(w http.ResponseWriter, r *http.Request)
	SearchDevice(w http.ResponseWriter, r *http.Request)
	UpdateDevice(w http.ResponseWriter, r *http.Request)
}

func NewHandler() Handler {
	return &handler{}
}

type handler struct{}

func (h *handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {}

func (h *handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {}

func (h *handler) SearchDevice(w http.ResponseWriter, r *http.Request) {}

func (h *handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {}
