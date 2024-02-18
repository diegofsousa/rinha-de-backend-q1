package api

type (
	TransactionRequest struct {
		Value       int64  `json:"valor"`
		Type        string `json:"tipo"`
		Description string `json:"descricao"`
	}
	TransactionResponse struct {
		Limit   int64 `json:"limite"`
		Balance int64 `json:"saldo"`
	}
)
