package repositories

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentMethodRepository defines the operations available on a payment method repository
type PaymentMethodRepository interface {
	FindAll(ctx context.Context) ([]models.PaymentMethod, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.PaymentMethod, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.PaymentMethod, error)
	Create(ctx context.Context, paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error)
	Update(ctx context.Context, id primitive.ObjectID, updates bson.M) (*models.PaymentMethod, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	SetDefault(ctx context.Context, userID, paymentMethodID primitive.ObjectID) error
}

// PaymentMethodRepositoryImpl implements the PaymentMethodRepository interface
type PaymentMethodRepositoryImpl struct {
	paymentMethodDAO daos.PaymentMethodDAOInterface
}

// NewPaymentMethodRepository creates a new instance of the payment method repository
func NewPaymentMethodRepository(paymentMethodDAO daos.PaymentMethodDAOInterface) *PaymentMethodRepositoryImpl {
	return &PaymentMethodRepositoryImpl{
		paymentMethodDAO: paymentMethodDAO,
	}
}

// FindAll retrieves all payment methods
func (r *PaymentMethodRepositoryImpl) FindAll(ctx context.Context) ([]models.PaymentMethod, error) {
	return r.paymentMethodDAO.GetAllPaymentMethods(ctx)
}

// FindByID retrieves a payment method by its ID
func (r *PaymentMethodRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*models.PaymentMethod, error) {
	return r.paymentMethodDAO.GetPaymentMethodByID(ctx, id)
}

// FindByUserID retrieves payment methods by user ID
func (r *PaymentMethodRepositoryImpl) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.PaymentMethod, error) {
	return r.paymentMethodDAO.GetPaymentMethodsByUserID(ctx, userID)
}

// Create creates a new payment method
func (r *PaymentMethodRepositoryImpl) Create(ctx context.Context, paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error) {
	result, err := r.paymentMethodDAO.CreatePaymentMethod(ctx, paymentMethod)
	if err != nil {
		return nil, err
	}

	// Get the created payment method with its MongoDB ID
	createdPaymentMethod, err := r.paymentMethodDAO.GetPaymentMethodByID(ctx, result.InsertedID.(primitive.ObjectID))
	if err != nil {
		return nil, err
	}

	return createdPaymentMethod, nil
}

// Update updates an existing payment method
func (r *PaymentMethodRepositoryImpl) Update(ctx context.Context, id primitive.ObjectID, updates bson.M) (*models.PaymentMethod, error) {
	_, err := r.paymentMethodDAO.UpdatePaymentMethod(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, id)
}

// Delete deletes a payment method by its ID
func (r *PaymentMethodRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.paymentMethodDAO.DeletePaymentMethod(ctx, id)
	return err
}

// SetDefault sets a payment method as the default for a user
func (r *PaymentMethodRepositoryImpl) SetDefault(ctx context.Context, userID, paymentMethodID primitive.ObjectID) error {
	return r.paymentMethodDAO.SetDefaultPaymentMethod(ctx, userID, paymentMethodID)
}
