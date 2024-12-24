package models

type Plan struct {
	Name        string   `json:"name" bson:"name"`
	Type        string   `json:"type" bson:"type"` // "monthly" or "annual"
	Price       float64  `json:"price" bson:"price"`
	Description string   `json:"description" bson:"description"`
	Features    []string `json:"features" bson:"features"`
}

// PredefinedPlans returns the list of available plans
func PredefinedPlans() map[string][]Plan {
	return map[string][]Plan{
		"monthly": {
			{
				Name:        "pro",
				Type:        "monthly",
				Price:       20.00,
				Description: "Professional plan for individual users",
				Features:    []string{"Unlimited cases", "Priority support", "Advanced analytics"},
			},
			{
				Name:        "team",
				Type:        "monthly",
				Price:       60.00,
				Description: "Perfect for small teams",
				Features:    []string{"Everything in Pro", "Team collaboration", "Admin dashboard"},
			},
			{
				Name:        "enterprise",
				Type:        "monthly",
				Price:       0,
				Description: "For large organizations",
				Features:    []string{"Everything in Team", "Custom integrations", "Dedicated support"},
			},
		},
		"annual": {
			{
				Name:        "pro",
				Type:        "annual",
				Price:       192.00,
				Description: "Professional plan for individual users",
				Features:    []string{"Unlimited cases", "Priority support", "Advanced analytics"},
			},
			{
				Name:        "team",
				Type:        "annual",
				Price:       576.00,
				Description: "Perfect for small teams",
				Features:    []string{"Everything in Pro", "Team collaboration", "Admin dashboard"},
			},
			{
				Name:        "enterprise",
				Type:        "annual",
				Price:       0,
				Description: "For large organizations",
				Features:    []string{"Everything in Team", "Custom integrations", "Dedicated support"},
			},
		},
	}
}
