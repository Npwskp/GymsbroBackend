package dashboardEnums

import exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"

type StrengthType string

type StrengthStandards map[string][]StrengthStandard

type StrengthStandard struct {
	Bodyweight float64 `json:"bodyweight"`
	Standards  struct {
		Beginner     float64 `json:"beginner"`
		Novice       float64 `json:"novice"`
		Intermediate float64 `json:"intermediate"`
		Advanced     float64 `json:"advanced"`
		Elite        float64 `json:"elite"`
	} `json:"standards"`
}

const (
	StrengthTypeUntrained    StrengthType = "untrained"    // Score = 0
	StrengthTypeBeginner     StrengthType = "beginner"     // Score < 30
	StrengthTypeNovice       StrengthType = "novice"       // Score 30-45
	StrengthTypeIntermediate StrengthType = "intermediate" // Score 45-75
	StrengthTypeAdvanced     StrengthType = "advanced"     // Score 75-112.5
	StrengthTypeElite        StrengthType = "elite"        // Score > 112.5
)

// Strength score thresholds
const (
	UntrainedScore    float64 = 0.0
	BeginnerScore     float64 = 30.0
	NoviceScore       float64 = 45.0
	IntermediateScore float64 = 75.0
	AdvancedScore     float64 = 112.5
	EliteScore        float64 = 125.0
)

type ConsiderExercise struct {
	Exercise  string
	Equipment exerciseEnums.Equipment
	ID        string
}

var ConsiderExercises = []ConsiderExercise{
	{Exercise: "Bench Press", Equipment: exerciseEnums.Barbell},
	{Exercise: "Bench Press", Equipment: exerciseEnums.Dumbbell},
	{Exercise: "Squat", Equipment: exerciseEnums.Barbell},
	{Exercise: "Deadlift", Equipment: exerciseEnums.Barbell},
	{Exercise: "Incline Bench Press", Equipment: exerciseEnums.Barbell},
	{Exercise: "Pull-up", Equipment: exerciseEnums.BodyWeight},
	{Exercise: "Shoulder Press", Equipment: exerciseEnums.Barbell},
	{Exercise: "Pulldown", Equipment: exerciseEnums.Cable},
	{Exercise: "Chest Dip", Equipment: exerciseEnums.BodyWeight}, // Triceps Dips or Chest Dips
	{Exercise: "Chin-up", Equipment: exerciseEnums.BodyWeight},
}

// ClassifyStrength returns the StrengthType based on the given strength score
func ClassifyStrength(score float64) StrengthType {
	switch {
	case score < BeginnerScore:
		return StrengthTypeBeginner
	case score < NoviceScore:
		return StrengthTypeNovice
	case score < IntermediateScore:
		return StrengthTypeIntermediate
	case score < AdvancedScore:
		return StrengthTypeAdvanced
	case score < EliteScore:
		return StrengthTypeElite
	default:
		return StrengthTypeUntrained
	}
}

// GetMinScoreForType returns the minimum score threshold for a given strength type
func GetMinScoreForType(st StrengthType) float64 {
	switch st {
	case StrengthTypeBeginner:
		return BeginnerScore
	case StrengthTypeNovice:
		return NoviceScore
	case StrengthTypeIntermediate:
		return IntermediateScore
	case StrengthTypeAdvanced:
		return AdvancedScore
	case StrengthTypeElite:
		return EliteScore
	default:
		return UntrainedScore
	}
}
