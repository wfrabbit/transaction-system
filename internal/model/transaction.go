package model

import "time"

const (
	TransactionCompleted = "completed"
	TransactionFailed    = "failed"
)

type Transaction struct {
	Id                 *uint64   `json:"id" gorm:"id;primaryKey;not null"`
	SourceAccount      *uint64   `json:"source_account" gorm:"source_account;index;not null"`
	DestinationAccount *uint64   `json:"destination_account" gorm:"destination_account;not null"`
	Amount             *float64  `json:"amount" gorm:"amount;not null"`
	CreatedAt          time.Time `json:"created_at"`
}
