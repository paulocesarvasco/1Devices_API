package server

import (
	"1Devices_API/internal/handler"
	"1Devices_API/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func(c chan os.Signal) {
		<-c
		log.Println("Shutting down server")
		os.Exit(0)
	}(c)

	handlers := handler.NewHandler()
	r := chi.NewMux()
	router.SetRoutes(r, handlers)
	http.ListenAndServe("localhost:8080", r)
}
