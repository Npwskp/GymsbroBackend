package workoutSession

type CreateWorkoutSessionDto struct {
	WorkoutID string `json:"workoutId" validate:"required"`
	Notes     string `json:"notes"`
}

type UpdateWorkoutSessionDto struct {
	Status    SessionStatus   `json:"status"`
	Exercises []ExerciseEntry `json:"exercises"`
	Notes     string          `json:"notes"`
}

type CompleteExerciseDto struct {
	ExerciseLogID string  `json:"exerciseLogId" validate:"required"`
	TotalVolume   float64 `json:"totalVolume"`
}
