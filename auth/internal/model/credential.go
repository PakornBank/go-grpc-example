package model

import (
	"github.com/google/uuid"
)

// Credential represent a credential record of a user.
type Credential struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" validate:"required"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-" validate:"required"`
}
