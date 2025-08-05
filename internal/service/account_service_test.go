package service

import (
	"testing"
	"wangfeng/transaction-system/internal/model"
	"wangfeng/transaction-system/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountService(t *testing.T) {
	mockRepository := mocks.NewAccountRepository(t)
	accountService, e := NewAccountService(mockRepository)

	assert.Nil(t, e)
	assert.NotNil(t, accountService)
}

func TestAccountSerivceCreateAccount(t *testing.T) {
	mockAccount := &model.Account{
		AccountId: ptrUint64(1),
		Balance:   ptrFloat64(100.0),
	}
	mockRepository := mocks.NewAccountRepository(t)
	mockRepository.On("Create", mockAccount).Return(nil)
	accountService, e := NewAccountService(mockRepository)

	assert.Nil(t, e)
	assert.NotNil(t, accountService)

	assert.Nil(t, accountService.CreateAccount(mockAccount))
}

func TestAccountSerivceGetAccount(t *testing.T) {
	var mockAccountId uint64 = 1
	mockAccount := &model.Account{
		AccountId: &mockAccountId,
		Balance:   ptrFloat64(100.0),
	}
	mockRepository := mocks.NewAccountRepository(t)
	mockRepository.On("GetByID", mockAccountId).Return(mockAccount, nil)
	accountService, e := NewAccountService(mockRepository)

	assert.Nil(t, e)
	assert.NotNil(t, accountService)

	response, e := accountService.GetAccount(mockAccountId)
	assert.Nil(t, e)
	assert.Equal(t, *response.AccountID, *mockAccount.AccountId)
	assert.Equal(t, *response.Balance, *mockAccount.Balance)
}

func TestAccountSerivceTransfer(t *testing.T) {
	var sourceAccountId uint64 = 1
	var destinationAccountId uint64 = 2
	var transferAmount float64 = 100.0

	mockRepository := mocks.NewAccountRepository(t)
	mockRepository.On("Transfer", sourceAccountId, destinationAccountId, &transferAmount).Return(nil)
	accountService, e := NewAccountService(mockRepository)

	assert.Nil(t, e)
	assert.NotNil(t, accountService)

	assert.Nil(t, accountService.Transfer(sourceAccountId, destinationAccountId, &transferAmount))
}

func ptrUint64(v uint64) *uint64 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}
