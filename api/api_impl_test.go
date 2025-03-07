package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alienxp03/teya-ledger/handler/transaction"
	"github.com/alienxp03/teya-ledger/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactions(t *testing.T) {
	type args struct {
		userToken string
	}
	type setup struct {
		mockTransactioner *MockTransactioner
	}

	tests := []struct {
		name    string
		args    args
		reqBody map[string]interface{}
		setup   setup
		want    GetTransactionsResponse
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"page": 1, "limit": 10},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetTransactionsFunc: func(userID string, req transaction.GetTransactionsRequest) (*transaction.GetTransactionsResponse, error) {
						return &transaction.GetTransactionsResponse{Transactions: []transaction.Transaction{{Amount: 100}}}, nil
					},
				}
				return setup{mockTransactioner}
			}(),
			want: GetTransactionsResponse{
				Transactions: []Transaction{{Amount: 100, CreatedAt: "0001-01-01T00:00:00Z", UpdatedAt: "0001-01-01T00:00:00Z"}},
			},
		},
		{
			name:    "logic error",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"page": 1, "limit": 10},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetTransactionsFunc: func(userID string, req transaction.GetTransactionsRequest) (*transaction.GetTransactionsResponse, error) {
						return nil, errors.New("logic error")
					},
				}
				return setup{mockTransactioner}
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New(tt.setup.mockTransactioner)

			reqBodyBytes, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("GET", "/api/v1/transactions", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.args.userToken)
			rr := httptest.NewRecorder()
			api.ServeHTTP(rr, req)

			if tt.wantErr {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				return
			}
			assert.Equal(t, http.StatusOK, rr.Code)
			var resp GetTransactionsResponse
			json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestCreateDeposit(t *testing.T) {
	type args struct {
		userToken string
	}
	type setup struct {
		mockTransactioner *MockTransactioner
	}

	tests := []struct {
		name    string
		args    args
		reqBody map[string]interface{}
		setup   setup
		want    CreateDepositResponse
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"transactionID": "idempotency-key", "accountNumber": "ACCOUNT_NUMBER_1", "amount": 100, "currency": "MYR", "description": "description"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					CreateDepositFunc: func(userID string, req transaction.CreateDepositRequest) (*transaction.CreateDepositResponse, error) {
						return &transaction.CreateDepositResponse{Transaction: transaction.Transaction{
							TransactionID: "idempotency-key",
							Status:        "pending",
							Amount:        100,
							Currency:      "MYR",
							Description:   "description",
						}}, nil
					},
				}
				return setup{mockTransactioner}
			}(),
			want: CreateDepositResponse{
				Transaction: Transaction{
					TransactionID: "idempotency-key",
					Status:        "pending",
					Amount:        100,
					Currency:      "MYR",
					Description:   "description",
					CreatedAt:     "0001-01-01T00:00:00Z",
					UpdatedAt:     "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name:    "negative amount",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"transactionID": "idempotency-key", "accountNumber": "ACCOUNT_NUMBER_1", "amount": -100, "currency": "MYR", "description": "description"},
			setup: func() setup {
				return setup{&MockTransactioner{}}
			}(),
			wantErr: true,
		},
		{
			name:    "invalid currency",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"transactionID": "idempotency-key", "accountNumber": "ACCOUNT_NUMBER_1", "amount": 100, "currency": "USD", "description": "description"},
			setup: func() setup {
				return setup{&MockTransactioner{}}
			}(),
			wantErr: true,
		},
		{
			name:    "logic error",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"transactionID": "idempotency-key", "accountNumber": "ACCOUNT_NUMBER_1", "amount": 100, "currency": "MYR", "description": "description"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					CreateDepositFunc: func(userID string, req transaction.CreateDepositRequest) (*transaction.CreateDepositResponse, error) {
						return nil, errors.New("logic error")
					},
				}
				return setup{mockTransactioner}
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New(tt.setup.mockTransactioner)

			reqBodyBytes, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/api/v1/deposits", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.args.userToken)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)

			if tt.wantErr {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				return
			}
			assert.Equal(t, http.StatusOK, w.Code)
			var resp CreateDepositResponse
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestCreateWithdrawal(t *testing.T) {
	type args struct {
		userToken string
	}
	type setup struct {
		mockTransactioner *MockTransactioner
	}

	tests := []struct {
		name    string
		args    args
		reqBody map[string]interface{}
		setup   setup
		want    CreateWithdrawalResponse
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"transactionID": "idempotency-key", "accountNumber": "ACCOUNT_NUMBER_1", "amount": -100, "currency": "MYR", "description": "withdrawal description"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					CreateWithdrawalFunc: func(userID string, req transaction.CreateWithdrawalRequest) (*transaction.CreateWithdrawalResponse, error) {
						return &transaction.CreateWithdrawalResponse{Transaction: transaction.Transaction{
							TransactionID: "idempotency-key",
							Status:        "pending",
							Amount:        -100,
							Currency:      "MYR",
							Description:   "withdrawal description",
						}}, nil
					},
				}
				return setup{mockTransactioner}
			}(),
			want: CreateWithdrawalResponse{
				Transaction: Transaction{
					TransactionID: "idempotency-key",
					Status:        "pending",
					Amount:        -100,
					Currency:      "MYR",
					Description:   "withdrawal description",
					CreatedAt:     "0001-01-01T00:00:00Z",
					UpdatedAt:     "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name:    "positive amount",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"transactionID": "idempotency-key", "accountNumber": "ACCOUNT_NUMBER_1", "amount": 100, "currency": "MYR", "description": "withdrawal description"},
			setup: func() setup {
				return setup{&MockTransactioner{}}
			}(),
			wantErr: true,
		},
		{
			name:    "invalid currency",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"transactionID": "idempotency-key", "accountNumber": "ACCOUNT_NUMBER_1", "amount": 100, "currency": "USD", "description": "withdrawal description"},
			setup: func() setup {
				return setup{&MockTransactioner{}}
			}(),
			wantErr: true,
		},
		{
			name:    "logic error",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"transactionID": "idempotency-key", "accountNumber": "ACCOUNT_NUMBER_1", "amount": -100, "currency": "MYR", "description": "withdrawal description"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					CreateWithdrawalFunc: func(userID string, req transaction.CreateWithdrawalRequest) (*transaction.CreateWithdrawalResponse, error) {
						return nil, errors.New("logic error")
					},
				}
				return setup{mockTransactioner}
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New(tt.setup.mockTransactioner)

			reqBodyBytes, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/api/v1/withdrawals", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.args.userToken)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)

			if tt.wantErr {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				return
			}
			assert.Equal(t, http.StatusOK, w.Code)
			var resp CreateWithdrawalResponse
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestGetBalance(t *testing.T) {
	type args struct {
		userToken string
	}
	type setup struct {
		mockTransactioner *MockTransactioner
	}

	tests := []struct {
		name    string
		args    args
		reqBody map[string]interface{}
		setup   setup
		want    GetBalanceResponse
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"accountNumber": "ACCOUNT_NUMBER_1"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetBalanceFunc: func(userID string, req transaction.GetBalanceRequest) (*transaction.GetBalanceResponse, error) {
						return &transaction.GetBalanceResponse{
							Amount:   1000,
							Currency: "MYR",
						}, nil
					},
				}
				return setup{mockTransactioner}
			}(),
			want: GetBalanceResponse{
				Balance: Balance{Amount: 1000, Currency: "MYR"},
			},
		},
		{
			name:    "invalid account",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"accountNumber": "INVALID_ACCOUNT"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetBalanceFunc: func(userID string, req transaction.GetBalanceRequest) (*transaction.GetBalanceResponse, error) {
						return nil, errors.New("account not found")
					},
				}
				return setup{mockTransactioner}
			}(),
			wantErr: true,
		},
		{
			name:    "logic error",
			args:    args{userToken: "USER_TOKEN_1"},
			reqBody: map[string]interface{}{"accountNumber": "ACCOUNT_NUMBER_1"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetBalanceFunc: func(userID string, req transaction.GetBalanceRequest) (*transaction.GetBalanceResponse, error) {
						return nil, errors.New("logic error")
					},
				}
				return setup{mockTransactioner}
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New(tt.setup.mockTransactioner)

			reqBodyBytes, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("GET", "/api/v1/balances", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.args.userToken)
			w := httptest.NewRecorder()
			api.ServeHTTP(w, req)

			if tt.wantErr {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				return
			}
			assert.Equal(t, http.StatusOK, w.Code)
			var resp GetBalanceResponse
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestGetTransactionStatus(t *testing.T) {
	type args struct {
		userToken string
	}
	type setup struct {
		mockTransactioner *MockTransactioner
	}

	tests := []struct {
		name    string
		args    args
		setup   setup
		want    GetTransactionResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{userToken: "USER_TOKEN_1"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetTransactionFunc: func(userID string, transactionID string) (*transaction.Transaction, error) {
						return &transaction.Transaction{
							TransactionID: "TRANSACTION_ID_1",
							Status:        "completed",
							Amount:        100,
							Currency:      "MYR",
							Description:   "test transaction",
						}, nil
					},
				}
				return setup{mockTransactioner}
			}(),
			want: GetTransactionResponse{
				Transaction{
					TransactionID: "TRANSACTION_ID_1",
					Status:        "completed",
					Amount:        100,
					Currency:      "MYR",
					Description:   "test transaction",
					CreatedAt:     "0001-01-01T00:00:00Z",
					UpdatedAt:     "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name: "transaction not found",
			args: args{userToken: "USER_TOKEN_1"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetTransactionFunc: func(userID string, transactionID string) (*transaction.Transaction, error) {
						return nil, storage.ErrNotFound
					},
				}
				return setup{mockTransactioner}
			}(),
			wantErr: true,
		},
		{
			name: "logic error",
			args: args{userToken: "USER_TOKEN_1"},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetTransactionFunc: func(userID string, transactionID string) (*transaction.Transaction, error) {
						return nil, errors.New("logic error")
					},
				}
				return setup{mockTransactioner}
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New(tt.setup.mockTransactioner)

			req, _ := http.NewRequest("GET", "/api/v1/transactions/TRANSACTION_ID_1", nil)
			req.Header.Set("Authorization", tt.args.userToken)
			r := httptest.NewRecorder()
			api.ServeHTTP(r, req)

			if tt.wantErr {
				assert.Equal(t, http.StatusBadRequest, r.Code)
				return
			}

			assert.Equal(t, http.StatusOK, r.Code)
			var resp GetTransactionResponse
			err := json.Unmarshal(r.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, resp)
		})
	}
}

// MockTransactioner is a mock implementation of the Transactioner interface
type MockTransactioner struct {
	GetTransactionsFunc  func(userID string, req transaction.GetTransactionsRequest) (*transaction.GetTransactionsResponse, error)
	CreateDepositFunc    func(userID string, req transaction.CreateDepositRequest) (*transaction.CreateDepositResponse, error)
	CreateWithdrawalFunc func(userID string, req transaction.CreateWithdrawalRequest) (*transaction.CreateWithdrawalResponse, error)
	GetBalanceFunc       func(userID string, req transaction.GetBalanceRequest) (*transaction.GetBalanceResponse, error)
	GetTransactionFunc   func(userID string, transactionID string) (*transaction.Transaction, error)
}

func (m *MockTransactioner) GetTransactions(userID string, req transaction.GetTransactionsRequest) (*transaction.GetTransactionsResponse, error) {
	return m.GetTransactionsFunc(userID, req)
}

func (m *MockTransactioner) CreateDeposit(userID string, req transaction.CreateDepositRequest) (*transaction.CreateDepositResponse, error) {
	return m.CreateDepositFunc(userID, req)
}

func (m *MockTransactioner) CreateWithdrawal(userID string, req transaction.CreateWithdrawalRequest) (*transaction.CreateWithdrawalResponse, error) {
	return m.CreateWithdrawalFunc(userID, req)
}

func (m *MockTransactioner) GetBalance(userID string, req transaction.GetBalanceRequest) (*transaction.GetBalanceResponse, error) {
	return m.GetBalanceFunc(userID, req)
}

func (m *MockTransactioner) GetTransaction(userID string, transactionID string) (*transaction.Transaction, error) {
	return m.GetTransactionFunc(userID, transactionID)
}
