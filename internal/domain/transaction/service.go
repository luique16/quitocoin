package transaction

type mempool struct {
	transactions []Transaction
}

func NewMemPool() MemPool {
	return &mempool{
		transactions: []Transaction{},
	}
}

func (m *mempool) CreateTransaction(from string, to string, amount float32) Transaction {
	return Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

func (m *mempool) PushTransaction(t Transaction) {
	m.transactions = append(m.transactions, t)
}

func (m *mempool) DeleteFirstTransactions(n int) {
	if n <= 0 {
		return
	}
	count := len(m.transactions)
	if n > count {
		n = count
	}
	m.transactions = m.transactions[n:]
}

func (m *mempool) PullFirstTransactions(n int) []Transaction {
	if n <= 0 {
		return []Transaction{}
	}
	count := len(m.transactions)
	if n > count {
		n = count
	}
	result := make([]Transaction, n)
	copy(result, m.transactions[:n])
	return result
}

func (m *mempool) Count() int {
	return len(m.transactions)
}

func (m *mempool) Clear() {
	m.transactions = []Transaction{}
}
