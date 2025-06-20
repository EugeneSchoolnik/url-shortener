package repo

import (
	"url-shortener/internal/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) Update(user *model.User) error {
	// Updates func ignore nil fields
	tx := r.db.Updates(user)
	if err := tx.Error; err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *UserRepo) ById(id string) (*model.User, error) {
	var user model.User

	return &user, r.db.Where("id = ?", id).First(&user).Error
}

func (r *UserRepo) ByEmail(email string) (*model.User, error) {
	var user model.User

	return &user, r.db.Where("email = ?", email).First(&user).Error
}

func (r *UserRepo) ContextById(id string) (*model.User, error) {
	var user model.User

	// TODO: add preloads when there will be other models
	return &user, r.db.Where("id = ?", id).First(&user).Error
}

func (r *UserRepo) ContextByEmail(email string) (*model.User, error) {
	var user model.User

	// TODO: add preloads when there will be other models
	return &user, r.db.Where("email = ?", email).First(&user).Error
}
