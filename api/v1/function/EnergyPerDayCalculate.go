package function

type EnergyConsumptionPlan struct {
	BMR                       float64
	ActivityLevel             string
	AllActivityCaloriesPerDay []*CalPerActivity
	Macronutrients            []*Macronutrients
}

type CalPerActivity struct {
	ActivityName string
	Calories     float64
}

type Macronutrients struct {
	Goal           string
	CarbPreference string
	Calories       float64
	Protein        float64
	Fat            float64
	Carbs          float64
}

var ActivityLevel = []string{"base", "sedentary", "lightly_active", "moderately_active", "very_active", "extra_active"}
var Goal = []string{"maintain", "cutting", "bulking"}
var CarbPreference = []string{"moderate_carb", "low_carb", "high_carb"}

func CalculateBMR(weight float64, height float64, age int, gender string) float64 {
	if gender == "male" {
		return (10 * weight) + (6.25 * height) - (5 * float64(age)) + 5
	} else {
		return (10 * weight) + (6.25 * height) - (5 * float64(age)) - 161
	}
}

func CalculateCaloriesPerDay(bmr float64, activityLevel string) []*CalPerActivity {
	caloriesPerDay := []*CalPerActivity{}

	// Add base BMR
	caloriesPerDay = append(caloriesPerDay, &CalPerActivity{
		ActivityName: "base",
		Calories:     bmr,
	})

	// Calculate calories for all activity levels
	activityMultipliers := map[string]float64{
		"sedentary":         1.2,
		"lightly_active":    1.375,
		"moderately_active": 1.55,
		"very_active":       1.725,
		"extra_active":      1.9,
	}

	for activity, multiplier := range activityMultipliers {
		caloriesPerDay = append(caloriesPerDay, &CalPerActivity{
			ActivityName: activity,
			Calories:     bmr * multiplier,
		})
	}

	return caloriesPerDay
}

func CalculateMacronutrients(calories float64) []*Macronutrients {
	macros := []*Macronutrients{}

	for _, goal := range Goal {
		for _, carbPreference := range CarbPreference {
			var proteinRatio, fatRatio, carbRatio float64

			switch goal {
			case "maintain":
				switch carbPreference {
				case "moderate_carb":
					// 30/30/40 (protein/fat/carb)
					proteinRatio = 0.30
					fatRatio = 0.30
					carbRatio = 0.40
				case "low_carb":
					// 35/35/30
					proteinRatio = 0.35
					fatRatio = 0.35
					carbRatio = 0.30
				case "high_carb":
					// 25/25/50
					proteinRatio = 0.25
					fatRatio = 0.25
					carbRatio = 0.50
				}

			case "cutting":
				calories = calories - 500
				switch carbPreference {
				case "moderate_carb":
					// 40/30/30
					proteinRatio = 0.40
					fatRatio = 0.30
					carbRatio = 0.30
				case "low_carb":
					// 45/35/20
					proteinRatio = 0.45
					fatRatio = 0.35
					carbRatio = 0.20
				case "high_carb":
					// 35/25/40
					proteinRatio = 0.35
					fatRatio = 0.25
					carbRatio = 0.40
				}

			case "bulking":
				calories = calories + 500
				switch carbPreference {
				case "moderate_carb":
					// 25/25/50
					proteinRatio = 0.25
					fatRatio = 0.25
					carbRatio = 0.50
				case "low_carb":
					// 30/35/35
					proteinRatio = 0.30
					fatRatio = 0.35
					carbRatio = 0.35
				case "high_carb":
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
	allActivityCaloriesPerDay := CalculateCaloriesPerDay(bmr, ActivityLevel[activityLevel])
	macronutrients := CalculateMacronutrients(bmr)

	return &EnergyConsumptionPlan{BMR: bmr, ActivityLevel: ActivityLevel[activityLevel], AllActivityCaloriesPerDay: allActivityCaloriesPerDay, Macronutrients: macronutrients}, nil
}
