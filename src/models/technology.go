package models

import (
	"backend/src/schemas"
)

type Technology struct {
	Model
	ImageUrl string `json:"image_url" gorm:"required"`
	Name     string `json:"name" gorm:"required"`
}

func (t Technology) ConvertToStruct() schemas.Technology {
	return schemas.Technology{
		ID:       t.ID,
		Name:     t.Name,
		ImageUrl: t.ImageUrl,
	}
}

func ConvertTechnologiesToStructArray(techs []*Technology) []schemas.Technology {
	var converted []schemas.Technology
	for _, t := range techs {
		converted = append(converted, t.ConvertToStruct())
	}
	return converted
}
