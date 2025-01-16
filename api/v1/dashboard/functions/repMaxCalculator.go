package dashboardFunctions

import (
	"errors"
	"math"
)

// CalculateOneRepMax calculates the theoretical one-rep maximum using the Brzycki formula
// weight: the weight used in the exercise (in any consistent unit like kg or lbs)
// reps: number of repetitions performed
// Returns the calculated 1RM and any error that occurred during calculation
func CalculateOneRepMax(weight, reps float64) (float64, error) {
	// Input validation
	if weight <= 0 {
		return 0, errors.New("weight must be greater than 0")
	}
	if reps < 1 {
		return 0, errors.New("reps must be at least 1")
	}
	if reps > 36 {
		return 0, errors.New("formula is not accurate for more than 36 reps")
	}

	// Brzycki formula: weight / (1.0278 – 0.0278 × reps)
	oneRepMax := weight / (1.0278 - 0.0278*reps)

	// Round to 2 decimal places for practical use
	oneRepMax = math.Round(oneRepMax*100) / 100

	return oneRepMax, nil
}

// CalculateAssistedOneRepMax calculates 1RM for assisted exercises (e.g., assisted pull-ups)
// bodyWeight: the person's body weight
// assistWeight: the weight of assistance provided by the machine
// reps: number of repetitions performed
// Returns the calculated 1RM (as body weight - required assistance) and any error
func CalculateAssistedOneRepMax(bodyWeight, assistWeight, reps float64) (float64, error) {
	// Input validation
	if bodyWeight <= 0 {
		return 0, errors.New("body weight must be greater than 0")
	}
	if assistWeight < 0 {
		return 0, errors.New("assist weight cannot be negative")
	}
	if assistWeight >= bodyWeight {
		return 0, errors.New("assist weight cannot be greater than or equal to body weight")
	}
	if reps < 1 {
		return 0, errors.New("reps must be at least 1")
	}
	if reps > 36 {
		return 0, errors.New("formula is not accurate for more than 36 reps")
	}

	// Calculate effective weight being lifted (body weight - assistance)
	effectiveWeight := bodyWeight - assistWeight

	// Apply Brzycki formula to the effective weight
	oneRepMax := effectiveWeight / (1.0278 - 0.0278*reps)

	// Round to 2 decimal places
	oneRepMax = math.Round(oneRepMax*100) / 100

	return oneRepMax, nil
}

// EstimateRepMax calculates the estimated max weight for a target number of reps
// oneRepMax: the calculated or known one rep maximum
// targetReps: the number of reps to estimate weight for
// Returns the estimated weight for the target reps and any error
func EstimateRepMax(oneRepMax float64, targetReps float64) (float64, error) {
	// Input validation
	if oneRepMax <= 0 {
		return 0, errors.New("one rep max must be greater than 0")
	}
	if targetReps < 1 {
		return 0, errors.New("target reps must be at least 1")
	}
	if targetReps > 36 {
		return 0, errors.New("formula is not accurate for more than 36 reps")
	}

	// Inverse Brzycki formula
	estimatedWeight := oneRepMax * (1.0278 - 0.0278*targetReps)

	// Round to 2 decimal places
	estimatedWeight = math.Round(estimatedWeight*100) / 100

	return estimatedWeight, nil
}
