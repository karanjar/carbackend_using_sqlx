package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func SecureHeaders(c *fiber.Ctx) error {

	c.Response().Header.Add("Content-Securty-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'")
	c.Response().Header.Add("X-XSS-Protection", "1; mode=block")
	c.Response().Header.Add("X-Frame-Options", "deny")
	c.Response().Header.Add("X-Content-Type", "application/x-www-form-urlencoded")

	return c.Next()
}
