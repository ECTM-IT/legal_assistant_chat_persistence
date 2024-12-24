package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/libs"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/templates"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
)

type HelpService interface {
	SendHelpEmailFromUser(ctx context.Context, request *dtos.HelpRequestDTO) error
}

type HelpServiceImpl struct {
	mailer libs.MailerService
	logger logs.Logger
}

func NewHelpService(mailer libs.MailerService, logger logs.Logger) *HelpServiceImpl {
	return &HelpServiceImpl{
		mailer: mailer,
		logger: logger,
	}
}

func (s *HelpServiceImpl) SendHelpEmailFromUser(ctx context.Context, request *dtos.HelpRequestDTO) error {
	s.logger.Info("Attempting to send help request email")

	emailData := templates.HelpRequestEmailData{
		FirstName:            request.FirstName,
		LastName:             request.LastName,
		Phone:                request.Phone,
		Email:                request.Email,
		Organization:         request.Organization,
		OrganizationLastName: request.OrganizationLastName,
		WebAddress:           request.WebAddress,
		Message:              request.Message,
	}

	htmlBody, err := templates.GenerateHelpRequestEmail(emailData)
	if err != nil {
		s.logger.Error("Failed to generate help request email", err)
		return err
	}

	// Send to the support email address
	err = s.mailer.SendHTMLEmail(
		[]string{"contatti@legalassist.it"},
		"New Help Request from "+request.FirstName+" "+request.LastName,
		htmlBody,
	)
	if err != nil {
		s.logger.Error("Failed to send help request email", err)
		return err
	}

	s.logger.Info("Help request email sent successfully")
	return nil
}
