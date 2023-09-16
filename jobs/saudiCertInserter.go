package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/qahta0/saudi-cert/db"
	"github.com/qahta0/saudi-cert/models"
	"go.mongodb.org/mongo-driver/bson"
)

func SaudiCertInserter() {
	rawData, err := os.ReadFile("./alerts/alerts.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	var records []map[string]interface{}
	err = json.Unmarshal(rawData, &records)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	var alerts []models.Alert
	for _, record := range records {
		alert, err := processRecord(record)
		if err != nil {
			log.Printf("Error processing record: %v", err)
		} else {
			alerts = append(alerts, alert)
		}
	}
	storeAlertsInMongoDB(alerts)
}

func storeAlertsInMongoDB(alerts []models.Alert) {
	ctx := context.TODO()
	collection := db.GetMongoClient().Database("saudi_cert").Collection("alerts")
	for _, alert := range alerts {
		filter := bson.M{"original_link": alert.OriginalLink}
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			log.Printf("Error checking for existing alert: %v", err)
			continue
		}
		if count == 0 {
			_, err := collection.InsertOne(ctx, alert)
			if err != nil {
				log.Printf("Error inserting alert into MongoDB: %v", err)
			}
		}
	}
	fmt.Println("Alerts successfully inserted into MongoDB.")
}

func processRecord(record map[string]interface{}) (models.Alert, error) {
	var alert models.Alert
	if severityLevel, ok := record["severity_level"].(string); ok {
		alert.SeverityLevel = strings.TrimSpace(severityLevel)
	}
	if name, ok := record["name"].(string); ok {
		alert.Name = name
	}
	if imgUrl, ok := record["image_url"].(string); ok {
		alert.ImageURL = imgUrl
	}
	if origLink, ok := record["original_link"].(string); ok {
		alert.OriginalLink = origLink
	}
	if details, ok := record["details"].(map[string]interface{}); ok {
		if affectedProducts, ok := details["affected_products"].(string); ok {
			for _, product := range strings.Split(affectedProducts, "#") {
				alert.Details.AffectedProducts = append(alert.Details.AffectedProducts, models.AffectedProduct{Name: strings.TrimSpace(product)})
			}
		}
		if threatList, ok := details["threat_list"].(string); ok {
			threatNames := strings.Split(threatList, "#")
			var threats []models.Threat
			for _, threatName := range threatNames {
				threats = append(threats, models.Threat{Name: strings.TrimSpace(threatName)})
			}
			alert.Details.ThreatList = threats
		}
		if recommendationLinks, ok := details["recommendation_links"].(string); ok {
			links := strings.Split(recommendationLinks, "#")
			var recommendations []models.Recommendation
			for _, link := range links {
				recommendations = append(recommendations, models.Recommendation{Link: strings.TrimSpace(link)})
			}
			alert.Details.Recommendations = recommendations
		}
		if warningNumberStr, ok := details["warning_number"].(string); ok {
			alert.Details.WarningNumber = warningNumberStr
		}
		if warningDate, ok := details["warning_date"].(string); ok {
			alert.Details.WarningDate = warningDate
		}
	}
	jsonData, _ := json.Marshal(record)
	json.Unmarshal(jsonData, &alert)
	return alert, nil
}
