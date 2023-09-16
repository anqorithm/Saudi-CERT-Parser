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
