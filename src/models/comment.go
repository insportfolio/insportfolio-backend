package models

type Comment struct {
	Model
	SenderId    int            `json:"sender_id"`
	Sender      *Profile       `json:"sender" gorm:"foreignKey:SenderId"`
	PortfolioId int            `json:"portfolio_id"`
	Portfolio   *Portfolio     `json:"portfolio" gorm:"foreignKey:PortfolioId"`
	AnswerId    *int           `json:"answer_id" gorm:"foreignKey:AnswerId"`
	Answer      *Comment       `json:"answer"`
	Likes       []*CommentLike `json:"likes" gorm:"foreignKey:CommentId"`
	Text        string         `json:"text" gorm:"not null"`
}

type CommentLike struct {
	Model
	SenderId  int      `json:"sender_id"`
	Sender    *Profile `json:"sender" gorm:"foreignKey:SenderId"`
	CommentId int      `json:"comment_id"`
	Comment   *Comment `json:"comment" gorm:"foreignKey:CommentId"`
}
