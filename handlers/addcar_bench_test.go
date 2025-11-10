package handlers

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karanjar/cargobackend_fibre_framework.git/config"
)

func BenchmarkGetcar(B *testing.B) {
	config.ConnectDb()

	app := fiber.New(fiber.Config{})
	app.Get("/cars/:id", Getcar)

	req, _ := http.NewRequest("GET", "/cars/19", nil)
	req.Header.Set("Content-Type", "application/json")
	for i := 0; i < B.N; i++ {
		_, _ = app.Test(req)
	}

}
