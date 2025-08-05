package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wangfeng/transaction-system/internal/model"
	"wangfeng/transaction-system/internal/service"
	"wangfeng/transaction-system/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAccountController(t *testing.T) {
	router := echo.New()
	mockService := mocks.NewAccountService(t)

	testCases := []struct {
		name           string
		router         *echo.Echo
		accountService service.AccountService
		wantError      bool
	}{
		{
			// OK
			name:           "OK",
			router:         router,
			accountService: mockService,
			wantError:      false,
		},
		{
			name:           "Param_Router_Required",
			router:         nil,
			accountService: mockService,
			wantError:      true,
		},
		{
			name:           "Param_AccountService_Required",
			router:         router,
			accountService: nil,
			wantError:      true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// execute
			c, e := NewAccountController(testCase.router, testCase.accountService)

			// log
			if e != nil {
				t.Log("error:", e)
			}

			// assert
			if testCase.wantError {
				assert.NotNil(t, e)
				assert.Nil(t, c)
			} else {
				assert.Nil(t, e)
				assert.NotNil(t, c)
			}
		})
	}

}

func TestAccountControllerCreateAccount(t *testing.T) {
	router := echo.New()
	mockService := mocks.NewAccountService(t)
	mockService.On("CreateAccount", mock.Anything).Return(nil)
	controller, _ := NewAccountController(router, mockService)

	testCases := []struct {
		name                   string
		accountId              *uint64
		initialBalance         *float64
		expectedHttpStatusCode int
	}{
		{
			name:                   "OK",
			accountId:              ptrUint64(1),
			initialBalance:         ptrFloat64(100.0),
			expectedHttpStatusCode: http.StatusOK,
		},
		{
			name:                   "Param_AccountID_Required",
			accountId:              nil,
			initialBalance:         ptrFloat64(100.0),
			expectedHttpStatusCode: http.StatusBadRequest,
		},
		{
			name:                   "Param_InitialBalance_Required",
			accountId:              ptrUint64(1),
			initialBalance:         nil,
			expectedHttpStatusCode: http.StatusBadRequest,
		},
		{
			name:                   "Param_InitialBalance_Invalid",
			accountId:              ptrUint64(1),
			initialBalance:         ptrFloat64(-100.0),
			expectedHttpStatusCode: http.StatusBadRequest,
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			requestBody := model.CreateAccountRequest{
				AccountID:      testCase.accountId,
				InitialBalance: testCase.initialBalance,
			}

			body, _ := json.Marshal(requestBody)
			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			ctx := router.NewContext(req, res)

			controller.create(ctx)
			assert.Equal(t, testCase.expectedHttpStatusCode, res.Code)
		})
	}
}

func TestAccountControllerGetAccount(t *testing.T) {
	router := echo.New()
	mockService := mocks.NewAccountService(t)
	mockService.On("GetAccount", mock.Anything).Return(&model.GetAccountResponse{}, nil)

	controller, _ := NewAccountController(router, mockService)

	testCases := []struct {
		name                   string
		accountId              string
		expectedHttpStatusCode int
	}{
		{
			name:                   "OK",
			accountId:              "1",
			expectedHttpStatusCode: http.StatusOK,
		},
		{
			name:                   "Param_Missing_AccountID",
			expectedHttpStatusCode: http.StatusBadRequest,
		},
		{
			name:                   "Param_Invalid_AccountID",
			accountId:              "invalid_account_id",
			expectedHttpStatusCode: http.StatusBadRequest,
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/accounts/"+testCase.accountId, nil)
			res := httptest.NewRecorder()
			ctx := router.NewContext(req, res)
			ctx.SetParamNames("id")
			ctx.SetParamValues(testCase.accountId)
			controller.get(ctx)
			assert.Equal(t, testCase.expectedHttpStatusCode, res.Code)
		})
	}
}

func TestAccountControllerTransfer(t *testing.T) {
	router := echo.New()
	mockService := mocks.NewAccountService(t)
	mockService.On("Transfer", uint64(1), uint64(2), ptrFloat64(100)).Return(nil)

	controller, _ := NewAccountController(router, mockService)
	testCases := []struct {
		name                   string
		sourceAccountId        *uint64
		destinationAccountId   *uint64
		amount                 *float64
		expectedHttpStatusCode int
	}{
		{
			name:                   "OK",
			sourceAccountId:        ptrUint64(1),
			destinationAccountId:   ptrUint64(2),
			amount:                 ptrFloat64(100),
			expectedHttpStatusCode: http.StatusOK,
		},
		{
			name:                   "Param_SourceAccountId_Missing",
			sourceAccountId:        nil,
			destinationAccountId:   ptrUint64(2),
			amount:                 ptrFloat64(100),
			expectedHttpStatusCode: http.StatusBadRequest,
		},
		{
			name:                   "Param_DestinationAccountId_Missing",
			sourceAccountId:        ptrUint64(1),
			destinationAccountId:   nil,
			amount:                 ptrFloat64(100),
			expectedHttpStatusCode: http.StatusBadRequest,
		},
		{
			name:                   "Param_SameAccounts_Transfer",
			sourceAccountId:        ptrUint64(1),
			destinationAccountId:   ptrUint64(1),
			amount:                 ptrFloat64(100),
			expectedHttpStatusCode: http.StatusBadRequest,
		},
		{
			name:                   "Param_Amount_Invalid",
			sourceAccountId:        ptrUint64(1),
			destinationAccountId:   ptrUint64(1),
			amount:                 ptrFloat64(-100),
			expectedHttpStatusCode: http.StatusBadRequest,
		},
		{
			name:                   "Param_Amount_Missing",
			sourceAccountId:        ptrUint64(1),
			destinationAccountId:   ptrUint64(1),
			amount:                 nil,
			expectedHttpStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			requestBody := model.TransferRequest{
				SourceAccountId:      testCase.sourceAccountId,
				DestinationAccountId: testCase.destinationAccountId,
				Amount:               testCase.amount,
			}

			body, _ := json.Marshal(requestBody)
			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			ctx := router.NewContext(req, res)

			controller.transfer(ctx)
			assert.Equal(t, testCase.expectedHttpStatusCode, res.Code)
		})
	}
}

func ptrUint64(v uint64) *uint64 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}
