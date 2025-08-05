package service

import (
	"wangfeng/transaction-system/internal/model"
	"wangfeng/transaction-system/internal/repository"
)

type AccountService interface {
	CreateAccount(account *model.Account) error
	GetAccount(id uint64) (*model.GetAccountResponse, error)
	Transfer(sourceAccountId, destinationAccountId uint64, amount *float64) error
}

type AccountServiceImpl struct {
	accountRepository repository.AccountRepository
}

var _ AccountService = (*AccountServiceImpl)(nil)

func NewAccountService(accountRepository repository.AccountRepository) (*AccountServiceImpl, error) {
	return &AccountServiceImpl{
		accountRepository: accountRepository,
	}, nil
}

func (s *AccountServiceImpl) CreateAccount(account *model.Account) error {
	return s.accountRepository.Create(account)
}

func (s *AccountServiceImpl) GetAccount(id uint64) (*model.GetAccountResponse, error) {
	var response model.GetAccountResponse

	foundAccount, err := s.accountRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	response.AccountID = foundAccount.AccountId
	response.Balance = foundAccount.Balance

	return &response, nil
}

func (s *AccountServiceImpl) Transfer(sourceAccountId, destinationAccountId uint64, amount *float64) error {
	return s.accountRepository.Transfer(sourceAccountId, destinationAccountId, amount)
}
