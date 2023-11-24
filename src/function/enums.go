package function

import (
	"strings"
)

var Day = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

var ExerciseType = []string{"Rest", "Push", "Pull", "Chest", "Back", "Legs", "Shoulders", "Arms", "Abs"}

var MuscleGroup = []string{"Triceps", "Biceps", "Forearms", "Upper Chest", "Middle Chest", "Lower Chest", "Lat", "Trap", "Front Shoulders", "Side Shoulders", "Rear Shoulders", "Abs", "Side Abs", "Quads", "Hamstrings", "Calves", "Glutes"}

func CheckDay(day string) bool {
	checked := false
	for _, d := range Day {
		if strings.Compare(d, day) == 0 {
			checked = true
		}
	}
	return checked
}

func CheckExerciseType(exerciseType string) bool {
	checked := false
	for _, et := range ExerciseType {
		if strings.Compare(et, exerciseType) == 0 {
			checked = true
		}
	}
	return checked
}

func CheckMuscleGroup(muscleGroup string) bool {
	checked := false
	for _, mg := range MuscleGroup {
		if strings.Compare(mg, muscleGroup) == 0 {
			checked = true
		}
	}
	return checked
}
