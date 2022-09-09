package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/gregorybai/strm/internal/handlers"
	"github.com/gregorybai/strm/pkg/helpers"
)

func main() {
	app := fiber.New()

	app.Static("/public/", "./internal/public")

	app.Get("/hello", handlers.Hello)

	// TODO: gofiber/websocket
	// app.Use("/room/*", func(c *fiber.Ctx) error {
	// 	// IsWebSocketUpgrade returns true if the client
	// 	// requested upgrade to the WebSocket protocol.
	// 	if websocket.IsWebSocketUpgrade(c) {
	// 		c.Locals("allowed", true) // ?
	// 		return c.Next()
	// 	}
	// 	return fiber.ErrUpgradeRequired
	// })

	app.Get("/room/create", handlers.CreateRoom)
	app.Get("/room/:id", websocket.New(handlers.JoinRoom))

	helpers.Must(app.Listen(":8000"))
}
