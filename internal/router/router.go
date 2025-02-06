package router

import (
	"1Devices_API/internal/handler"

	"github.com/go-chi/chi/v5"
)

func SetRoutes(r chi.Router, h handler.Handler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", h.HomePage)
		r.Get("/devices", h.SearchDevice)
		r.Post("/devices", h.RegisterDevice)
		r.Delete("/devices", h.DeleteDevice)
		r.Put("/devices", h.UpdateDevice)
		r.Patch("/devices", h.PatchDevice)
	})
}
