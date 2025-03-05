package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alienxp03/teya-ledger/handler/transaction"
	"github.com/stretchr/testify/assert"
)

func TestTransactionsList(t *testing.T) {
	type setup struct {
		mockTransactioner *MockTransactioner
	}

	tests := []struct {
		name    string
		reqBody map[string]interface{}
		setup   setup
		want    GetTransactionsResponse
		wantErr bool
	}{
		{
			name:    "success",
			reqBody: map[string]interface{}{"page": 1, "limit": 10},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetTransactionsFunc: func(w http.ResponseWriter, req transaction.GetTransactionsRequest) (transaction.GetTransactionsResponse, error) {
						return transaction.GetTransactionsResponse{Transactions: []transaction.Transaction{{Amount: 100}}}, nil
					},
				}
				return setup{mockTransactioner}
			}(),
			want: GetTransactionsResponse{
				Transactions: []Transaction{{Amount: 100, CreatedAt: "0001-01-01T00:00:00Z", UpdatedAt: "0001-01-01T00:00:00Z"}},
			},
		},
		{
			name:    "invalid page",
			reqBody: map[string]interface{}{"page": 1, "limit": 10},
			setup: func() setup {
				mockTransactioner := &MockTransactioner{
					GetTransactionsFunc: func(w http.ResponseWriter, req transaction.GetTransactionsRequest) (transaction.GetTransactionsResponse, error) {
						return transaction.GetTransactionsResponse{}, errors.New("logic error")
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

// MockTransactioner is a mock implementation of the Transactioner interface
type MockTransactioner struct {
	GetTransactionsFunc func(w http.ResponseWriter, req transaction.GetTransactionsRequest) (transaction.GetTransactionsResponse, error)
}

func (m *MockTransactioner) GetTransactions(w http.ResponseWriter, req transaction.GetTransactionsRequest) (transaction.GetTransactionsResponse, error) {
	return m.GetTransactionsFunc(w, req)
}
