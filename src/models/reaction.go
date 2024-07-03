package models

type PortfolioReaction struct {
	Model
	SenderId    int        `json:"-"`
	Sender      *Profile   `json:"sender" gorm:"foreignKey:SenderId"`
	PortfolioId int        `json:"portfolio_id"`
	Portfolio   *Portfolio `json:"-" gorm:"foreignKey:PortfolioId"`
	Value       int        `json:"value" gorm:"default:0"` // between 1 and 5
}
