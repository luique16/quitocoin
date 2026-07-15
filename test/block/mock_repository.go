package block_test

import (
	"context"

	"github.com/luique16/quitocoin/ent"
)

type mockCreateArgs struct {
	ctx context.Context
	b   *ent.Block
}

type mockGetByHashArgs struct {
	ctx  context.Context
	hash string
}

type mockGetByIndexArgs struct {
	ctx   context.Context
	index int
}

type mockGetLastArgs struct {
	ctx context.Context
}

type mockGetLastNArgs struct {
	ctx context.Context
	n   int
}

type mockListArgs struct {
	ctx context.Context
}

type mockCountArgs struct {
	ctx context.Context
}

type MockRepository struct {
	CreateFn    func(ctx context.Context, b *ent.Block) (*ent.Block, error)
	GetByHashFn func(ctx context.Context, hash string) (*ent.Block, error)
	GetByIndexFn func(ctx context.Context, index int) (*ent.Block, error)
	GetLastFn   func(ctx context.Context) (*ent.Block, error)
	GetLastNFn  func(ctx context.Context, n int) ([]*ent.Block, error)
	ListFn      func(ctx context.Context) ([]*ent.Block, error)
	CountFn     func(ctx context.Context) (int, error)

	createCalls    []mockCreateArgs
	getByHashCalls []mockGetByHashArgs
	getByIndexCalls []mockGetByIndexArgs
	getLastCalls   []mockGetLastArgs
	getLastNCalls  []mockGetLastNArgs
	listCalls      []mockListArgs
	countCalls     []mockCountArgs
}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) Create(ctx context.Context, b *ent.Block) (*ent.Block, error) {
	m.createCalls = append(m.createCalls, mockCreateArgs{ctx: ctx, b: b})
	if m.CreateFn != nil {
		return m.CreateFn(ctx, b)
	}
	return nil, nil
}

func (m *MockRepository) GetByHash(ctx context.Context, hash string) (*ent.Block, error) {
	m.getByHashCalls = append(m.getByHashCalls, mockGetByHashArgs{ctx: ctx, hash: hash})
	if m.GetByHashFn != nil {
		return m.GetByHashFn(ctx, hash)
	}
	return nil, nil
}

func (m *MockRepository) GetByIndex(ctx context.Context, index int) (*ent.Block, error) {
	m.getByIndexCalls = append(m.getByIndexCalls, mockGetByIndexArgs{ctx: ctx, index: index})
	if m.GetByIndexFn != nil {
		return m.GetByIndexFn(ctx, index)
	}
	return nil, nil
}

func (m *MockRepository) GetLast(ctx context.Context) (*ent.Block, error) {
	m.getLastCalls = append(m.getLastCalls, mockGetLastArgs{ctx: ctx})
	if m.GetLastFn != nil {
		return m.GetLastFn(ctx)
	}
	return nil, nil
}

func (m *MockRepository) GetLastN(ctx context.Context, n int) ([]*ent.Block, error) {
	m.getLastNCalls = append(m.getLastNCalls, mockGetLastNArgs{ctx: ctx, n: n})
	if m.GetLastNFn != nil {
		return m.GetLastNFn(ctx, n)
	}
	return nil, nil
}

func (m *MockRepository) List(ctx context.Context) ([]*ent.Block, error) {
	m.listCalls = append(m.listCalls, mockListArgs{ctx: ctx})
	if m.ListFn != nil {
		return m.ListFn(ctx)
	}
	return nil, nil
}

func (m *MockRepository) Count(ctx context.Context) (int, error) {
	m.countCalls = append(m.countCalls, mockCountArgs{ctx: ctx})
	if m.CountFn != nil {
		return m.CountFn(ctx)
	}
	return 0, nil
}

func (m *MockRepository) CreateCallCount() int      { return len(m.createCalls) }
func (m *MockRepository) GetByHashCallCount() int    { return len(m.getByHashCalls) }
func (m *MockRepository) GetByIndexCallCount() int   { return len(m.getByIndexCalls) }
func (m *MockRepository) GetLastCallCount() int      { return len(m.getLastCalls) }
func (m *MockRepository) GetLastNCallCount() int     { return len(m.getLastNCalls) }
func (m *MockRepository) ListCallCount() int         { return len(m.listCalls) }
func (m *MockRepository) CountCallCount() int        { return len(m.countCalls) }
