package main

import (
	"fiber.misetaku.net/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//init route
	routes.Route{
		F: app,
	}.InitRoute()

	app.Listen(":3000")
}
