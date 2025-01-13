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
	AssistedMachine   Equipment = "Assisted (machine)"
	AssistedPartner   Equipment = "Assisted (partner)"
	BandResistive     Equipment = "Band Resistive"
	BandAssisted      Equipment = "Band-assisted"
	Barbell           Equipment = "Barbell"
	BodyWeight        Equipment = "Body Weight"
	Cable             Equipment = "Cable"
	CableStandingFly  Equipment = "Cable Standing Fly"
	CablePullSide     Equipment = "Cable (pull side)"
	Dumbbell          Equipment = "Dumbbell"
	Isometric         Equipment = "Isometric"
	Lever             Equipment = "Lever"
	LeverPlateLoaded  Equipment = "Lever (plate loaded)"
	LeverSelectorized Equipment = "Lever (selectorized)"
	Plyometric        Equipment = "Plyometric"
	SelfAssisted      Equipment = "Self-assisted"
	Sled              Equipment = "Sled"
	SledPlateLoaded   Equipment = "Sled (plate loaded)"
	SledSelectorized  Equipment = "Sled (selectorized)"
	Smith             Equipment = "Smith"
	Suspended         Equipment = "Suspended"
	Suspension        Equipment = "Suspension"
	Weighted          Equipment = "Weighted"
	WeightedChestDip  Equipment = "Weighted Chest Dip"
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
		Assisted, AssistedMachine, AssistedPartner, BandResistive, BandAssisted,
		Barbell, BodyWeight, Cable, CableStandingFly, CablePullSide, Dumbbell,
		Isometric, Lever, LeverPlateLoaded, LeverSelectorized, Plyometric,
		SelfAssisted, Sled, SledPlateLoaded, SledSelectorized, Smith,
		Suspended, Suspension, Weighted, WeightedChestDip,
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
	case Assisted, AssistedMachine, AssistedPartner, BandResistive, BandAssisted,
		Barbell, BodyWeight, Cable, CableStandingFly, CablePullSide, Dumbbell,
		Isometric, Lever, LeverPlateLoaded, LeverSelectorized, Plyometric,
		SelfAssisted, Sled, SledPlateLoaded, SledSelectorized, Smith,
		Suspended, Suspension, Weighted, WeightedChestDip:
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
