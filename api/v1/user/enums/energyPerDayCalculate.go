package userFitnessPreferenceEnums

import (
	"fmt"
)

// ActivityLevelType represents the activity level enum
type ActivityLevelType string

// GoalType represents the goal enum
type GoalType string

// CarbPreferenceType represents the carb preference enum
type CarbPreferenceType string

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

// Update the struct to use the new types
type EnergyConsumptionPlan struct {
	BMR                       float64
	ActivityLevel             ActivityLevelType
	AllActivityCaloriesPerDay []*CalPerActivity
	Macronutrients            []*Macronutrients
}

type CalPerActivity struct {
	ActivityName string
	Calories     float64
}

type Macronutrients struct {
	Goal           GoalType
	CarbPreference CarbPreferenceType
	Calories       float64
	Protein        float64
	Fat            float64
	Carbs          float64
}

func CalculateBMR(weight float64, height float64, age int, gender string) float64 {
	if gender == "male" {
		return (10 * weight) + (6.25 * height) - (5 * float64(age)) + 5
	} else {
		return (10 * weight) + (6.25 * height) - (5 * float64(age)) - 161
	}
}

func CalculateCaloriesPerDay(bmr float64, activityLevel ActivityLevelType) []*CalPerActivity {
	caloriesPerDay := []*CalPerActivity{}

	// Add base BMR
	caloriesPerDay = append(caloriesPerDay, &CalPerActivity{
		ActivityName: string(ActivitySedentary),
		Calories:     bmr,
	})

	// Calculate calories for all activity levels
	activityMultipliers := map[ActivityLevelType]float64{
		ActivitySedentary:     1.2,
		ActivityLightlyActive: 1.375,
		ActivityModerate:      1.55,
		ActivityVeryActive:    1.725,
		ActivityExtraActive:   1.9,
	}

	for activity, multiplier := range activityMultipliers {
		caloriesPerDay = append(caloriesPerDay, &CalPerActivity{
			ActivityName: string(activity),
			Calories:     bmr * multiplier,
		})
	}

	return caloriesPerDay
}

func CalculateMacronutrients(calories float64) []*Macronutrients {
	macros := []*Macronutrients{}

	for _, goal := range GetAllGoals() {
		for _, carbPreference := range GetAllCarbPreferences() {
			var proteinRatio, fatRatio, carbRatio float64

			switch goal {
			case GoalMaintain:
				switch carbPreference {
				case CarbModerate:
					// 30/30/40 (protein/fat/carb)
					proteinRatio = 0.30
					fatRatio = 0.30
					carbRatio = 0.40
				case CarbLow:
					// 35/35/30
					proteinRatio = 0.35
					fatRatio = 0.35
					carbRatio = 0.30
				case CarbHigh:
					// 25/25/50
					proteinRatio = 0.25
					fatRatio = 0.25
					carbRatio = 0.50
				}

			case GoalCutting:
				calories = calories - 500
				switch carbPreference {
				case CarbModerate:
					// 40/30/30
					proteinRatio = 0.40
					fatRatio = 0.30
					carbRatio = 0.30
				case CarbLow:
					// 45/35/20
					proteinRatio = 0.45
					fatRatio = 0.35
					carbRatio = 0.20
				case CarbHigh:
					// 35/25/40
					proteinRatio = 0.35
					fatRatio = 0.25
					carbRatio = 0.40
				}

			case GoalBulking:
				calories = calories + 500
				switch carbPreference {
				case CarbModerate:
					// 25/25/50
					proteinRatio = 0.25
					fatRatio = 0.25
					carbRatio = 0.50
				case CarbLow:
					// 30/35/35
					proteinRatio = 0.30
					fatRatio = 0.35
					carbRatio = 0.35
				case CarbHigh:
					// 20/20/60
					proteinRatio = 0.20
					fatRatio = 0.20
					carbRatio = 0.60
				}
			}

			proteinCals := calories * proteinRatio
			fatCals := calories * fatRatio
			carbsCals := calories * carbRatio

			macros = append(macros, &Macronutrients{
				Goal:           goal,
				CarbPreference: carbPreference,
				Calories:       calories,
				Protein:        proteinCals / 4,
				Fat:            fatCals / 9,
				Carbs:          carbsCals / 4,
			})
		}
	}

	return macros
}

func GetUserEnergyConsumePlan(weight float64, height float64, age int, gender string, activityLevel int, goal string) (*EnergyConsumptionPlan, error) {
	bmr := CalculateBMR(weight, height, age, gender)
	allActivityLevels := GetAllActivityLevels()
	if activityLevel < 0 || activityLevel >= len(allActivityLevels) {
		return nil, fmt.Errorf("invalid activity level index")
	}

	selectedActivity := allActivityLevels[activityLevel]
	allActivityCaloriesPerDay := CalculateCaloriesPerDay(bmr, selectedActivity)
	macronutrients := CalculateMacronutrients(bmr)

	return &EnergyConsumptionPlan{
		BMR:                       bmr,
		ActivityLevel:             selectedActivity,
		AllActivityCaloriesPerDay: allActivityCaloriesPerDay,
		Macronutrients:            macronutrients,
	}, nil
}
