package types

type Nutrient struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Category string  `json:"category"`
	Value    float64 `json:"value"`
	Unit     string  `json:"unit"`
}
