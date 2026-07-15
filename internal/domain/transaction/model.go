package transaction

type Transaction struct {
	From       string
	To         string
	Amount     float32
}

type MemPool interface {
	CreateTransaction(from string, to string, amount float32) Transaction
	PushTransaction(t Transaction)
	DeleteFirstTransactions(n int)
	PullFirstTransactions(n int) []Transaction
	Count() int
	Clear()
}
