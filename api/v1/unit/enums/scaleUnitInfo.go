package unitEnums

// ScaleUnitType represents the category of measurement
type ScaleUnitType string

// Unit type constants
const (
	MassUnit    ScaleUnitType = "mass"
	VolumeUnit  ScaleUnitType = "volume"
	CookingUnit ScaleUnitType = "cooking"
)

// ScaleUnitInfo stores information about a unit
type ScaleUnitInfo struct {
	Symbol      string        `json:"symbol"`
	Type        ScaleUnitType `json:"type"`
	DisplayName string        `json:"displayName"`
	ToGrams     float64       `json:"toGrams"`
}

// ScaleUnitInfoMap stores detailed information about each unit
var ScaleUnitInfoMap = map[string]ScaleUnitInfo{
	// Mass units
	"µg":  {Symbol: "µg", Type: MassUnit, DisplayName: "Microgram", ToGrams: 1e-6},
	"mg":  {Symbol: "mg", Type: MassUnit, DisplayName: "Milligram", ToGrams: 1e-3},
	"g":   {Symbol: "g", Type: MassUnit, DisplayName: "Gram", ToGrams: 1},
	"kg":  {Symbol: "kg", Type: MassUnit, DisplayName: "Kilogram", ToGrams: 1000},
	"oz":  {Symbol: "oz", Type: MassUnit, DisplayName: "Ounce", ToGrams: 28.3495},
	"lbs": {Symbol: "lbs", Type: MassUnit, DisplayName: "Pound", ToGrams: 453.592},
	"t":   {Symbol: "t", Type: MassUnit, DisplayName: "Metric Ton", ToGrams: 1e6},

	// Volume units
	"ml":    {Symbol: "ml", Type: VolumeUnit, DisplayName: "Milliliter", ToGrams: 1},
	"l":     {Symbol: "l", Type: VolumeUnit, DisplayName: "Liter", ToGrams: 1000},
	"fl_oz": {Symbol: "fl_oz", Type: VolumeUnit, DisplayName: "Fluid Ounce", ToGrams: 29.5735},
	"cup":   {Symbol: "cup", Type: VolumeUnit, DisplayName: "Cup", ToGrams: 236.588},
	"pt":    {Symbol: "pt", Type: VolumeUnit, DisplayName: "Pint", ToGrams: 473.176},
	"qt":    {Symbol: "qt", Type: VolumeUnit, DisplayName: "Quart", ToGrams: 946.353},
	"gal":   {Symbol: "gal", Type: VolumeUnit, DisplayName: "Gallon", ToGrams: 3785.41},

	// Cooking units
	"tsp":     {Symbol: "tsp", Type: CookingUnit, DisplayName: "Teaspoon", ToGrams: 4.92892},
	"tbsp":    {Symbol: "tbsp", Type: CookingUnit, DisplayName: "Tablespoon", ToGrams: 14.7868},
	"1/4_cup": {Symbol: "1/4_cup", Type: CookingUnit, DisplayName: "1/4 Cup", ToGrams: 59.1470},
	"1/3_cup": {Symbol: "1/3_cup", Type: CookingUnit, DisplayName: "1/3 Cup", ToGrams: 78.8627},
	"1/2_cup": {Symbol: "1/2_cup", Type: CookingUnit, DisplayName: "1/2 Cup", ToGrams: 118.294},
}

type ExerciseWeightUnit string

const (
	ExerciseWeightUnitPound ExerciseWeightUnit = "lbs"
	ExerciseWeightUnitKg    ExerciseWeightUnit = "kg"
)

func GetAllExerciseWeightUnit() []ExerciseWeightUnit {
	return []ExerciseWeightUnit{ExerciseWeightUnitPound, ExerciseWeightUnitKg}
}
