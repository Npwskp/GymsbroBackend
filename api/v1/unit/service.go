package unit

import (
	"fmt"

	unitEnums "github.com/Npwskp/GymsbroBackend/api/v1/unit/enums"
)

type UnitService struct{}

type IUnitService interface {
	GetUnit(symbol string) (unitEnums.UnitInfo, bool)
	IsValidUnit(symbol string) bool
	GetAllUnits() []unitEnums.UnitInfo
	ConvertToGrams(value float64, fromUnit string) (float64, error)
	ConvertBetweenUnits(value float64, fromUnit, toUnit string) (float64, error)
}

// GetUnit returns unit info for a given symbol
func (s *UnitService) GetUnit(symbol string) (unitEnums.UnitInfo, bool) {
	unit, exists := unitEnums.UnitInfoMap[symbol]
	return unit, exists
}

// IsValidUnit checks if a unit symbol is valid
func (s *UnitService) IsValidUnit(symbol string) bool {
	_, exists := unitEnums.UnitInfoMap[symbol]
	return exists
}

func (s *UnitService) GetAllUnits() []unitEnums.UnitInfo {
	units := make([]unitEnums.UnitInfo, 0, len(unitEnums.UnitInfoMap))
	for _, unit := range unitEnums.UnitInfoMap {
		units = append(units, unit)
	}
	return units
}

// ConvertToGrams converts a value from one unit to grams
func (s *UnitService) ConvertToGrams(value float64, fromUnit string) (float64, error) {
	unit, exists := unitEnums.UnitInfoMap[fromUnit]
	if !exists {
		return 0, fmt.Errorf("invalid unit: %s", fromUnit)
	}
	return value * unit.ToGrams, nil
}

// ConvertBetweenUnits converts a value from one unit to another
func (s *UnitService) ConvertBetweenUnits(value float64, fromUnit, toUnit string) (float64, error) {
	// First validate both units
	fromUnitInfo, fromExists := unitEnums.UnitInfoMap[fromUnit]
	toUnitInfo, toExists := unitEnums.UnitInfoMap[toUnit]

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
