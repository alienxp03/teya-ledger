package storage

type Storage interface {
	CreateTransaction(transaction *Transaction) error
	GetTransactions(accountNumber string, limit, page int) ([]*Transaction, error)
}

type MemoryStorage struct {
	transactions []*Transaction
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		transactions: []*Transaction{},
	}
}
