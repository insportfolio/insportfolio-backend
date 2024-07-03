package schemas

type PortfolioReactionRequest struct {
	Value       uint `binding:"required,min=0,max=5"`
	SenderId    uint `binding:"required,min=0"`
	PortfolioId uint `binding:"required,min=0"`
}
