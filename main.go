package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/karanjar/cargobackend_fibre_framework.git/config"
	_ "github.com/karanjar/cargobackend_fibre_framework.git/docs"
	"github.com/karanjar/cargobackend_fibre_framework.git/handlers"
)

// @title Car Backend API
// @version 1.0
// @description API documentation for the Cargo Car Backend built with Fiber.
// @host localhost:8080
// @BasePath /

func main() {
	config.ConnectDb()
	config.ConnectCache()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	//app.Use(middleware.SecureHeaders)
	//app.Use(basicauth.New(basicauth.Config{
	//	Users: map[string]string{
	//		"admin":   "12345",
	//		"manager": "dustin",
	//		"john":    "doe",
	//	},
	//	Unauthorized: func(c *fiber.Ctx) error {
	//		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	//			"error": "user is unauthorized",
	//		})
	//	},
	//}))
	//app.Use(etag.New())
	app.Use("/swagger/*", swagger.HandlerDefault)
	//POST /cars
	//GET /cars/:id
	//DELETE /cars/:id

	app.Post("/cars", handlers.Createcar)
	app.Get("/cars/:id", handlers.Getcar)
	app.Delete("/cars/:id", handlers.Deletecar)
	app.Put("/cars/:id", handlers.Updatecar)

	fmt.Println("HTTP server listenning on port 8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Couldn't listen to the port 8080,error: %v", err)
	}
}
