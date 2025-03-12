package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"itv-movie/internal/models"
	"itv-movie/internal/pkg/utils/constants"
	"itv-movie/internal/storage/database"
	"time"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(postgres *database.PostgresDB) *UserRepository {
	return &UserRepository{
		db: postgres.DB,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", now).Error
}

func (r *UserRepository) UpdateStatus(ctx context.Context, userID uuid.UUID, active bool) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("active", active).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *UserRepository) FindAdmin(ctx context.Context) (int64, error) {
	var adminCount int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", constants.AdminRole).Count(&adminCount).Error; err != nil {
		return 0, err
	}

	return adminCount, nil
}
