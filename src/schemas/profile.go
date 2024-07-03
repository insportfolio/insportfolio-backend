package schemas

type ProfileSchema struct {
	UserId      int    `json:"user_id,omitempty" gorm:"required;unique"`
	FullName    string `json:"full_name" gorm:"required"`
	Propic      string `json:"propic"`
	Description string `json:"description"`
	Profession  string `json:"profession"`
	OpenToWork  bool   `json:"open_to_work" gorm:"default:false"`
}

type NewProfileRequest struct {
	Record struct {
		ID    string `json:"id" binding:"required"`
		Email string `json:"email" binding:"required"`
	} `json:"record"`
}
