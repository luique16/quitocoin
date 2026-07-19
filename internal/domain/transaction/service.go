package transaction

import "context"

type service struct {
	repo Repository
	ctx  context.Context
}

func NewService(repo Repository) MemPool {
	return &service{repo: repo, ctx: context.Background()}
}

type inMemoryRepo struct {
	txs []Transaction
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepo{txs: []Transaction{}}
}

func (r *inMemoryRepo) Push(_ context.Context, t Transaction) error {
	r.txs = append(r.txs, t)
	return nil
}

func (r *inMemoryRepo) PullFirst(_ context.Context, n int) ([]Transaction, error) {
	if n <= 0 {
		return []Transaction{}, nil
	}
	count := len(r.txs)
	if n > count {
		n = count
	}
	result := make([]Transaction, n)
	copy(result, r.txs[:n])
	return result, nil
}

func (r *inMemoryRepo) DeleteFirst(_ context.Context, n int) error {
	if n <= 0 {
		return nil
	}
	count := len(r.txs)
	if n > count {
		n = count
	}
	r.txs = r.txs[n:]
	return nil
}

func (r *inMemoryRepo) Count(_ context.Context) (int, error) {
	return len(r.txs), nil
}

func (r *inMemoryRepo) Clear(_ context.Context) error {
	r.txs = []Transaction{}
	return nil
}

func (s *service) CreateTransaction(from string, to string, amount float32) Transaction {
	return Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

func (s *service) PushTransaction(t Transaction) {
	s.repo.Push(s.ctx, t)
}

func (s *service) DeleteFirstTransactions(n int) {
	s.repo.DeleteFirst(s.ctx, n)
}

func (s *service) PullFirstTransactions(n int) []Transaction {
	txs, err := s.repo.PullFirst(s.ctx, n)
	if err != nil {
		return []Transaction{}
	}
	return txs
}

func (s *service) Count() int {
	count, err := s.repo.Count(s.ctx)
	if err != nil {
		return 0
	}
	return count
}

func (s *service) Clear() {
	s.repo.Clear(s.ctx)
}
