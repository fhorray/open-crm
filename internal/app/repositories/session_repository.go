package repositories

import (
	"open-crm/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository interface {
	FindByUserAndProvider(userID uuid.UUID, provider string) (*models.Account, error)
	Create(session *models.Session) error
	FindByToken(token string) (*models.Session, error)
	DeleteSession(token string) error
	DeleteUserSessions(userId uuid.UUID) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

// Find By user and Provider
func (r *sessionRepository) FindByUserAndProvider(userID uuid.UUID, provider string) (*models.Account, error) {
	var account models.Account
	err := r.db.
		Where("user_id = ? AND provider = ?", userID, provider).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Create Session
func (r *sessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

// Find Session by Token
func (r *sessionRepository) FindByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("access_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Delete Session
func (r *sessionRepository) DeleteSession(token string) error {
	return r.db.Where("access_token = ?", token).Delete(&models.Session{}).Error
}

// Delete User Sessions
func (r *sessionRepository) DeleteUserSessions(userId uuid.UUID) error {
	return r.db.Where("user_id = ?", userId).Delete(&models.Session{}).Error
}
