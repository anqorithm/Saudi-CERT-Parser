package models

type Alert struct {
	ID            string  `bson:"_id,omitempty" json:"id"`
	SeverityLevel string  `bson:"severity_level" json:"severity_level"`
	Name          string  `bson:"name" json:"name"`
	ImageURL      string  `bson:"image_url" json:"image_url"`
	OriginalLink  string  `bson:"original_link" json:"original_link"`
	Details       Details `bson:"details" json:"details"`
}

type Details struct {
	BestPractice     string            `bson:"best_practice" json:"best_practice"`
	Description      string            `bson:"description" json:"description"`
	TargetedSector   string            `bson:"targeted_sector" json:"targeted_sector"`
	Threats          string            `bson:"threats" json:"threats"`
	WarningNumber    string            `bson:"warning_number" json:"warning_number"`
	WarningDate      string            `bson:"warning_date" json:"warning_date"`
	AffectedProducts []AffectedProduct `bson:"affected_products" json:"affected_products"`
	ThreatList       []Threat          `bson:"threat_list" json:"threat_list"`
	Recommendations  []Recommendation  `bson:"recommendations" json:"recommendation_links"`
}

type AffectedProduct struct {
	Name string `bson:"name" json:"name"`
}

type Threat struct {
	Name string `bson:"name" json:"name"`
}

type Recommendation struct {
	Link string `bson:"link" json:"link"`
}
