package dtos

type HelpRequestDTO struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Phone                string `json:"phone"`
	Email                string `json:"email"`
	Organization         string `json:"organization"`
	OrganizationLastName string `json:"organization_last_name"`
	WebAddress           string `json:"web_address"`
	Message              string `json:"message"`
}
