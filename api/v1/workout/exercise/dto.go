package exercise

import (
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
)

type CreateExerciseDto struct {
	Name        string                       `json:"name" validate:"required"`
	Description string                       `json:"description" validate:"required"`
	Type        []exerciseEnums.ExerciseType `json:"type" validate:"required"`
	Muscle      []exerciseEnums.MuscleGroup  `json:"muscle" validate:"required"`
	Image       string                       `json:"image" validate:"required"`
}

type UpdateExerciseDto struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Type        []exerciseEnums.ExerciseType `json:"type"`
	Muscle      []exerciseEnums.MuscleGroup  `json:"muscle"`
	Image       string                       `json:"image"`
}

type SearchExerciseFilters struct {
	Types   string `query:"types"`   // Comma-separated exercise types
	Muscles string `query:"muscles"` // Comma-separated muscle groups
	UserID  string
}
