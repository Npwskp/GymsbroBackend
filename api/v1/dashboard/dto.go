package dashboard

import (
	"time"

	dashboardEnums "github.com/Npwskp/GymsbroBackend/api/v1/dashboard/enums"
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
)

type DashboardResponse struct {
	FrequencyGraph FrequencyGraphData `json:"frequency_graph"`
	Analysis       WorkoutAnalysis    `json:"analysis"`
}

// FrequencyGraphData provides data points for plotting exercise frequency
type FrequencyGraphData struct {
	// Data points for last 30 days
	Labels    []string  `json:"labels"`    // Date labels in "2024-01-01" format
	Values    []int     `json:"values"`    // Number of exercises/workouts per day
	TrendLine []float64 `json:"trendline"` // 7-day moving average
}

// WorkoutAnalysis provides overall workout statistics
type WorkoutAnalysis struct {
	// General Stats
	TotalWorkouts  int     `json:"total_workouts"`
	TotalExercises int     `json:"total_exercises"`
	TotalVolume    float64 `json:"total_volume"`

	// Time-based Stats
	AveragePerWeek float64 `json:"average_per_week"`
	BestStreak     int     `json:"best_streak"`    // Best consecutive days streak
	CurrentStreak  int     `json:"current_streak"` // Current consecutive days streak

	// Pattern Analysis
	MostActiveDay  string `json:"most_active_day"`  // e.g., "Monday"
	MostActiveTime string `json:"most_active_time"` // "Morning", "Afternoon", "Evening", "Night"

	// Recent Trends
	LastWeekCount  int `json:"last_week_count"`  // Workouts in last 7 days
	LastMonthCount int `json:"last_month_count"` // Workouts in last 30 days
}

type UserStrengthStandards struct {
	ExerciseStandards    []UserStrengthStandardPerExercise    `json:"exerciseStandards"`
	MuscleGroupStrengths []UserStrengthStandardPerMuscleGroup `json:"muscleGroupStrengths"`
}

type UserStrengthStandardPerExercise struct {
	Exercise         string                      `json:"exercise"`
	Equipment        exerciseEnums.Equipment     `json:"equipment"`
	RepMax           float64                     `json:"repmax"`
	RelativeStrength float64                     `json:"relativeStrength"`
	StrengthLevel    dashboardEnums.StrengthType `json:"strengthLevel"`
	Score            float64                     `json:"score"`
	LastPerformed    time.Time                   `json:"lastPerformed"`
}

type UserStrengthStandardPerMuscleGroup struct {
	TargetMuscle  exerciseEnums.TargetMuscle  `json:"target_muscle"`
	StrengthLevel dashboardEnums.StrengthType `json:"strength_level"`
	Score         float64                     `json:"score"`
}

type RepMaxResponse struct {
	OneRepMax    float64   `json:"oneRepMax"`
	EightRepMax  float64   `json:"eightRepMax"`
	TwelveRepMax float64   `json:"twelveRepMax"`
	LastUpdated  time.Time `json:"lastUpdated"`
}
