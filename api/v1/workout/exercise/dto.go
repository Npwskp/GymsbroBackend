package exercise

import (
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
)

type CreateExerciseDto struct {
	Name         string                       `json:"name" validate:"required"`
	Equipment    exerciseEnums.Equipment      `json:"equipment" validate:"required"`
	Mechanics    exerciseEnums.Mechanics      `json:"mechanics" validate:"required"`
	Force        exerciseEnums.Force          `json:"force" validate:"required"`
	Preparation  []string                     `json:"preparation" validate:"required"`
	Execution    []string                     `json:"execution" validate:"required"`
	Image        string                       `json:"image"`
	BodyPart     []exerciseEnums.BodyPart     `json:"body_part" validate:"required"`
	TargetMuscle []exerciseEnums.TargetMuscle `json:"target_muscle" validate:"required"`
}

type UpdateExerciseDto struct {
	Name         string                       `json:"name"`
	Equipment    exerciseEnums.Equipment      `json:"equipment"`
	Mechanics    exerciseEnums.Mechanics      `json:"mechanics"`
	Force        exerciseEnums.Force          `json:"force"`
	Preparation  []string                     `json:"preparation"`
	Execution    []string                     `json:"execution"`
	Image        string                       `json:"image"`
	BodyPart     []exerciseEnums.BodyPart     `json:"body_part"`
	TargetMuscle []exerciseEnums.TargetMuscle `json:"target_muscle"`
}

type SearchExerciseFilters struct {
	Equipment    string `query:"equipment"`
	Mechanics    string `query:"mechanics"`
	Force        string `query:"force"`
	BodyPart     string `query:"body_part"`
	TargetMuscle string `query:"target_muscle"` // Comma-separated target muscles
	Query        string `query:"query"`         // Search query for exercise name
}
