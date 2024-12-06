package daos

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CaseDAOInterface defines the interface for the CaseDAO
type CaseDAOInterface interface {
	FindAll(ctx context.Context) ([]models.Case, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (models.Case, error)
	FindByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]models.Case, error)
	Create(ctx context.Context, caseRequest *models.Case) (*mongo.InsertOneResult, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	AddCollaborator(ctx context.Context, caseID primitive.ObjectID, collaborator map[string]interface{}) (*mongo.UpdateResult, error)
	RemoveCollaborator(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error)
}

// CaseDAO implements the CaseDAOInterface
type CaseDAO struct {
	collection *mongo.Collection
	logger     logs.Logger
}

// NewCaseDAO creates a new CaseDAO
func NewCaseDAO(db *mongo.Database, logger logs.Logger) *CaseDAO {
	return &CaseDAO{
		collection: db.Collection("cases"),
		logger:     logger,
	}
}

// FindAll retrieves all cases from the database
func (dao *CaseDAO) FindAll(ctx context.Context) ([]models.Case, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve all cases")
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve cases", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var cases []models.Case
	if err := cursor.All(ctx, &cases); err != nil {
		dao.logger.Error("DAO Level: Failed to decode cases", err)
		return nil, err
	}

	dao.logger.Info("DAO Level: Successfully retrieved all cases")
	return cases, nil
}

// FindByID retrieves a case by its ID from the database
func (dao *CaseDAO) FindByID(ctx context.Context, id primitive.ObjectID) (models.Case, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve case by ID")
	var caseResponse models.Case
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&caseResponse)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Warn("Case not found")
			return models.Case{}, errors.New("case not found")
		}
		dao.logger.Error("DAO Level: Failed to retrieve case", err)
		return models.Case{}, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved case")
	return caseResponse, nil
}

// FindByCreatorID retrieves cases by creator ID from the database
func (dao *CaseDAO) FindByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]models.Case, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve cases by creator ID")
	cursor, err := dao.collection.Find(ctx, bson.M{"creator_id": creatorID})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve cases by creator ID", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var cases []models.Case
	if err := cursor.All(ctx, &cases); err != nil {
		dao.logger.Error("DAO Level: Failed to decode cases", err)
		return nil, err
	}

	dao.logger.Info("DAO Level: Successfully retrieved cases by creator ID")
	return cases, nil
}

// Create creates a new case in the database
func (dao *CaseDAO) Create(ctx context.Context, caseRequest *models.Case) (*mongo.InsertOneResult, error) {
	dao.logger.Info("DAO Level: Attempting to create new case")
	result, err := dao.collection.InsertOne(ctx, caseRequest)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to create case", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully created new case")
	return result, nil
}

// Update updates an existing case in the database
func (dao *CaseDAO) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to update case")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to update case", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully updated case")
	return result, nil
}

// Delete deletes a case by its ID from the database
func (dao *CaseDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	dao.logger.Info("DAO Level: Attempting to delete case")
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to delete case", err)
		return err
	}
	dao.logger.Info("DAO Level: Successfully deleted case")
	return nil
}

// AddCollaborator adds a collaborator to a case in the database
func (dao *CaseDAO) AddCollaborator(ctx context.Context, caseID primitive.ObjectID, collaborator map[string]interface{}) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to add collaborator to case")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": caseID}, bson.M{"$addToSet": bson.M{"collaborator_ids": collaborator}})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to add collaborator to case", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully added collaborator to case")
	return result, nil
}

// RemoveCollaborator removes a collaborator from a case in the database
func (dao *CaseDAO) RemoveCollaborator(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to remove collaborator from case")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": caseID}, bson.M{"$pull": bson.M{"collaborator_ids": collaboratorID}})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to remove collaborator from case", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully removed collaborator from case")
	return result, nil
}

// AddDocument adds a document to a case in the database
func (dao *CaseDAO) AddDocument(ctx context.Context, caseID primitive.ObjectID, document *models.Document) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to add document to case")

	// Generate new ID for the document
	document.ID = primitive.NewObjectID()

	// Set the upload date for the document
	document.UploadDate = time.Now()
	document.ModifiedDate = time.Now()

	// Update the case by adding the document to the documents array
	result, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": caseID},
		bson.M{"$addToSet": bson.M{"documents": document}},
	)

	if err != nil {
		dao.logger.Error("DAO Level: Failed to add document to case", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully added document to case")
	return result, nil
}

func (dao *CaseDAO) UpdateDocument(ctx context.Context, caseID primitive.ObjectID, documentID primitive.ObjectID, updatedDocument *models.Document) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to update document in case")

	// Set the modified date for the document
	updatedDocument.ModifiedDate = time.Now()

	// Build the query to match the case ID and the document ID within the documents array
	filter := bson.M{
		"_id":           caseID,
		"documents._id": documentID,
	}

	// Dynamically construct the update document based on non-empty fields in `updatedDocument`
	updateFields := bson.M{}
	docValue := reflect.ValueOf(*updatedDocument)
	docType := reflect.TypeOf(*updatedDocument)

	for i := 0; i < docValue.NumField(); i++ {
		field := docType.Field(i)
		bsonTag := field.Tag.Get("bson")

		// Skip fields with no bson tag or `omitempty` that are zero-valued
		if bsonTag == "" || bsonTag == "-" {
			continue
		}

		fieldValue := docValue.Field(i)
		zeroValue := reflect.Zero(field.Type)

		// Check if the field is zero (handles slices, structs, etc.)
		if isZeroValue(fieldValue, zeroValue) {
			continue
		}

		// Add non-zero value to the update fields map
		updateFields["documents.$."+bsonTag] = fieldValue.Interface()
	}

	if len(updateFields) == 0 {
		dao.logger.Warn("DAO Level: No fields to update")
		return nil, fmt.Errorf("no fields to update")
	}

	// Build the update operation
	update := bson.M{
		"$set": updateFields,
	}

	// Perform the update operation
	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to update document in case", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully updated document in case")
	return result, nil
}

// AddDocumentCollaborator adds a collaborator to a document in a case
func (dao *CaseDAO) AddDocumentCollaborator(ctx context.Context, caseID primitive.ObjectID, documentID primitive.ObjectID, collaborator *models.DocumentCollaborator) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to add collaborator to document in case")

	// Build the query to match the case ID and the specific document ID
	filter := bson.M{
		"_id":           caseID,
		"documents._id": documentID,
	}

	// Build the update operation to add the collaborator to the document's collaborators array
	update := bson.M{
		"$addToSet": bson.M{
			"documents.$.collaborators": collaborator,
		},
	}

	// Perform the update operation
	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to add collaborator to document in case", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully added collaborator to document in case")
	return result, nil
}

// DeleteDocument removes a document from a case in the database
func (dao *CaseDAO) DeleteDocument(ctx context.Context, caseID, documentID primitive.ObjectID) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to delete document from case")

	// Update the case by removing the document with the given documentID from the documents array
	result, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": caseID},
		bson.M{"$pull": bson.M{"documents": bson.M{"_id": documentID}}},
	)

	if err != nil {
		dao.logger.Error("DAO Level: Failed to delete document from case", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully deleted document from case")
	return result, nil
}

// AddFeedback adds a feedback entry to a specific message within a case in the database
func (dao *CaseDAO) AddFeedback(ctx context.Context, caseID primitive.ObjectID, messageID string, feedback models.Feedback) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to add feedback to message in case")

	// Prepare the filter to find the specific message within the case
	filter := bson.M{
		"_id":          caseID,
		"messages._id": messageID,
	}

	// Prepare the update to add the feedback to the message's feedback array
	update := bson.M{
		"$addToSet": bson.M{"messages.$.feedbacks": feedback},
	}

	// Perform the update
	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to add feedback to message in case", err)
		return nil, err
	}

	dao.logger.Info("DAO Level: Successfully added feedback to message in case")
	return result, nil
}

// GetFeedbackByUserAndMessage retrieves feedback from a specific user for a specific message in a case.
func (dao *CaseDAO) GetFeedbackByUserAndMessage(ctx context.Context, creatorID primitive.ObjectID, messageID string) ([]models.Feedback, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve feedback by user and message")

	// Define the filter for feedbacks matching the user and message IDs
	filter := bson.M{
		"messages.feedbacks": bson.M{
			"$elemMatch": bson.M{
				"creator_id": creatorID,
				"message_id": messageID,
			},
		},
	}

	// Find the case containing the specified feedback
	var caseDoc models.Case
	err := dao.collection.FindOne(ctx, filter).Decode(&caseDoc)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve feedback by user and message", err)
		return nil, err
	}

	// Extract and return the feedbacks that match
	var feedbacks []models.Feedback
	for _, msg := range caseDoc.Messages {
		if msg.ID == messageID {
			for _, fb := range msg.Feedbacks {
				if fb.CreatorID == creatorID {
					feedbacks = append(feedbacks, fb)
				}
			}
			break
		}
	}

	return feedbacks, nil
}

// Helper function to check if a field value is zero
func isZeroValue(value, zero reflect.Value) bool {
	switch value.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		return value.Len() == 0 // Zero if the length is 0
	case reflect.Struct:
		return reflect.DeepEqual(value.Interface(), zero.Interface()) // Compare structs deeply
	case reflect.Ptr, reflect.Interface:
		return value.IsNil() // Zero if the pointer or interface is nil
	default:
		return value.Interface() == zero.Interface() // Default comparison for other types
	}
}
