package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/qahta0/saudi-cert/config"
	"github.com/qahta0/saudi-cert/db"
	"github.com/qahta0/saudi-cert/jobs"
	"github.com/qahta0/saudi-cert/routes"
)

func main() {
	db.ConnectToMongoDB()
	jobs.RunSaudiCertCrawler(50, 51)
	jobs.SaudiCertInserter()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	routes.SetupRoutes(app)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
	PORT := config.Config("PORT")
	log.Fatal(app.Listen(":" + PORT))
}
