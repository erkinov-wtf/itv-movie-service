package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"itv-movie/internal/models"
	"itv-movie/internal/storage/database"
	"time"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(postgres *database.PostgresDB) *SessionRepository {
	return &SessionRepository{
		db: postgres.DB,
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) (*models.Session, error) {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) GetByToken(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	if err := r.db.WithContext(ctx).Where("token = ? AND is_revoked = false", token).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	var sessions []*models.Session
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now()).
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) RevokeByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&models.Session{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}

func (r *SessionRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Session{}).
		Where("user_id = ? AND is_revoked = false", userID).
		Update("is_revoked", true).Error
}

func (r *SessionRepository) GetByAccessToken(ctx context.Context, accessToken string) (*models.Session, error) {
	var session models.Session
	if err := r.db.WithContext(ctx).Where("access_token = ? AND is_revoked = false", accessToken).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var session models.Session
	if err := r.db.WithContext(ctx).Where("refresh_token = ? AND is_revoked = false", refreshToken).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) RevokeByID(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("is_revoked", true).Error
}
