package function

import (
	"fmt"
)

// UnitToGramMap maps different weight units to their conversion factor to grams
var UnitToGramMap = map[string]float64{
	// Small mass units
	"µg":  1e-6, // microgram to gram
	"Âµg": 1e-6, // alternative encoding of microgram to gram
	"mg":  1e-3, // milligram to gram
	"g":   1,    // gram to gram

	// Larger mass units
	"kg": 1000,    // kilogram to gram
	"oz": 28.3495, // ounce to gram
	"lb": 453.592, // pound to gram
	"t":  1e6,     // metric ton to gram

	// Volume units (based on water density at room temperature)
	"ml":    1,       // milliliter to gram (assuming water density)
	"l":     1000,    // liter to gram (assuming water density)
	"fl_oz": 29.5735, // fluid ounce to gram (assuming water density)
	"cup":   236.588, // US cup to gram (assuming water density)
	"pt":    473.176, // US pint to gram (assuming water density)
	"qt":    946.353, // US quart to gram (assuming water density)
	"gal":   3785.41, // US gallon to gram (assuming water density)

	// Common cooking units (dry measurements)
	"tsp":     4.92892, // teaspoon to gram (assuming water density)
	"tbsp":    14.7868, // tablespoon to gram (assuming water density)
	"1/4_cup": 59.1470, // 1/4 cup to gram (assuming water density)
	"1/3_cup": 78.8627, // 1/3 cup to gram (assuming water density)
	"1/2_cup": 118.294, // 1/2 cup to gram (assuming water density)
}

// ConvertUnit converts a value from one unit to another
// Returns the converted value and any error that occurred
func ConvertUnit(value float64, fromUnit string, toUnit string) (float64, error) {
	// Check if both units exist in our conversion map
	fromFactor, fromExists := UnitToGramMap[fromUnit]
	if !fromExists {
		return 0, fmt.Errorf("unsupported source unit: %s", fromUnit)
	}

	toFactor, toExists := UnitToGramMap[toUnit]
	if !toExists {
		return 0, fmt.Errorf("unsupported target unit: %s", toUnit)
	}

	// Convert to grams first, then to target unit
	valueInGrams := value * fromFactor
	result := valueInGrams / toFactor

	return result, nil
}
