package plans

type CreatePlanDto struct {
	UserID     string   `json:"userid" validate:"required"`
	TypeOfPlan string   `json:"typeofplan" default:"Rest"`
	DayOfWeek  string   `json:"dayofweek" validate:"required"`
	Exercise   []string `json:"exercise"`
}

type UpdatePlanDto struct {
	TypeOfPlan string   `json:"typeofplan"`
	Exercise   []string `json:"exercise"`
}
