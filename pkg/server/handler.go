package server

import "github.com/gofiber/fiber/v2"

type APIServerHandler struct {
	GoRestfulApp *fiber.App
}

func NewAPIServerHandler() *APIServerHandler {
	gorestfulApp := fiber.New()
	
	return &APIServerHandler{
		GoRestfulApp: gorestfulApp,
	}
}