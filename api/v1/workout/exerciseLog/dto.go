package exerciseLog

type CreateExerciseLogDto struct {
	ExerciseID       string   `json:"exerciseId" validate:"required"`
	WorkoutSessionID string   `json:"workoutSessionId" validate:"required"`
	Sets             []SetLog `json:"sets" validate:"required,dive"`
	Notes            string   `json:"notes"`
}

type UpdateExerciseLogDto struct {
	Sets  []SetLog `json:"sets" validate:"required,dive"`
	Notes string   `json:"notes"`
}
