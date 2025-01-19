package workoutSession

import "time"

type CreateWorkoutSessionDto struct {
	WorkoutID string      `json:"workoutId"`
	Type      SessionType `json:"type" validate:"required,oneof=planned custom"`
	Notes     string      `json:"notes"`
}

type UpdateWorkoutSessionDto struct {
	Status    SessionStatus     `json:"status"`
	Exercises []SessionExercise `json:"exercises"`
	Notes     string            `json:"notes"`
}

type ReorderExercisesDto struct {
	Exercises []ExerciseOrder `json:"exercises" validate:"required,dive"`
}

type ExerciseOrder struct {
	ExerciseID string `json:"exerciseId" validate:"required"`
	Order      int    `json:"order" validate:"required,min=0"`
}

type CompleteExerciseDto struct {
	ExerciseLogID string `json:"exerciseLogId" validate:"required"`
}

type LoggedSessionDto struct {
	WorkoutID string            `json:"workoutId"`
	StartTime time.Time         `json:"startTime" validate:"required"`
	EndTime   time.Time         `json:"endTime" validate:"required"`
	Status    SessionStatus     `json:"status" validate:"required"`
	Exercises []SessionExercise `json:"exercises"`
	Notes     string            `json:"notes"`
}
