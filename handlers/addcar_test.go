package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karanjar/cargobackend_fibre_framework.git/config"
	"github.com/stretchr/testify/assert"
)

func TestCreatecar(t *testing.T) {
	config.ConnectDb()

	app := fiber.New(fiber.Config{})
	app.Post("/cars", Createcar)

	body := `
    {
		"name": "Corolla",
		"model":"toyota",
		"year":2021,
		"price":4979307
	}`

	req, _ := http.NewRequest("POST", "/cars", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

}

func TestGetcar(t *testing.T) {
	config.ConnectDb()

	app := fiber.New(fiber.Config{})
	app.Get("/cars/:id", Getcar)

	req, _ := http.NewRequest("GET", "/cars/19", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

}
