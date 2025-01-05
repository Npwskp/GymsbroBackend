package authEnums

// GenderType represents the gender enum
type GenderType string

const (
	GenderMale   GenderType = "male"
	GenderFemale GenderType = "female"
)

func GetAllGenders() []GenderType {
	return []GenderType{
		GenderMale,
		GenderFemale,
	}
}
