package dto

import "url-shortener/internal/model"

type CreateUrl struct {
	Link  string `validate:"required,url"`
	Alias string `validate:"omitempty,ascii,max=16"`
}

func (dto *CreateUrl) Model(userID string) *model.Url {
	return &model.Url{ID: dto.Alias, Link: dto.Link, UserID: userID}
}
