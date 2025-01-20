package exerciseLog

import "time"

type CreateExerciseLogDto struct {
	ExerciseID string    `json:"exerciseId" validate:"required"`
	DateTime   time.Time `json:"dateTime"`
	Sets       []SetLog  `json:"sets" validate:"required,dive"`
	Notes      string    `json:"notes"`
}

type UpdateExerciseLogDto struct {
	Sets     []SetLog  `json:"sets" validate:"required,dive"`
	DateTime time.Time `json:"dateTime"`
	Notes    string    `json:"notes"`
}
