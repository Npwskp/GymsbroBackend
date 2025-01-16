package workout

type CreateWorkoutDto struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	Exercises   []WorkoutExercise `json:"exercises" validate:"dive"`
}

type UpdateWorkoutDto struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exercises   []WorkoutExercise `json:"exercises" validate:"dive"`
}
