package dashboard

import (
	"time"

	dashboardEnums "github.com/Npwskp/GymsbroBackend/api/v1/dashboard/enums"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise"
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
)

type DashboardResponse struct {
	FrequencyGraph FrequencyGraphData  `json:"frequency_graph"`
	Analysis       WorkoutAnalysis     `json:"analysis"`
	TopProgress    []ExerciseProgress  `json:"top_progress"`
	TopFrequency   []ExerciseFrequency `json:"top_frequency"`
}

// FrequencyGraphData provides data points for plotting exercise frequency
type FrequencyGraphData struct {
	Labels    []string  `json:"labels"`    // Date labels in "2024-01-01" format
	Values    []int     `json:"values"`    // Number of exercises/workouts per day
	TrendLine []float64 `json:"trendline"` // 7-day moving average
}

// WorkoutAnalysis provides overall workout statistics
type WorkoutAnalysis struct {
	// General Stats
	TotalWorkouts          int     `json:"total_workouts"`
	TotalExercises         int     `json:"total_exercises"`
	TotalVolume            float64 `json:"total_volume"`
	AverageWorkoutDuration float64 `json:"average_workout_duration"`
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

type ExerciseProgress struct {
	ExerciseID     string            `json:"exerciseId"`
	Exercise       exercise.Exercise `json:"exercise"`
	StartVolume    float64           `json:"startVolume"`
	EndVolume      float64           `json:"endVolume"`
	VolumeProgress float64           `json:"volumeProgress"` // Percentage increase in volume
	StartOneRM     float64           `json:"startOneRM"`
	EndOneRM       float64           `json:"endOneRM"`
	OneRMProgress  float64           `json:"oneRMProgress"` // Percentage increase in 1RM
	Progress       float64           `json:"progress"`      // Average of volume and 1RM progress
	StartDate      time.Time         `json:"startDate"`
	EndDate        time.Time         `json:"endDate"`
}

type ExerciseFrequency struct {
	ExerciseID string            `json:"exerciseId"`
	Exercise   exercise.Exercise `json:"exercise"`
	Frequency  float64           `json:"frequency"`
}

type DailyNutritionSummary struct {
	Date          string  `json:"date"`
	TotalCalories float64 `json:"total_calories"`
	TotalProtein  float64 `json:"total_protein"`
	TotalCarbs    float64 `json:"total_carbs"`
	TotalFat      float64 `json:"total_fat"`
}

type NutritionSummaryResponse struct {
	DailySummaries  []DailyNutritionSummary `json:"daily_summaries"`
	AverageCalories float64                 `json:"average_calories"`
	AverageProtein  float64                 `json:"average_protein"`
	AverageCarbs    float64                 `json:"average_carbs"`
	AverageFat      float64                 `json:"average_fat"`
}
