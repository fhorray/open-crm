package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	Domain    *string    `gorm:"unique" json:"domain"`
	IsActive  bool       `gorm:"default:false" json:"is_active"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Users     []User     `gorm:"foreignKey:OrganizationID" json:"users,omitempty"`
}

func (Organization) TableName() string {
	return "core.organizations"
}

// Invitations
type Invitation struct {
	ID             string     `db:"id"`
	OrganizationID string     `db:"organization_id"`
	InvitedBy      string     `db:"invited_by"`
	InvitedEmail   string     `db:"invited_email"`
	Code           string     `db:"code"`
	ExpiresAt      time.Time  `db:"expires_at"`
	CreatedAt      *time.Time `db:"created_at"`
}

// DTOS

// -- Organizations
type CreateOrganizationDTO struct {
	Name   string  `json:"name"`
	Domain *string `json:"domain,omitempty"`
}

type OrganizationResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name,omitempty"`
	Domain    string    `json:"domain"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// -- Invitations
type CreateInvitationDTO struct {
	OrganizationID string    `json:"organization_id"`
	InvitedBy      string    `json:"invited_by"`
	InvitedEmail   string    `json:"invited_email"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type InvitationResponseDTO struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	InvitedBy      string    `json:"invited_by"`
	InvitedEmail   string    `json:"invited_email"`
	Code           string    `json:"code"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type ValidateInvitationDTO struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
