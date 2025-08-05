package controller

import (
	"errors"
	"net/http"
	"strconv"
	"wangfeng/transaction-system/internal/model"
	"wangfeng/transaction-system/internal/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AccountController struct {
	accountService service.AccountService
}

func NewAccountController(router *echo.Echo, accountService service.AccountService) (*AccountController, error) {
	if router == nil {
		return nil, errors.New("routere can't be nil")
	}
	if accountService == nil {
		return nil, errors.New("accountService can't be nil")
	}

	c := &AccountController{accountService: accountService}

	router.POST("/accounts", c.create)
	router.GET("/accounts/:id", c.get)
	router.POST("/transactions", c.transfer)

	return c, nil
}

func (c *AccountController) create(ctx echo.Context) error {
	var request model.CreateAccountRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid account request")
	}
	if request.AccountID == nil || request.InitialBalance == nil {
		return ctx.JSON(http.StatusBadRequest, "Account ID and initial balance are required")
	}
	if *request.InitialBalance < 0 {
		return ctx.JSON(http.StatusBadRequest, "Initial balance must be greater than zero")
	}
	account := model.Account{
		AccountId: request.AccountID,
		Balance:   request.InitialBalance,
	}
	err := c.accountService.CreateAccount(&account)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to create account: "+err.Error())
	}
	return ctx.JSON(http.StatusOK, map[string]any{})
}

func (c *AccountController) get(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, "Invalid account ID")
	}

	accountID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid account ID format")
	}
	account, err := c.accountService.GetAccount(accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, "Account not found")
		}
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) transfer(ctx echo.Context) error {
	var request model.TransferRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid transaction request")
	}
	if request.SourceAccountId == nil || request.DestinationAccountId == nil || request.Amount == nil {
		return ctx.JSON(http.StatusBadRequest, "Source account, destination account, and amount are required")
	}
	if *request.Amount <= 0 {
		return ctx.JSON(http.StatusBadRequest, "Amount must be greater than zero")
	}
	if *request.SourceAccountId == *request.DestinationAccountId {
		return ctx.JSON(http.StatusBadRequest, "Source account and destination account must be different")
	}
	err := c.accountService.Transfer(*request.SourceAccountId, *request.DestinationAccountId, request.Amount)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to transfer: "+err.Error())
	}
	return ctx.JSON(200, "Transfer successful")
}
