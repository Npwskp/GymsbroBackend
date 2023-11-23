package function

import (
	"strings"
)

var Day = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

var ExerciseType = []string{"Rest", "Push", "Pull", "Chest", "Back", "Legs", "Shoulders", "Arms", "Abs"}

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
