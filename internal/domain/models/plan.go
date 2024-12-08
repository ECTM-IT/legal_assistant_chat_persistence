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
				Name:        "Pro",
				Type:        "monthly",
				Price:       49.99,
				Description: "Professional plan for individual users",
				Features:    []string{"Unlimited cases", "Priority support", "Advanced analytics"},
			},
			{
				Name:        "Team",
				Type:        "monthly",
				Price:       99.99,
				Description: "Perfect for small teams",
				Features:    []string{"Everything in Pro", "Team collaboration", "Admin dashboard"},
			},
			{
				Name:        "Enterprise",
				Type:        "monthly",
				Price:       199.99,
				Description: "For large organizations",
				Features:    []string{"Everything in Team", "Custom integrations", "Dedicated support"},
			},
		},
		"annual": {
			{
				Name:        "Pro",
				Type:        "annual",
				Price:       499.99,
				Description: "Professional plan for individual users",
				Features:    []string{"Unlimited cases", "Priority support", "Advanced analytics"},
			},
			{
				Name:        "Team",
				Type:        "annual",
				Price:       999.99,
				Description: "Perfect for small teams",
				Features:    []string{"Everything in Pro", "Team collaboration", "Admin dashboard"},
			},
			{
				Name:        "Enterprise",
				Type:        "annual",
				Price:       1999.99,
				Description: "For large organizations",
				Features:    []string{"Everything in Team", "Custom integrations", "Dedicated support"},
			},
		},
	}
}
