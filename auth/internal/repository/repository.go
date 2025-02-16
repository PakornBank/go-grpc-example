package repository

import (
	"context"
	"errors"

	"github.com/PakornBank/go-grpc-example/auth/internal/model"
	"github.com/PakornBank/go-grpc-example/auth/internal/service"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, credential *model.Credential) error
	FindByEmail(ctx context.Context, email string) (*model.Credential, error)
	DeleteByID(ctx context.Context, id string) error
	// FindByID(ctx context.Context, id string) (*model.Credential, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateUser inserts a new user record into the database.
func (r *repository) CreateUser(ctx context.Context, credential *model.Credential) error {
	return r.db.WithContext(ctx).Create(credential).Error
}

// FindByEmail retrieves a user from the database by their email address.
func (r *repository) FindByEmail(ctx context.Context, email string) (*model.Credential, error) {
	var credential model.Credential

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &credential, nil
}

// DeleteByID deletes a user record from the database.
func (r *repository) DeleteByID(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.Credential{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrRecordNotFound
	}
	return nil
}

//// FindByID retrieves a user from the database by their email address.
//func (r *repository) FindByID(ctx context.Context, id string) (*model.User, error) {
//	uid, err := uuid.Parse(id)
//	if err != nil {
//		return nil, err
//	}
//
//	var user model.User
//	if err := r.db.WithContext(ctx).Where("id = ?", uid).First(&user).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, nil
//		}
//		return nil, err
//	}
//
//	return &user, nil
//}
