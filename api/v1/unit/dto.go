package unit

// ConversionRequest represents a request to convert between units
type ConversionRequest struct {
	Value    float64 `json:"value"`
	FromUnit string  `json:"from_unit"`
	ToUnit   string  `json:"to_unit"`
}

// ConversionResponse represents the response from a unit conversion
type ConversionResponse struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}
