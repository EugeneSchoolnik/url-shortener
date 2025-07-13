package dto

import "url-shortener/internal/model"

type CreateUrl struct {
	Alias string `validate:"omitempty,ascii,max=16"`
	Link  string `validate:"required,url"`
}

type PublicUrl struct {
	Alias     string `json:"alias"`
	Link      string `json:"link"`
	TotalHits int64  `json:"totalHits"`
}

func (dto *CreateUrl) Model(userID string) *model.Url {
	return &model.Url{ID: dto.Alias, Link: dto.Link, UserID: userID}
}

func ToPublicUrl(url *model.Url) *PublicUrl {
	return &PublicUrl{Alias: url.ID, Link: url.Link, TotalHits: url.TotalHits}
}
