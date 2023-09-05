package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
)

type APIServerHandler struct {
	GoRestfulApp *fiber.App
	NonGoRestfulMux *mux.Router
}

func NewAPIServerHandler() *APIServerHandler {
	gorestfulApp := fiber.New()
	
	return &APIServerHandler{
		GoRestfulApp: gorestfulApp,
		NonGoRestfulMux: mux.NewRouter(),
	}
}