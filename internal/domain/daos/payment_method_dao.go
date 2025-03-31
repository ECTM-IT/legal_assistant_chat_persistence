package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PaymentMethodDAOInterface defines the interface for payment method data access
type PaymentMethodDAOInterface interface {
	GetAllPaymentMethods(ctx context.Context) ([]models.PaymentMethod, error)
	GetPaymentMethodByID(ctx context.Context, id primitive.ObjectID) (*models.PaymentMethod, error)
	GetPaymentMethodsByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.PaymentMethod, error)
	CreatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) (*mongo.InsertOneResult, error)
	UpdatePaymentMethod(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error)
	DeletePaymentMethod(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
	SetDefaultPaymentMethod(ctx context.Context, userID, paymentMethodID primitive.ObjectID) error
}

// PaymentMethodDAO implements the PaymentMethodDAOInterface
type PaymentMethodDAO struct {
	collection *mongo.Collection
	logger     logs.Logger
}

// NewPaymentMethodDAO creates a new PaymentMethodDAO
func NewPaymentMethodDAO(db *mongo.Database, logger logs.Logger) *PaymentMethodDAO {
	return &PaymentMethodDAO{
		collection: db.Collection("payment_methods"),
		logger:     logger,
	}
}

// GetAllPaymentMethods retrieves all payment methods from the database
func (dao *PaymentMethodDAO) GetAllPaymentMethods(ctx context.Context) ([]models.PaymentMethod, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve all payment methods")
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve payment methods", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var paymentMethods []models.PaymentMethod
	if err := cursor.All(ctx, &paymentMethods); err != nil {
		dao.logger.Error("DAO Level: Failed to decode payment methods", err)
		return nil, err
	}

	dao.logger.Info("DAO Level: Successfully retrieved all payment methods")
	return paymentMethods, nil
}

// GetPaymentMethodByID retrieves a payment method by its ID from the database
func (dao *PaymentMethodDAO) GetPaymentMethodByID(ctx context.Context, id primitive.ObjectID) (*models.PaymentMethod, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve payment method by ID")
	var paymentMethod models.PaymentMethod
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&paymentMethod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Warn("Payment method not found")
			return nil, errors.New("payment method not found")
		}
		dao.logger.Error("DAO Level: Failed to retrieve payment method", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved payment method")
	return &paymentMethod, nil
}

// GetPaymentMethodsByUserID retrieves payment methods by user ID from the database
func (dao *PaymentMethodDAO) GetPaymentMethodsByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.PaymentMethod, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve payment methods by user ID")
	cursor, err := dao.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve payment methods by user ID", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var paymentMethods []models.PaymentMethod
	if err := cursor.All(ctx, &paymentMethods); err != nil {
		dao.logger.Error("DAO Level: Failed to decode payment methods", err)
		return nil, err
	}

	dao.logger.Info("DAO Level: Successfully retrieved payment methods by user ID")
	return paymentMethods, nil
}

// CreatePaymentMethod creates a new payment method in the database
func (dao *PaymentMethodDAO) CreatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) (*mongo.InsertOneResult, error) {
	dao.logger.Info("DAO Level: Attempting to create new payment method")
	result, err := dao.collection.InsertOne(ctx, paymentMethod)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to create payment method", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully created new payment method")
	return result, nil
}

// UpdatePaymentMethod updates an existing payment method in the database
func (dao *PaymentMethodDAO) UpdatePaymentMethod(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to update payment method")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to update payment method", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully updated payment method")
	return result, nil
}

// DeletePaymentMethod deletes a payment method by its ID from the database
func (dao *PaymentMethodDAO) DeletePaymentMethod(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	dao.logger.Info("DAO Level: Attempting to delete payment method")
	result, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to delete payment method", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully deleted payment method")
	return result, nil
}

// SetDefaultPaymentMethod sets a payment method as default and unsets others for the user
func (dao *PaymentMethodDAO) SetDefaultPaymentMethod(ctx context.Context, userID, paymentMethodID primitive.ObjectID) error {
	dao.logger.Info("DAO Level: Setting default payment method")

	// Start a session for transaction
	session, err := dao.collection.Database().Client().StartSession()
	if err != nil {
		dao.logger.Error("DAO Level: Failed to start session", err)
		return err
	}
	defer session.EndSession(ctx)

	// Run transaction
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// First, unset all payment methods as default for this user
		_, err := dao.collection.UpdateMany(
			sessionContext,
			bson.M{
				"user_id": userID,
			},
			bson.M{
				"$set": bson.M{
					"is_default": false,
				},
			},
		)
		if err != nil {
			dao.logger.Error("DAO Level: Failed to unset previous default payment methods", err)
			return err
		}

		// Then, set the specified payment method as default
		opts := options.Update().SetUpsert(false)
		result, err := dao.collection.UpdateOne(
			sessionContext,
			bson.M{
				"_id":     paymentMethodID,
				"user_id": userID,
			},
			bson.M{
				"$set": bson.M{
					"is_default": true,
				},
			},
			opts,
		)
		if err != nil {
			dao.logger.Error("DAO Level: Failed to set payment method as default", err)
			return err
		}

		if result.MatchedCount == 0 {
			dao.logger.Error("DAO Level: No payment method found to set as default", errors.New("payment method not found"))
			return errors.New("payment method not found")
		}

		return nil
	})

	if err != nil {
		return err
	}

	dao.logger.Info("DAO Level: Successfully set default payment method")
	return nil
}
