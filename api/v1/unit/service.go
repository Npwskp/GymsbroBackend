package unit

import (
	"fmt"

	unitEnums "github.com/Npwskp/GymsbroBackend/api/v1/unit/enums"
)

type UnitService struct{}

type unitType string

const (
	scaleUnitType   unitType = "scale"
	measureUnitType unitType = "measure"
)

type IUnitService interface {
	GetUnit(symbol string, unitType unitType) (interface{}, bool)
	GetAllUnits(unitType unitType) interface{}
	ConvertUnits(value float64, fromUnit, toUnit string, unitType unitType) (float64, error)
}

// GetUnit returns unit info for a given symbol
func (s *UnitService) GetUnit(symbol string, uType unitType) (interface{}, bool) {
	switch uType {
	case scaleUnitType:
		return unitEnums.ScaleUnitInfoMap[symbol], symbol != ""
	case measureUnitType:
		return unitEnums.MeasureUnitInfoMap[symbol], symbol != ""
	default:
		return nil, false
	}
}

// GetAllUnits returns all units of specified type
func (s *UnitService) GetAllUnits(uType unitType) interface{} {
	switch uType {
	case scaleUnitType:
		units := make([]unitEnums.ScaleUnitInfo, 0, len(unitEnums.ScaleUnitInfoMap))
		for _, unit := range unitEnums.ScaleUnitInfoMap {
			units = append(units, unit)
		}
		return units
	case measureUnitType:
		units := make([]unitEnums.MeasureUnitInfo, 0, len(unitEnums.MeasureUnitInfoMap))
		for _, unit := range unitEnums.MeasureUnitInfoMap {
			units = append(units, unit)
		}
		return units
	default:
		return nil
	}
}

// ConvertUnits converts a value between units of the same type
func (s *UnitService) ConvertUnits(value float64, fromUnit, toUnit string, uType unitType) (float64, error) {
	switch uType {
	case scaleUnitType:
		fromScaleUnitInfo, fromExists := unitEnums.ScaleUnitInfoMap[fromUnit]
		toScaleUnitInfo, toExists := unitEnums.ScaleUnitInfoMap[toUnit]

		if !fromExists {
			return 0, fmt.Errorf("invalid source unit: %s", fromUnit)
		}
		if !toExists {
			return 0, fmt.Errorf("invalid target unit: %s", toUnit)
		}

		valueInGrams := value * fromScaleUnitInfo.ToGrams
		return valueInGrams / toScaleUnitInfo.ToGrams, nil

	case measureUnitType:
		fromMeasureUnitInfo, fromExists := unitEnums.MeasureUnitInfoMap[fromUnit]
		toMeasureUnitInfo, toExists := unitEnums.MeasureUnitInfoMap[toUnit]

		if !fromExists {
			return 0, fmt.Errorf("invalid source unit: %s", fromUnit)
		}
		if !toExists {
			return 0, fmt.Errorf("invalid target unit: %s", toUnit)
		}

		valueInMeters := value * fromMeasureUnitInfo.ToMeter
		return valueInMeters / toMeasureUnitInfo.ToMeter, nil

	default:
		return 0, fmt.Errorf("invalid unit type")
	}
}
