package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/qahta0/saudi-cert/config"
	"github.com/qahta0/saudi-cert/db"
	_ "github.com/qahta0/saudi-cert/docs"
	"github.com/qahta0/saudi-cert/routes"
)

// @title Saudi CERT
// @version 2.0
// @description This is a documentation for Saudi CERT RESTfulAPI
// @termsOfService http://swagger.io/terms/

// @contact.name @qahta0
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
// @schemes http

func main() {

	// connect to mongodb

	db.ConnectToMongoDB()

	// jobs.RunSaudiCertCrawler(50, 51)
	// jobs.SaudiCertInserter()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// routes

	app.Get("/swagger/*", swagger.HandlerDefault)

	// health checks

	app.Get("/", HealthCheck)
	routes.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	PORT := config.Config("PORT")
	log.Fatal(app.Listen(":" + PORT))
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Health Check
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data":   "Server is up and running 🚀",
		"status": true,
	}
	if err := c.JSON(res); err != nil {
		return err
	}
	return nil
}
