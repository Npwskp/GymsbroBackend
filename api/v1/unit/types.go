package unit

// UnitType represents the category of measurement
type UnitType string

// UnitInfo stores information about a unit
type UnitInfo struct {
	Symbol      string   `json:"symbol"`
	Type        UnitType `json:"type"`
	DisplayName string   `json:"displayName"`
	ToGrams     float64  `json:"toGrams"`
}
