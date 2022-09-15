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

	app.Get("/room/create", handlers.CreateRoom)
	app.Get("/room/:id", websocket.New(handlers.JoinRoom))
	// How to log ws events ?
	app.Get("/rtc/ws", websocket.New(handlers.InitWebRTC))

	// helpers.Must(app.Listen(":8000"))
	go helpers.Must(app.Listen(":8000"))
}
