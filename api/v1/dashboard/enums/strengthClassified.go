package dashboardEnums

type strengthType string

const (
	strengthTypeBeginner     strengthType = "beginner"     // Score < 30
	strengthTypeNovice       strengthType = "novice"       // Score 30-45
	strengthTypeIntermediate strengthType = "intermediate" // Score 45-75
	strengthTypeAdvanced     strengthType = "advanced"     // Score 75-112.5
	strengthTypeElite        strengthType = "elite"        // Score > 112.5
)

// Strength score thresholds
const (
	BeginnerScore     float64 = 30.0
	NoviceScore       float64 = 45.0
	IntermediateScore float64 = 75.0
	AdvancedScore     float64 = 112.5
)

// ClassifyStrength returns the strengthType based on the given strength score
func ClassifyStrength(score float64) strengthType {
	switch {
	case score < BeginnerScore:
		return strengthTypeBeginner
	case score < NoviceScore:
		return strengthTypeNovice
	case score < IntermediateScore:
		return strengthTypeIntermediate
	case score < AdvancedScore:
		return strengthTypeAdvanced
	default:
		return strengthTypeElite
	}
}

// GetMinScoreForType returns the minimum score threshold for a given strength type
func GetMinScoreForType(st strengthType) float64 {
	switch st {
	case strengthTypeBeginner:
		return 0.0
	case strengthTypeNovice:
		return BeginnerScore
	case strengthTypeIntermediate:
		return NoviceScore
	case strengthTypeAdvanced:
		return IntermediateScore
	case strengthTypeElite:
		return AdvancedScore
	default:
		return 0.0
	}
}
