package transaction

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/alienxp03/teya-ledger/storage"
)

func TestCreateDeposit(t *testing.T) {
	type args struct {
		userID string
	}
	type setup struct {
		mockStorage *MockStorage
	}

	tests := []struct {
		name    string
		args    args
		setup   setup
		req     CreateDepositRequest
		want    *CreateDepositResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{userID: "USER_ID_1"},
			req: CreateDepositRequest{
				TransactionID: "idempotency-key",
				Amount:        100,
				Currency:      "MYR",
				Description:   "description",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetAccountFunc: func(userID, accountNumber string) (*storage.Account, error) {
						return &storage.Account{Number: "ACCOUNT_NUMBER_1"}, nil
					},
					CreateDepositFunc: func(transaction *storage.Transaction) (*storage.Transaction, error) {
						return &storage.Transaction{
							TransactionID: "idempotency-key",
							Status:        "pending",
							Amount:        100,
							Currency:      "MYR",
							Description:   "description",
							CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						}, nil
					},
					UpdateBalanceFunc: func(userID string, accountNumber string, amount int64) error {
						return nil
					},
					UpdateTransactionFunc: func(transactionID string, status string) error {
						return nil
					},
				}
				return setup{mockStorage}
			}(),
			want: &CreateDepositResponse{
				Transaction: Transaction{
					TransactionID: "idempotency-key",
					Status:        "pending",
					Amount:        100,
					Currency:      "MYR",
					Description:   "description",
					CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "invalid user account",
			req: CreateDepositRequest{
				TransactionID: "idempotency-key",
				Amount:        100,
				Currency:      "MYR",
				Description:   "description",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetAccountFunc: func(userID string, accountNumber string) (*storage.Account, error) {
						return nil, storage.ErrNotFound
					},
				}
				return setup{mockStorage}
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "error saving data",
			req: CreateDepositRequest{
				TransactionID: "idempotency-key",
				Amount:        100,
				Currency:      "MYR",
				Description:   "description",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetAccountFunc: func(userID string, accountNumber string) (*storage.Account, error) {
						return &storage.Account{Number: "account-number"}, nil
					},
					CreateDepositFunc: func(transaction *storage.Transaction) (*storage.Transaction, error) {
						return nil, errors.New("saving error")
					},
				}
				return setup{mockStorage}
			}(),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.setup.mockStorage)
			got, err := handler.CreateDeposit(tt.args.userID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDeposit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDeposit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	type args struct {
		userID string
		req    GetTransactionsRequest
	}
	type setup struct {
		mockStorage *MockStorage
	}

	tests := []struct {
		name    string
		args    args
		setup   setup
		want    *GetTransactionsResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userID: "USER_ID_1",
				req: GetTransactionsRequest{
					AccountNumber: "account-number",
					Limit:         10,
					Page:          1,
				},
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetTransactionsFunc: func(userID, accountNumber string, limit, page int) ([]*storage.Transaction, error) {
						return []*storage.Transaction{
							{
								TransactionID: "idempotency-key",
								Status:        "pending",
								Amount:        100,
								Currency:      "MYR",
								Description:   "description",
								CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					},
				}
				return setup{mockStorage}
			}(),
			want: &GetTransactionsResponse{
				Transactions: []Transaction{
					{
						TransactionID: "idempotency-key",
						Status:        "pending",
						Amount:        100,
						Currency:      "MYR",
						Description:   "description",
						CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			name: "logic error",
			args: args{
				userID: "USER_ID_1",
				req: GetTransactionsRequest{
					AccountNumber: "account-number",
					Limit:         10,
					Page:          1,
				},
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetTransactionsFunc: func(userID, accountNumber string, limit, page int) ([]*storage.Transaction, error) {
						return nil, errors.New("error")
					},
				}
				return setup{mockStorage}
			}(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.setup.mockStorage)

			got, err := handler.GetTransactions(tt.args.userID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDeposit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDeposit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateWithdrawal(t *testing.T) {
	type args struct {
		userID string
	}
	type setup struct {
		mockStorage *MockStorage
	}

	tests := []struct {
		name    string
		args    args
		setup   setup
		req     CreateWithdrawalRequest
		want    *CreateWithdrawalResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{userID: "USER_ID_1"},
			req: CreateWithdrawalRequest{
				TransactionID: "idempotency-key",
				AccountNumber: "ACCOUNT_NUMBER_1",
				Amount:        100,
				Currency:      "MYR",
				Description:   "withdrawal description",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetAccountFunc: func(userID, accountNumber string) (*storage.Account, error) {
						return &storage.Account{Number: "ACCOUNT_NUMBER_1"}, nil
					},
					GetBalanceFunc: func(userID string, accountNumber string) (*storage.Balance, error) {
						return &storage.Balance{Amount: 1000, Currency: "MYR"}, nil
					},
					CreateWithdrawalFunc: func(transaction *storage.Transaction) (*storage.Transaction, error) {
						return &storage.Transaction{
							TransactionID: "idempotency-key",
							Status:        "pending",
							Amount:        -100,
							Currency:      "MYR",
							Description:   "withdrawal description",
							CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						}, nil
					},
					UpdateBalanceFunc: func(userID string, accountNumber string, amount int64) error {
						return nil
					},
					UpdateTransactionFunc: func(transactionID string, status string) error {
						return nil
					},
				}
				return setup{mockStorage}
			}(),
			want: &CreateWithdrawalResponse{
				Transaction: Transaction{
					TransactionID: "idempotency-key",
					Status:        "pending",
					Amount:        -100,
					Currency:      "MYR",
					Description:   "withdrawal description",
					CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "invalid user account",
			req: CreateWithdrawalRequest{
				TransactionID: "idempotency-key",
				Amount:        100,
				Currency:      "MYR",
				Description:   "withdrawal description",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetAccountFunc: func(userID string, accountNumber string) (*storage.Account, error) {
						return nil, storage.ErrNotFound
					},
				}
				return setup{mockStorage}
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "error saving data",
			req: CreateWithdrawalRequest{
				TransactionID: "idempotency-key",
				Amount:        100,
				Currency:      "MYR",
				Description:   "withdrawal description",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetAccountFunc: func(userID string, accountNumber string) (*storage.Account, error) {
						return &storage.Account{Number: "account-number"}, nil
					},
					CreateWithdrawalFunc: func(transaction *storage.Transaction) (*storage.Transaction, error) {
						return nil, errors.New("saving error")
					},
					GetBalanceFunc: func(userID string, accountNumber string) (*storage.Balance, error) {
						return &storage.Balance{Amount: 1000, Currency: "MYR"}, nil
					},
				}
				return setup{mockStorage}
			}(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.setup.mockStorage)
			got, err := handler.CreateWithdrawal(tt.args.userID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWithdrawal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateWithdrawal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTransaction(t *testing.T) {
	type args struct {
		userID        string
		transactionID string
	}
	type setup struct {
		mockStorage *MockStorage
	}

	tests := []struct {
		name    string
		args    args
		setup   setup
		want    *Transaction
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userID:        "USER_ID_1",
				transactionID: "TRANSACTION_ID_1",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetTransactionFunc: func(userID, transactionID string) (*storage.Transaction, error) {
						return &storage.Transaction{
							TransactionID: "TRANSACTION_ID_1",
							UserID:        "USER_ID_1",
							Status:        "pending",
							Amount:        100,
							Currency:      "MYR",
							Description:   "test transaction",
							CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						}, nil
					},
				}
				return setup{mockStorage}
			}(),
			want: &Transaction{
				TransactionID: "TRANSACTION_ID_1",
				Status:        "pending",
				Amount:        100,
				Currency:      "MYR",
				Description:   "test transaction",
				CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "transaction not found",
			args: args{
				userID:        "USER_ID_1",
				transactionID: "TRANSACTION_ID_2",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetTransactionFunc: func(userID, transactionID string) (*storage.Transaction, error) {
						return nil, storage.ErrNotFound
					},
				}
				return setup{mockStorage}
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "transaction belongs to different user",
			args: args{
				userID:        "USER_ID_1",
				transactionID: "TRANSACTION_ID_3",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					GetTransactionFunc: func(userID, transactionID string) (*storage.Transaction, error) {
						return nil, storage.ErrNotFound
					},
				}
				return setup{mockStorage}
			}(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.setup.mockStorage)
			got, err := handler.GetTransaction(tt.args.userID, tt.args.transactionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransactionStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTransactionStatus(t *testing.T) {
	type args struct {
		transactionID string
	}
	type setup struct {
		mockStorage *MockStorage
	}

	tests := []struct {
		name    string
		args    args
		setup   setup
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				transactionID: "TRANSACTION_ID_1",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					UpdateTransactionFunc: func(transactionID string, status string) error {
						return nil
					},
				}
				return setup{mockStorage}
			}(),
			wantErr: false,
		},
		{
			name: "error updating status",
			args: args{
				transactionID: "TRANSACTION_ID_2",
			},
			setup: func() setup {
				mockStorage := &MockStorage{
					UpdateTransactionFunc: func(transactionID string, status string) error {
						return errors.New("update error")
					},
				}
				return setup{mockStorage}
			}(),
			wantErr: false, // Should not return error since it's a background task
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.setup.mockStorage)
			handler.updateTransaction(tt.args.transactionID)
			// Wait for the goroutine to complete
			time.Sleep(300 * time.Millisecond)
		})
	}
}

type MockStorage struct {
	CreateAccountFunc     func(accountNumber storage.Account) (*storage.Account, error)
	GetAccountFunc        func(userID string, accountNumber string) (*storage.Account, error)
	CreateTransactionFunc func(transaction *storage.Transaction) error
	GetTransactionsFunc   func(userID, accountNumber string, limit, page int) ([]*storage.Transaction, error)
	CreateDepositFunc     func(transaction *storage.Transaction) (*storage.Transaction, error)
	CreateWithdrawalFunc  func(transaction *storage.Transaction) (*storage.Transaction, error)
	GetBalanceFunc        func(userID, accountNumber string) (*storage.Balance, error)
	UpdateBalanceFunc     func(userID, accountNumber string, amount int64) error
	GetTransactionFunc    func(useriD, transactionID string) (*storage.Transaction, error)
	UpdateTransactionFunc func(transactionID string, status string) error
}

func (m *MockStorage) CreateAccount(account storage.Account) (*storage.Account, error) {
	return m.CreateAccountFunc(account)
}

func (m *MockStorage) CreateDeposit(transaction *storage.Transaction) (*storage.Transaction, error) {
	return m.CreateDepositFunc(transaction)
}

func (m *MockStorage) CreateTransaction(transaction *storage.Transaction) error {
	return m.CreateTransactionFunc(transaction)
}

func (m *MockStorage) GetTransactions(userID, accountNumber string, limit, page int) ([]*storage.Transaction, error) {
	return m.GetTransactionsFunc(userID, accountNumber, limit, page)
}

func (m *MockStorage) GetAccount(userID string, accountNumber string) (*storage.Account, error) {
	return m.GetAccountFunc(userID, accountNumber)
}

func (m *MockStorage) CreateWithdrawal(transaction *storage.Transaction) (*storage.Transaction, error) {
	return m.CreateWithdrawalFunc(transaction)
}

func (m *MockStorage) GetBalance(userID string, accountNumber string) (*storage.Balance, error) {
	return m.GetBalanceFunc(userID, accountNumber)
}

func (m *MockStorage) UpdateBalance(userID string, accountNumber string, amount int64) error {
	return m.UpdateBalanceFunc(userID, accountNumber, amount)
}

func (m *MockStorage) GetTransaction(userID, transactionID string) (*storage.Transaction, error) {
	return m.GetTransactionFunc(userID, transactionID)
}

func (m *MockStorage) UpdateTransaction(transactionID string, status string) error {
	return m.UpdateTransactionFunc(transactionID, status)
}
