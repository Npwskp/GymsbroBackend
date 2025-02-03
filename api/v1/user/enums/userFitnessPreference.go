package userFitnessPreferenceEnums

import (
	"math"

	authEnums "github.com/Npwskp/GymsbroBackend/api/v1/auth/enums"
)

// ActivityLevelType represents the activity level enum
type ActivityLevelType string

// GoalType represents the goal enum
type GoalType string

// CarbPreferenceType represents the carb preference enum
type CarbPreferenceType string

type NutritionInfo struct {
	BMR           float64           `json:"bmr"`
	ActivityLevel ActivityLevelType `json:"activity_level"`
	Goal          GoalType          `json:"goal"`
}

type EnergyConsumptionPlan struct {
	BMR            float64           `json:"bmr"`
	ActivityLevel  ActivityLevelType `json:"activity_level"`
	Goal           GoalType          `json:"goal"`
	Macronutrients []*Macronutrients `json:"macronutrients"`
}

type Macronutrients struct {
	CarbPreference CarbPreferenceType `json:"carb_preference"`
	Calories       float64            `json:"calories"`
	Protein        float64            `json:"protein"`
	Fat            float64            `json:"fat"`
	Carbs          float64            `json:"carbs"`
}

const (
	// Activity Levels
	ActivitySedentary     ActivityLevelType = "sedentary"
	ActivityLightlyActive ActivityLevelType = "lightly_active"
	ActivityModerate      ActivityLevelType = "moderately_active"
	ActivityVeryActive    ActivityLevelType = "very_active"
	ActivityExtraActive   ActivityLevelType = "extra_active"

	// Goals
	GoalMaintain GoalType = "maintain"
	GoalCutting  GoalType = "cutting"
	GoalBulking  GoalType = "bulking"

	// Carb Preferences
	CarbModerate CarbPreferenceType = "moderate_carb"
	CarbLow      CarbPreferenceType = "low_carb"
	CarbHigh     CarbPreferenceType = "high_carb"
	CarbManual   CarbPreferenceType = "manual"
)

// Replace the slice variables with functions that return all possible values
func GetAllActivityLevels() []ActivityLevelType {
	return []ActivityLevelType{
		ActivitySedentary,
		ActivityLightlyActive,
		ActivityModerate,
		ActivityVeryActive,
		ActivityExtraActive,
	}
}

func GetAllGoals() []GoalType {
	return []GoalType{
		GoalMaintain,
		GoalCutting,
		GoalBulking,
	}
}

func GetAllCarbPreferences() []CarbPreferenceType {
	return []CarbPreferenceType{
		CarbModerate,
		CarbLow,
		CarbHigh,
	}
}

func CalculateBMR(weight float64, height float64, age int, gender authEnums.GenderType) float64 {
	if gender == authEnums.GenderMale {
		return (10 * weight) + (6.25 * height) - (5 * float64(age)) + 5
	} else if gender == authEnums.GenderFemale {
		return (10 * weight) + (6.25 * height) - (5 * float64(age)) - 161
	}

	return 0
}

func CalculateCaloriesPerDay(bmr float64, activityLevel ActivityLevelType, goal GoalType) float64 {
	var calories float64

	switch activityLevel {
	case ActivitySedentary:
		calories = bmr * 1.2
	case ActivityLightlyActive:
		calories = bmr * 1.375
	case ActivityModerate:
		calories = bmr * 1.55
	case ActivityVeryActive:
		calories = bmr * 1.725
	case ActivityExtraActive:
		calories = bmr * 1.9
	default:
		calories = bmr
	}

	if goal == GoalMaintain {
		return math.Round(calories)
	} else if goal == GoalCutting {
		return math.Round(calories) - 500
	} else if goal == GoalBulking {
		return math.Round(calories) + 500
	}

	return math.Round(calories)
}

func CalculateMacronutrients(calories float64) []*Macronutrients {
	macros := []*Macronutrients{}

	var proteinRatio, fatRatio, carbRatio float64

	for _, carbPreference := range GetAllCarbPreferences() {
		switch carbPreference {
		case CarbModerate:
			// 30/35/35 (protein/fat/carb)
			proteinRatio = 0.30
			fatRatio = 0.35
			carbRatio = 0.35
		case CarbLow:
			// 40/40/20
			proteinRatio = 0.40
			fatRatio = 0.40
			carbRatio = 0.20
		case CarbHigh:
			// 30/20/50
			proteinRatio = 0.30
			fatRatio = 0.20
			carbRatio = 0.50
		}

		proteinCals := calories * proteinRatio
		fatCals := calories * fatRatio
		carbsCals := calories * carbRatio

		macros = append(macros, &Macronutrients{
			CarbPreference: carbPreference,
			Calories:       math.Round(calories),
			Protein:        math.Round(proteinCals / 4),
			Fat:            math.Round(fatCals / 9),
			Carbs:          math.Round(carbsCals / 4),
		})
	}

	return macros
}

func GetUserEnergyConsumePlan(weight float64, height float64, age int, gender authEnums.GenderType, activityLevel ActivityLevelType, goal GoalType) (*EnergyConsumptionPlan, error) {
	bmr := CalculateBMR(weight, height, age, gender)

	caloriesPerDay := CalculateCaloriesPerDay(bmr, activityLevel, goal)
	macronutrients := CalculateMacronutrients(caloriesPerDay)

	return &EnergyConsumptionPlan{
		BMR:            bmr,
		ActivityLevel:  activityLevel,
		Goal:           goal,
		Macronutrients: macronutrients,
	}, nil
}
