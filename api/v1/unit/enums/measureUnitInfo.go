package unitEnums

type MeasureUnitType string

const (
	Imperial MeasureUnitType = "IMPERIAL"
	Metric   MeasureUnitType = "METRIC"
)

type MeasureUnitInfo struct {
	Symbol      string          `json:"symbol"`
	Type        MeasureUnitType `json:"type"`
	DisplayName string          `json:"displayName"`
	ToMeter     float64         `json:"toMeter"`
}

var MeasureUnitInfoMap = map[string]MeasureUnitInfo{
	"in": {Symbol: "in", Type: Imperial, DisplayName: "Inch", ToMeter: 0.0254},
	"ft": {Symbol: "ft", Type: Imperial, DisplayName: "Foot", ToMeter: 0.3048},
	"yd": {Symbol: "yd", Type: Imperial, DisplayName: "Yard", ToMeter: 0.9144},
	"mi": {Symbol: "mi", Type: Imperial, DisplayName: "Mile", ToMeter: 1609.34},

	"km": {Symbol: "km", Type: Metric, DisplayName: "Kilometer", ToMeter: 1000},
	"cm": {Symbol: "cm", Type: Metric, DisplayName: "Centimeter", ToMeter: 0.01},
	"mm": {Symbol: "mm", Type: Metric, DisplayName: "Millimeter", ToMeter: 0.001},
	"µm": {Symbol: "µm", Type: Metric, DisplayName: "Micrometer", ToMeter: 1e-6},
}

type BodyPartMeasureUnit string

const (
	BodyPartMeasureUnitInch BodyPartMeasureUnit = "in"
	BodyPartMeasureUnitCm   BodyPartMeasureUnit = "cm"
)

func GetAllBodyPartMeasureUnit() []BodyPartMeasureUnit {
	return []BodyPartMeasureUnit{BodyPartMeasureUnitInch, BodyPartMeasureUnitCm}
}
