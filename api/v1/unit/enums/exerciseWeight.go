package unitEnums

type ExerciseWeightUnit string

const (
	ExerciseWeightUnitPound ExerciseWeightUnit = "lbs"
	ExerciseWeightUnitKg    ExerciseWeightUnit = "kg"
)

func GetAllExerciseWeightUnit() []ExerciseWeightUnit {
	return []ExerciseWeightUnit{ExerciseWeightUnitPound, ExerciseWeightUnitKg}
}
