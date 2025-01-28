package userFitnessPreferenceEnums

import "math"

type BodyCompositionInfo struct {
	BMI                float64 `json:"bmi" default:"0"`
	BodyFatMass        float64 `json:"bodyfat_mass" default:"0"`
	BodyFatPercentage  float64 `json:"bodyfat_percentage" default:"0"`
	SkeletalMuscleMass float64 `json:"skeletal_muscle_mass" default:"0"`
	ExtracellularWater float64 `json:"extracellular_water" default:"0"`
	ECWRatio           float64 `json:"ecw_ratio" default:"0"`
}

func CalculateBMI(weight float64, height float64) float64 {
	// Convert height from cm to meters
	heightInMeters := height / 100
	if heightInMeters == 0 {
		return 0
	}
	return math.Round((weight/(heightInMeters*heightInMeters))*10) / 10
}
