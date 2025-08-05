package repository

import (
	"fmt"
	"wangfeng/transaction-system/internal/db"
	"wangfeng/transaction-system/internal/model"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *model.Account) error
	GetByID(id uint64) (*model.Account, error)
	Transfer(sourceAccountId, destinationAccountId uint64, amount *float64) error
}

type AccountRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepository() (*AccountRepositoryImpl, error) {
	db := db.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	return &AccountRepositoryImpl{db: db}, nil
}

func (r *AccountRepositoryImpl) Create(account *model.Account) error {
	return r.db.Create(account).Error
}

func (r *AccountRepositoryImpl) GetByID(id uint64) (*model.Account, error) {
	var account model.Account
	if err := r.db.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepositoryImpl) Transfer(sourceAccountId, destinationAccountId uint64, amount *float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var sourceAccount model.Account
		if err := tx.First(&sourceAccount, sourceAccountId).Error; err != nil {
			return fmt.Errorf("failed to find source account: %w", err)
		}
		var destinationAccount model.Account
		if err := tx.First(&destinationAccount, destinationAccountId).Error; err != nil {
			return fmt.Errorf("failed to find destination account: %w", err)
		}
		if *sourceAccount.Balance < *amount {
			return fmt.Errorf("insufficient balance in source account")
		}
		// Update balances
		*sourceAccount.Balance -= *amount
		*destinationAccount.Balance += *amount
		// Save changes
		if err := tx.Model(&sourceAccount).UpdateColumn("balance", sourceAccount.Balance).Error; err != nil {
			return err
		}
		if err := tx.Model(destinationAccount).UpdateColumn("balance", destinationAccount.Balance).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Transaction{}).Create(&model.Transaction{
			SourceAccount:      sourceAccount.AccountId,
			DestinationAccount: destinationAccount.AccountId,
			Amount:             amount,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}
