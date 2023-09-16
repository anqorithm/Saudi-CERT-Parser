package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/qahta0/saudi-cert/db"
	"github.com/qahta0/saudi-cert/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAlerts godoc
// @Summary Get a list of alerts
// @Description Fetches alerts with optional limit
// @Tags Alerts
// @Accept json
// @Produce json
// @Param limit query int false "Limit for the number of alerts returned, defaults to 10" default(10)
// @Success 200 {array} models.Alert "Successfully retrieved alerts"
// @Failure 400 {object} map[string]interface{} "Invalid limit value"
// @Failure 500 {object} map[string]interface{} "Error fetching or decoding the alerts"
// @Router /api/v1/alerts [get]
func GetAlerts(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid limit value", "data": nil})
	}
	collection := db.GetMongoClient().Database("saudi_cert").Collection("alerts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetLimit(int64(limit))
	cur, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error fetching the alerts", "data": nil})
	}
	defer cur.Close(ctx)
	var alerts []models.Alert
	if err = cur.All(ctx, &alerts); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error decoding the alerts", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Alerts fetched successfully", "data": alerts})
}

// GetAlertByID godoc
// @Summary Get a specific alert by ID
// @Description Fetches an alert based on the provided ID
// @Tags Alerts
// @Accept json
// @Produce json
// @Param id path string true "ID of the alert to fetch"
// @Success 200 {object} models.Alert "Successfully retrieved the alert"
// @Failure 400 {object} map[string]interface{} "Invalid ID format or ID is required"
// @Failure 404 {object} map[string]interface{} "Alert not found"
// @Failure 500 {object} map[string]interface{} "Error fetching the alert"
// @Router /api/v1/alerts/{id} [get]
func GetAlertByID(c *fiber.Ctx) error {
	alertID := c.Params("id")
	if alertID == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "ID is required", "data": nil})
	}
	objectID, err := primitive.ObjectIDFromHex(alertID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid ID format", "data": nil})
	}
	collection := db.GetMongoClient().Database("saudi_cert").Collection("alerts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var alert models.Alert
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&alert)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Alert not found", "data": nil})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error fetching the alert", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Alert fetched successfully", "data": alert})
}
