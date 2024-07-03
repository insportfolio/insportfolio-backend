package models

type Profile struct {
	Model
	UserId       string               `json:"user_id,omitempty" gorm:"required;unique"`
	Email        string               `json:"-"`
	FullName     string               `json:"full_name"`
	Propic       string               `json:"propic"`
	Description  string               `json:"description"`
	Profession   string               `json:"profession"`
	OpenToWork   bool                 `json:"open_to_work" gorm:"default:false"`
	Portfolios   []*Portfolio         `json:"portfolios,omitempty" gorm:"foreignKey:ProfileId"`
	Reactions    []*PortfolioReaction `json:"reactions,omitempty" gorm:"foreignKey:SenderId"`
	Comments     []*Comment           `json:"comments,omitempty" gorm:"foreignKey:SenderId"`
	CommentLikes []*CommentLike       `json:"comment_likes,omitempty" gorm:"foreignKey:SenderId"`
}
