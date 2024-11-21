package http

import (
	"digibank/delivery/http/router"
	"digibank/infrastructure"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
)

func NewHttpDelivery(dir string) *fiber.App {
	app := fiber.New()
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:3001, http://localhost:3000",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	var c *fiber.Ctx

	db := infrastructure.NewDatabaseConnect(dir)

	router.NewRouter(app, db, c)

	return app
}
