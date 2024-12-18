package function

import (
	"fmt"
)

type Day string
type ExerciseType string
type MuscleGroup string

const (
	Monday    Day = "Monday"
	Tuesday   Day = "Tuesday"
	Wednesday Day = "Wednesday"
	Thursday  Day = "Thursday"
	Friday    Day = "Friday"
	Saturday  Day = "Saturday"
	Sunday    Day = "Sunday"
)

const (
	Rest      ExerciseType = "Rest"
	Push      ExerciseType = "Push"
	Pull      ExerciseType = "Pull"
	Chest     ExerciseType = "Chest"
	Back      ExerciseType = "Back"
	Legs      ExerciseType = "Legs"
	Shoulders ExerciseType = "Shoulders"
	Arms      ExerciseType = "Arms"
	Abs       ExerciseType = "Abs"
)

const (
	Triceps        MuscleGroup = "Triceps"
	Biceps         MuscleGroup = "Biceps"
	Forearms       MuscleGroup = "Forearms"
	UpperChest     MuscleGroup = "Upper Chest"
	MiddleChest    MuscleGroup = "Middle Chest"
	LowerChest     MuscleGroup = "Lower Chest"
	Latissimus     MuscleGroup = "Latissimus"
	Trapezius      MuscleGroup = "Trapezius"
	LowerBack      MuscleGroup = "Lower Back"
	FrontShoulders MuscleGroup = "Front Shoulders"
	SideShoulders  MuscleGroup = "Side Shoulders"
	RearShoulders  MuscleGroup = "Rear Shoulders"
	Abdominal      MuscleGroup = "Abdominal"
	SideAbs        MuscleGroup = "Side Abs"
	Quads          MuscleGroup = "Quads"
	Hamstrings     MuscleGroup = "Hamstrings"
	Calves         MuscleGroup = "Calves"
	Glutes         MuscleGroup = "Glutes"
)

func CheckDay(day string) bool {
	switch Day(day) {
	case Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday:
		return true
	default:
		return false
	}
}

func ParseExerciseType(s string) (ExerciseType, error) {
	switch ExerciseType(s) {
	case Rest, Push, Pull, Chest, Back, Legs, Shoulders, Arms, Abs:
		return ExerciseType(s), nil
	default:
		return "", fmt.Errorf("invalid exercise type: %s", s)
	}
}

func ParseMuscleGroup(s string) (MuscleGroup, error) {
	switch MuscleGroup(s) {
	case Triceps, Biceps, Forearms, UpperChest, MiddleChest, LowerChest,
		Latissimus, Trapezius, LowerBack, FrontShoulders, SideShoulders,
		RearShoulders, Abdominal, SideAbs, Quads, Hamstrings, Calves, Glutes:
		return MuscleGroup(s), nil
	default:
		return "", fmt.Errorf("invalid muscle group: %s", s)
	}
}
