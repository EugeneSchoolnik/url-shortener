package dto

import "url-shortener/internal/model"

type CreateUser struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}

type UpdateUser struct {
	Email string `validate:"omitempty,email"`
}

type PublicUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (dto *CreateUser) Model() *model.User {
	return &model.User{Email: dto.Email, Password: dto.Password}
}

func (dto *UpdateUser) Model() *model.User {
	return &model.User{Email: dto.Email}
}

func ToPublicUser(u *model.User) *PublicUser {
	return &PublicUser{ID: u.ID, Email: u.Email}
}
