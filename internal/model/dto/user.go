package dto

import "url-shortener/internal/model"

type CreateUser struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=255"`
}

type UpdateUser struct {
	Email    string `validate:"omitempty,email"`
	Password string `validate:"omitempty,min=8,max=255"`
}

func (dto *CreateUser) Model() *model.User {
	return &model.User{Email: dto.Email, Password: dto.Password}
}

func (dto *UpdateUser) Model() *model.User {
	return &model.User{Email: dto.Email, Password: dto.Password}
}
