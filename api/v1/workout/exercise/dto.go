package exercise

type CreateExerciseDto struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Type        []string `json:"type" validate:"required"`
	Muscle      []string `json:"muscle" validate:"required"`
	Image       string   `json:"image" validate:"required"`
}

type UpdateExerciseDto struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        []string `json:"type"`
	Muscle      []string `json:"muscle"`
	Image       string   `json:"image"`
}
