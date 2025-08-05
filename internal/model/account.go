package model

type Account struct {
	AccountId *uint64  `json:"id" gorm:"id;primaryKey;not null"`
	Balance   *float64 `json:"balance" gorm:"balance;not null"`
}
