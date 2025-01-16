package workoutPlan

type CreatePlanByDaysOfWeekDto struct {
	MondayWorkoutID    string `json:"mondayWorkoutId" validate:"required"`
	TuesdayWorkoutID   string `json:"tuesdayWorkoutId" validate:"required"`
	WednesdayWorkoutID string `json:"wednesdayWorkoutId" validate:"required"`
	ThursdayWorkoutID  string `json:"thursdayWorkoutId" validate:"required"`
	FridayWorkoutID    string `json:"fridayWorkoutId" validate:"required"`
	SaturdayWorkoutID  string `json:"saturdayWorkoutId" validate:"required"`
	SundayWorkoutID    string `json:"sundayWorkoutId" validate:"required"`
	WeeksDuration      int    `json:"weeksDuration" validate:"required"`
}

type CreatePlanByCyclicWorkoutDto struct {
	WorkoutIDs    []string `json:"workoutIds" validate:"required"`
	WeeksDuration int      `json:"weeksDuration" validate:"required"`
}
