package exercise

import (
	"fmt"
	"time"

	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Add this helper function at package level
func debugType(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

type Exercise struct {
	ID           string                       `json:"id,omitempty" bson:"_id,omitempty"`
	UserID       string                       `json:"userid" bson:"userid" validate:"required"`
	Name         string                       `json:"name" validate:"required"`
	Equipment    exerciseEnums.Equipment      `json:"equipment" validate:"required"`
	Mechanics    exerciseEnums.Mechanics      `json:"mechanics" validate:"required"`
	Force        exerciseEnums.Force          `json:"force" validate:"required"`
	Preparation  []string                     `json:"preparation" validate:"required"`
	Execution    []string                     `json:"execution" validate:"required"`
	Image        string                       `json:"image"`
	BodyPart     []exerciseEnums.BodyPart     `json:"body_part" bson:"body_part" validate:"required"`
	TargetMuscle []exerciseEnums.TargetMuscle `json:"target_muscle" bson:"target_muscle" validate:"required"`
	CreatedAt    time.Time                    `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time                    `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt    time.Time                    `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// MarshalBSON implements the bson.Marshaler interface.
func (e Exercise) MarshalBSON() ([]byte, error) {
	type ExerciseAlias Exercise // prevent recursive marshaling

	// Convert enums to strings for BSON
	bodyPartStrings := make([]string, len(e.BodyPart))
	for i, bp := range e.BodyPart {
		bodyPartStrings[i] = string(bp)
	}

	targetMuscleStrings := make([]string, len(e.TargetMuscle))
	for i, tm := range e.TargetMuscle {
		targetMuscleStrings[i] = string(tm)
	}

	return bson.Marshal(struct {
		ExerciseAlias `bson:",inline"`
		BodyPart      []string `bson:"body_part"`
		TargetMuscle  []string `bson:"target_muscle"`
	}{
		ExerciseAlias: ExerciseAlias(e),
		BodyPart:      bodyPartStrings,
		TargetMuscle:  targetMuscleStrings,
	})
}

// UnmarshalBSON implements the bson.Unmarshaler interface.
func (e *Exercise) UnmarshalBSON(data []byte) error {
	type ExerciseAlias Exercise // prevent recursive unmarshaling

	// First try with string arrays
	temp := struct {
		ID            primitive.ObjectID `bson:"_id,omitempty"`
		ExerciseAlias `bson:",inline"`
		BodyPart      interface{} `bson:"body_part"`
		TargetMuscle  interface{} `bson:"target_muscle"`
	}{}

	if err := bson.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Convert ID to string
	if !temp.ID.IsZero() {
		e.ID = temp.ID.Hex()
	}

	// Copy all other fields
	*e = Exercise(temp.ExerciseAlias)

	// Handle BodyPart conversion
	switch v := temp.BodyPart.(type) {
	case []interface{}:
		e.BodyPart = make([]exerciseEnums.BodyPart, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				e.BodyPart[i] = exerciseEnums.BodyPart(str)
			}
		}
	case string:
		e.BodyPart = []exerciseEnums.BodyPart{exerciseEnums.BodyPart(v)}
	case nil:
		e.BodyPart = []exerciseEnums.BodyPart{}
	default:
		fmt.Printf("Unexpected type for BodyPart: %T\n", v)
		e.BodyPart = []exerciseEnums.BodyPart{}
	}

	// Handle TargetMuscle conversion
	switch v := temp.TargetMuscle.(type) {
	case primitive.A: // Handle MongoDB array type directly
		e.TargetMuscle = make([]exerciseEnums.TargetMuscle, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				e.TargetMuscle[i] = exerciseEnums.TargetMuscle(str)
			}
		}
	case []interface{}:
		e.TargetMuscle = make([]exerciseEnums.TargetMuscle, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				e.TargetMuscle[i] = exerciseEnums.TargetMuscle(str)
			}
		}
	case string:
		e.TargetMuscle = []exerciseEnums.TargetMuscle{exerciseEnums.TargetMuscle(v)}
	case nil:
		e.TargetMuscle = []exerciseEnums.TargetMuscle{}
	default:
		fmt.Printf("Unexpected type for TargetMuscle: %T\n", v)
		e.TargetMuscle = []exerciseEnums.TargetMuscle{}
	}

	return nil
}
