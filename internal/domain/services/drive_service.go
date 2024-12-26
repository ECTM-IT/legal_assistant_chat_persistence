package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	fileURLstring         = "https://drive.google.com/file/d/%s/view"
	defaultPermissionType = "user"
	defaultPermissionRole = "reader"
)

var fileMimeTypes = map[string]string{
	// Microsoft Office Formats
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",

	// Google Workspace Formats
	"google-docs":     "application/vnd.google-apps.document",
	"google-sheets":   "application/vnd.google-apps.spreadsheet",
	"google-slides":   "application/vnd.google-apps.presentation",
	"google-forms":    "application/vnd.google-apps.form",
	"google-drawings": "application/vnd.google-apps.drawing",

	// PDF and Text Files
	".pdf": "application/pdf",
	".txt": "text/plain",
	".rtf": "application/rtf",

	// Image Formats
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".bmp":  "image/bmp",

	// Compressed and Archive Formats
	".zip": "application/zip",
	".gz":  "application/gzip",
	".tar": "application/x-tar",

	// Other Document Formats
	".md":   "text/markdown",
	".xml":  "application/xml",
	".json": "application/json",
	".csv":  "text/csv",
}

// DriveService defines the operations available for managing Google Drive files.
type DriveService interface {
	UploadDocumentToDrive(ctx context.Context, document *models.Document, userEmail string) (*drive.File, string, error)
	ShareDocument(ctx context.Context, userEmail string) error
}

// DriveServiceImpl implements the DriveService interface.
type DriveServiceImpl struct {
	driveService *drive.Service
	logger       logs.Logger
}

// NewDriveService creates a new instance of the DriveService.
func NewDriveService(gcpSaKeyPath string, logger logs.Logger) *DriveServiceImpl {
	// Create a context
	ctx := context.Background()

	// Authenticate using the service account
	srv, err := drive.NewService(ctx, option.WithCredentialsFile(gcpSaKeyPath))
	if err != nil {
		logger.Error("Unable to create Drive client: %v", err)
	}

	return &DriveServiceImpl{
		driveService: srv,
		logger:       logger,
	}
}

// UploadDocumentToDrive uploads a single document into Google Drive.
func (s *DriveServiceImpl) UploadDocumentToDrive(ctx context.Context, document *models.Document, userEmail string) (*drive.File, string, error) {
	// Define the file metadata
	driveFile := &drive.File{
		Name:     document.FileName,
		MimeType: fileMimeTypes[getFileExtension(document.FileName)],
	}

	file, err := base64toIOReader(document.FileContent)
	if err != nil {
		s.logger.Error("Unable to convert file from base64 into []byte: %v", err)
		return nil, "", err
	}

	// Upload the file
	uploadedFile, err := s.driveService.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		s.logger.Error("Unable to upload file: %v", err)
		return nil, "", err
	}
	s.logger.Info("File uploaded successfully!")

	// userEmail comes points to the owner of this case and is used here to identify
	// the user to share the current file with.
	if userEmail != "" {
		// Permission to user's email
		permission := &drive.Permission{
			Type:         defaultPermissionType, // "anyone" to make file accessibel to anyone with the link
			Role:         defaultPermissionRole, // options: reader, commenter, writer
			EmailAddress: userEmail,
		}

		// Create permission
		_, err = s.driveService.Permissions.Create(uploadedFile.Id, permission).Do()
		if err != nil {
			s.logger.Error("Failed to share file: %v", err)
		}
		s.logger.Info("File shared successfully!")
	}

	fileURL := fmt.Sprintf(fileURLstring, uploadedFile.Id)

	return uploadedFile, fileURL, nil
}

// ShareDocument adds collaborator emails to the google drive file.
func (s *DriveServiceImpl) ShareDocument(ctx context.Context, fileID string, userEmail string) error {
	// Permission to user's email
	permission := &drive.Permission{
		Type:         defaultPermissionType, // "anyone" to make file accessibel to anyone with the link
		Role:         defaultPermissionRole, // options: reader, commenter, writer
		EmailAddress: userEmail,
	}

	// Create permission
	_, err := s.driveService.Permissions.Create(fileID, permission).Do()
	if err != nil {
		s.logger.Error("Failed to share file: %v", err)
	}
	s.logger.Info("File shared successfully!")
	return nil
}

func getFileExtension(fileName string) string {
	// Use strings.LastIndex to find the last dot in the file name
	index := strings.LastIndex(fileName, ".")
	if index == -1 || index == len(fileName)-1 {
		// Return an empty string if there's no extension or the dot is the last character
		return ""
	}
	return fileName[index+1:]
}
