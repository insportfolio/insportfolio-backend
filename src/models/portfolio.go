package models

import (
	"backend/src/utils"
)

type Portfolio struct {
	Model
	Url          string               `json:"url" gorm:"required"`
	Preview      string               `json:"preview" gorm:"required"`
	Views        int                  `json:"views" gorm:"default:0"`
	ProfileId    int                  `json:"profile_id"`
	Active       bool                 `json:"active" gorm:"default:true"`
	Profile      *Profile             `json:"profile" gorm:"foreignKey:ProfileId"`
	Technologies []*Technology        `json:"technologies" gorm:"many2many:portfolio_technologies"`
	Reactions    []*PortfolioReaction `json:"-" gorm:"foreignKey:PortfolioId"`
	Comments     []*Comment           `json:"-" gorm:"foreignKey:PortfolioId"`
}

func (p *Portfolio) SetFullPreview() {
	p.Preview = utils.GetImageFullPath(p.Preview)
}
