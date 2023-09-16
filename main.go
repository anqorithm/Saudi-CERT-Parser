package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/qahta0/saudi-cert/db"
	"github.com/qahta0/saudi-cert/jobs"
	"github.com/qahta0/saudi-cert/routes"
)

func main() {
	db.ConnectToMongoDB()
	jobs.RunSaudiCertCrawler(1, 50)
	jobs.SaudiCertInserter()
	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
