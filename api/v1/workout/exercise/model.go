package exercise

type Exercise struct {
	ID          string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Type        []string `json:"type" validate:"required"`
	Muscle      []string `json:"muscle" validate:"required"`
	Image       string   `json:"image" validate:"required"`
}
