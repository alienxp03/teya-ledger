package storage

type Storage interface {
	CreateAccount(account Account) (*Account, error)
	GetAccount(userID string, accountNumber string) (*Account, error)

	CreateTransaction(transaction *Transaction) error
	GetTransactions(userID, accountNumber string, limit, page int) ([]*Transaction, error)

	CreateDeposit(transaction *Transaction) (*Transaction, error)
}

type MemoryStorage struct {
	accounts     []*Account
	transactions []*Transaction
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		transactions: []*Transaction{},
	}
}
