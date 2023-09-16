package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qahta0/saudi-cert/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/api/v1/alerts", controllers.GetAlerts)
	app.Get("/api/v1/alerts/:id", controllers.GetAlertByID)
}
