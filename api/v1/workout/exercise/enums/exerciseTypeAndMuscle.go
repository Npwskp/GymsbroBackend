package exerciseEnums

import (
	"fmt"
)

type Equipment string
type Mechanics string
type Force string
type BodyPart string
type TargetMuscle string

const (
	// Equipment
	Assisted          Equipment = "Assisted"
	BandAssisted      Equipment = "Band-assisted"
	Barbell           Equipment = "Barbell"
	BodyWeight        Equipment = "Body Weight"
	Cable             Equipment = "Cable"
	Dumbbell          Equipment = "Dumbbell"
	Lever             Equipment = "Lever"
	LeverPlateLoaded  Equipment = "Lever (plate loaded)"
	LeverSelectorized Equipment = "Lever (selectorized)"
	Plyometric        Equipment = "Plyometric"
	SelfAssisted      Equipment = "Self-assisted"
	Sled              Equipment = "Sled"
	Smith             Equipment = "Smith"
	Suspended         Equipment = "Suspended"
	Suspension        Equipment = "Suspension"
	Weighted          Equipment = "Weighted"
)

const (
	// Mechanics
	Compound Mechanics = "Compound"
	Isolated Mechanics = "Isolated"
)

const (
	// Force
	Push        Force = "Push"
	Pull        Force = "Pull"
	PushAndPull Force = "Push & Pull"
)

const (
	// Target Muscles
	Adductors              TargetMuscle = "Adductors"
	AnteriorDeltoid        TargetMuscle = "Anterior Deltoid"
	BicepsBrachii          TargetMuscle = "Biceps Brachii"
	Brachialis             TargetMuscle = "Brachialis"
	Brachioradialis        TargetMuscle = "Brachioradialis"
	ErectorSpinae          TargetMuscle = "Erector Spinae"
	Gastrocnemius          TargetMuscle = "Gastrocnemius"
	GluteusMaximus         TargetMuscle = "Gluteus Maximus"
	Hamstrings             TargetMuscle = "Hamstrings"
	HipFlexors             TargetMuscle = "Hip Flexors"
	Iliopsoas              TargetMuscle = "Iliopsoas"
	Infraspinatus          TargetMuscle = "Infraspinatus"
	LateralDeltoid         TargetMuscle = "Lateral Deltoid"
	LatissimusDorsi        TargetMuscle = "Latissimus Dorsi"
	LevatorScapulae        TargetMuscle = "Levator Scapulae"
	LowerTrapezius         TargetMuscle = "Lower Trapezius"
	MiddleTrapezius        TargetMuscle = "Middle Trapezius"
	Obliques               TargetMuscle = "Obliques"
	PectoralisMajorClav    TargetMuscle = "Pectoralis Major Clavicular"
	PectoralisMajorSternal TargetMuscle = "Pectoralis Major Sternal"
	PectoralisMinor        TargetMuscle = "Pectoralis Minor"
	PosteriorDeltoid       TargetMuscle = "Posterior Deltoid"
	QuadratusLumborum      TargetMuscle = "Quadratus Lumborum"
	Quadriceps             TargetMuscle = "Quadriceps"
	RectusAbdominis        TargetMuscle = "Rectus Abdominis"
	Rhomboids              TargetMuscle = "Rhomboids"
	SerratusAnterior       TargetMuscle = "Serratus Anterior"
	Soleus                 TargetMuscle = "Soleus"
	Splenius               TargetMuscle = "Splenius"
	Sternocleidomastoid    TargetMuscle = "Sternocleidomastoid"
	Subscapularis          TargetMuscle = "Subscapularis"
	Supraspinatus          TargetMuscle = "Supraspinatus"
	TeresMinor             TargetMuscle = "Teres Minor"
	TibialisAnterior       TargetMuscle = "Tibialis Anterior"
	Trapezius              TargetMuscle = "Trapezius"
	TricepsBrachii         TargetMuscle = "Triceps Brachii"
	UpperTrapezius         TargetMuscle = "Upper Trapezius"
)

const (
	// Body Parts
	Back      BodyPart = "Back"
	Calves    BodyPart = "Calves"
	Chest     BodyPart = "Chest"
	Forearm   BodyPart = "Forearm"
	Hips      BodyPart = "Hips"
	Neck      BodyPart = "Neck"
	Shoulder  BodyPart = "Shoulder"
	Thighs    BodyPart = "Thighs"
	UpperArms BodyPart = "Upper Arms"
)

func GetAllEquipment() []Equipment {
	return []Equipment{
		Assisted, BandAssisted, Barbell, BodyWeight, Cable, Dumbbell,
		Lever, LeverPlateLoaded, LeverSelectorized, Plyometric,
		SelfAssisted, Sled, Smith, Suspended, Suspension, Weighted,
	}
}

func GetAllMechanics() []Mechanics {
	return []Mechanics{Compound, Isolated}
}

func GetAllForces() []Force {
	return []Force{Push, Pull, PushAndPull}
}

func GetAllBodyParts() []BodyPart {
	return []BodyPart{Back, Calves, Chest, Forearm, Hips, Neck, Shoulder, Thighs, UpperArms}
}

func GetAllTargetMuscles() []TargetMuscle {
	return []TargetMuscle{
		Adductors, AnteriorDeltoid, BicepsBrachii, Brachialis, Brachioradialis,
		ErectorSpinae, Gastrocnemius, GluteusMaximus, Hamstrings, HipFlexors,
		Iliopsoas, Infraspinatus, LateralDeltoid, LatissimusDorsi, LevatorScapulae,
		LowerTrapezius, MiddleTrapezius, Obliques, PectoralisMajorClav,
		PectoralisMajorSternal, PectoralisMinor, PosteriorDeltoid, QuadratusLumborum,
		Quadriceps, RectusAbdominis, Rhomboids, SerratusAnterior, Soleus,
		Splenius, Sternocleidomastoid, Subscapularis, Supraspinatus, TeresMinor,
		TibialisAnterior, Trapezius, TricepsBrachii, UpperTrapezius,
	}
}

func ParseEquipment(s string) (Equipment, error) {
	switch Equipment(s) {
	case Assisted, BandAssisted, Barbell, BodyWeight, Cable, Dumbbell,
		Lever, LeverPlateLoaded, LeverSelectorized, Plyometric,
		SelfAssisted, Sled, Smith, Suspended, Suspension, Weighted:
		return Equipment(s), nil
	default:
		return "", fmt.Errorf("invalid equipment: %s", s)
	}
}

func ParseMechanics(s string) (Mechanics, error) {
	switch Mechanics(s) {
	case Compound, Isolated:
		return Mechanics(s), nil
	default:
		return "", fmt.Errorf("invalid mechanics: %s", s)
	}
}

func ParseForce(s string) (Force, error) {
	switch Force(s) {
	case Push, Pull, PushAndPull:
		return Force(s), nil
	default:
		return "", fmt.Errorf("invalid force: %s", s)
	}
}

func ParseBodyPart(s string) (BodyPart, error) {
	switch BodyPart(s) {
	case Back, Calves, Chest, Forearm, Hips, Neck, Shoulder, Thighs, UpperArms:
		return BodyPart(s), nil
	default:
		return "", fmt.Errorf("invalid body part: %s", s)
	}
}

func ParseTargetMuscle(s string) (TargetMuscle, error) {
	for _, v := range GetAllTargetMuscles() {
		if string(v) == s {
			return v, nil
		}
	}
	return "", fmt.Errorf("invalid target muscle: %s", s)
}

// Generic validation function for any string-based enum type
func isValidEnum[T ~string](value string, getAllFn func() []T) bool {
	for _, v := range getAllFn() {
		if string(v) == value {
			return true
		}
	}
	return false
}

// Specific validation functions using the generic function
func IsValidEquipment(value string) bool {
	return isValidEnum(value, GetAllEquipment)
}

func IsValidMechanics(value string) bool {
	return isValidEnum(value, GetAllMechanics)
}

func IsValidForce(value string) bool {
	return isValidEnum(value, GetAllForces)
}

func IsValidBodyPart(value string) bool {
	return isValidEnum(value, GetAllBodyParts)
}

func IsValidTargetMuscle(value string) bool {
	return isValidEnum(value, GetAllTargetMuscles)
}
