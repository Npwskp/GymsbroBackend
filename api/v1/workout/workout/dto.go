package workout

type CreateWorkoutDto struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	Exercises   []WorkoutExercise `json:"exercises" validate:"required"`
}

type UpdateWorkoutDto struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	Exercises   []WorkoutExercise `json:"exercises" validate:"required"`
}

type SearchWorkoutFilters struct {
	Query string `json:"query" query:"query"`
}
