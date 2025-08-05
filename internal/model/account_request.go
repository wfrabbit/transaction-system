package model

type CreateAccountRequest struct {
	AccountID      *uint64  `json:"account_id"`
	InitialBalance *float64 `json:"initial_balance"`
}

type TransferRequest struct {
	SourceAccountId      *uint64  `json:"source_account_id"`
	DestinationAccountId *uint64  `json:"destination_account_id"`
	Amount               *float64 `json:"amount"`
}

type GetAccountResponse struct {
	AccountID *uint64  `json:"account_id"`
	Balance   *float64 `json:"balance"`
}
