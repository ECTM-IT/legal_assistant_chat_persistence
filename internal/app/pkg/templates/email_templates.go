package templates

import (
	"bytes"
	"html/template"
)

type WelcomeEmailData struct {
	Username string
	PlanName string
}

type SubscriptionEmailData struct {
	Username    string
	PlanName    string
	ExpiryDate  string
	TotalAmount float64
	InvoiceURL  string
}

type HelpRequestEmailData struct {
	FirstName            string
	LastName             string
	Phone                string
	Email                string
	Organization         string
	OrganizationLastName string
	WebAddress           string
	Message              string
}

const welcomeEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 30px; }
        .content { margin-bottom: 30px; }
        .footer { text-align: center; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Legal Assistant!</h1>
        </div>
        <div class="content">
            <p>Dear {{.Username}},</p>
            <p>Welcome to Legal Assistant! We're thrilled to have you on board with our {{.PlanName}} plan.</p>
            <p>You now have access to all the features included in your subscription. If you need any assistance, don't hesitate to reach out to our support team.</p>
            <p>Best regards,<br>The Legal Assistant Team</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply directly to this email.</p>
        </div>
    </div>
</body>
</html>
`

const subscriptionEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 30px; }
        .content { margin-bottom: 30px; }
        .details { background: #f9f9f9; padding: 15px; margin: 20px 0; }
        .button { display: inline-block; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; }
        .footer { text-align: center; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Subscription Confirmation</h1>
        </div>
        <div class="content">
            <p>Dear {{.Username}},</p>
            <p>Thank you for subscribing to our {{.PlanName}} plan!</p>
            <div class="details">
                <p><strong>Plan:</strong> {{.PlanName}}</p>
                <p><strong>Amount:</strong> €{{.TotalAmount}}</p>
                <p><strong>Valid Until:</strong> {{.ExpiryDate}}</p>
            </div>
            <p>You can view your invoice here:</p>
            <p><a href="{{.InvoiceURL}}" class="button">View Invoice</a></p>
            <p>Best regards,<br>The Legal Assistant Team</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply directly to this email.</p>
        </div>
    </div>
</body>
</html>
`

const helpRequestEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 30px; }
        .content { margin-bottom: 30px; }
        .details { background: #f9f9f9; padding: 15px; margin: 20px 0; }
        .footer { text-align: center; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>New Help Request</h1>
        </div>
        <div class="content">
            <div class="details">
                <p><strong>Name:</strong> {{.FirstName}} {{.LastName}}</p>
                <p><strong>Phone:</strong> {{.Phone}}</p>
                <p><strong>Email:</strong> {{.Email}}</p>
                <p><strong>Organization:</strong> {{.Organization}}</p>
                <p><strong>Organization Contact:</strong> {{.OrganizationLastName}}</p>
                <p><strong>Website:</strong> {{.WebAddress}}</p>
                <h3>Message:</h3>
                <p>{{.Message}}</p>
            </div>
        </div>
        <div class="footer">
            <p>This message was sent from the Legal Assistant help form.</p>
        </div>
    </div>
</body>
</html>
`

func GenerateWelcomeEmail(data WelcomeEmailData) (string, error) {
	tmpl, err := template.New("welcome").Parse(welcomeEmailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GenerateSubscriptionEmail(data SubscriptionEmailData) (string, error) {
	tmpl, err := template.New("subscription").Parse(subscriptionEmailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GenerateHelpRequestEmail(data HelpRequestEmailData) (string, error) {
	tmpl, err := template.New("helprequest").Parse(helpRequestEmailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
