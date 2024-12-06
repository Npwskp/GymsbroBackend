package unit

import "fmt"

// Global instance of UnitService
var Service *UnitService

// init function will be called automatically when the package is imported
func init() {
	Service = NewUnitService()
}

// UnitService provides helper methods for unit operations
type UnitService struct {
	unitMap map[string]UnitInfo
}

// NewUnitService creates a new instance of UnitService
func NewUnitService() *UnitService {
	return &UnitService{
		unitMap: UnitInfoMap,
	}
}

// GetUnit returns unit info for a given symbol
func (s *UnitService) GetUnit(symbol string) (UnitInfo, bool) {
	unit, exists := s.unitMap[symbol]
	return unit, exists
}

// IsValidUnit checks if a unit symbol is valid
func (s *UnitService) IsValidUnit(symbol string) bool {
	_, exists := s.unitMap[symbol]
	return exists
}

// ConvertToGrams converts a value from one unit to grams
func (s *UnitService) ConvertToGrams(value float64, fromUnit string) (float64, error) {
	unit, exists := s.unitMap[fromUnit]
	if !exists {
		return 0, fmt.Errorf("invalid unit: %s", fromUnit)
	}
	return value * unit.ToGrams, nil
}

// ConvertBetweenUnits converts a value from one unit to another
func (s *UnitService) ConvertBetweenUnits(value float64, fromUnit, toUnit string) (float64, error) {
	// First validate both units
	fromUnitInfo, fromExists := s.unitMap[fromUnit]
	toUnitInfo, toExists := s.unitMap[toUnit]

	if !fromExists {
		return 0, fmt.Errorf("invalid source unit: %s", fromUnit)
	}
	if !toExists {
		return 0, fmt.Errorf("invalid target unit: %s", toUnit)
	}

	// Check if units are of the same type
	if fromUnitInfo.Type != toUnitInfo.Type {
		return 0, fmt.Errorf("cannot convert between different unit types: %s (%s) to %s (%s)",
			fromUnit, fromUnitInfo.Type, toUnit, toUnitInfo.Type)
	}

	// Convert to grams first, then to target unit
	valueInGrams := value * fromUnitInfo.ToGrams
	result := valueInGrams / toUnitInfo.ToGrams

	return result, nil
}
