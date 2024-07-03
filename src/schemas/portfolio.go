package schemas

type GetPortfolioResponse struct {
	DefaultSchema
	Url            string       `json:"url"`
	Preview        string       `json:"preview"`
	Views          int          `json:"views"`
	ProfileId      int          `json:"profile_id"`
	Active         bool         `json:"active"`
	ReactionsCount int64        `json:"reactions_count"`
	CommentsCount  int64        `json:"comments_count"`
	Technologies   []Technology `json:"technologies"`
}

type Portfolio struct {
	DefaultSchema
	Url          string        `json:"url"`
	Preview      string        `json:"preview"`
	Views        int           `json:"views"`
	ProfileId    int           `json:"profile_id"`
	Active       bool          `json:"active"`
	Technologies []*Technology `json:"technologies"`
}
