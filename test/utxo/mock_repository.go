package utxo_test

import (
	"context"
)

type mockGetBalanceArgs struct {
	ctx    context.Context
	userId string
}

type mockSetBalanceArgs struct {
	ctx    context.Context
	userId string
	amount float32
}

type MockRepository struct {
	GetBalanceFn func(ctx context.Context, userId string) (float32, error)
	SetBalanceFn func(ctx context.Context, userId string, amount float32) error
	ClearFn      func(ctx context.Context) error

	getBalanceCalls []mockGetBalanceArgs
	setBalanceCalls []mockSetBalanceArgs
	clearCalls      []struct{ ctx context.Context }
}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) GetBalance(ctx context.Context, userId string) (float32, error) {
	m.getBalanceCalls = append(m.getBalanceCalls, mockGetBalanceArgs{ctx: ctx, userId: userId})
	if m.GetBalanceFn != nil {
		return m.GetBalanceFn(ctx, userId)
	}
	return 0, nil
}

func (m *MockRepository) SetBalance(ctx context.Context, userId string, amount float32) error {
	m.setBalanceCalls = append(m.setBalanceCalls, mockSetBalanceArgs{ctx: ctx, userId: userId, amount: amount})
	if m.SetBalanceFn != nil {
		return m.SetBalanceFn(ctx, userId, amount)
	}
	return nil
}

func (m *MockRepository) Clear(ctx context.Context) error {
	m.clearCalls = append(m.clearCalls, struct{ ctx context.Context }{ctx: ctx})
	if m.ClearFn != nil {
		return m.ClearFn(ctx)
	}
	return nil
}

func (m *MockRepository) GetBalanceCallCount() int { return len(m.getBalanceCalls) }
func (m *MockRepository) SetBalanceCallCount() int { return len(m.setBalanceCalls) }
func (m *MockRepository) ClearCallCount() int     { return len(m.clearCalls) }
