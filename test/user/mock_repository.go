package user_test

import (
	"context"

	"github.com/luique16/quitocoin/ent"
)

type mockCreateArgs struct {
	ctx context.Context
	u   *ent.User
}

type mockGetArgs struct {
	ctx context.Context
	id  string
}

type mockGetByEmailArgs struct {
	ctx   context.Context
	email string
}

type mockGetByPublicIDArgs struct {
	ctx      context.Context
	publicID string
}

type mockListArgs struct {
	ctx context.Context
}

type MockRepository struct {
	CreateFn  func(ctx context.Context, u *ent.User) (*ent.User, error)
	GetFn     func(ctx context.Context, id string) (*ent.User, error)
	GetByEmailFn func(ctx context.Context, email string) (*ent.User, error)
	GetByPublicIDFn func(ctx context.Context, publicID string) (*ent.User, error)
	ListFn    func(ctx context.Context) ([]*ent.User, error)
	UpdateFn  func(ctx context.Context, u *ent.User) (*ent.User, error)
	DeleteFn  func(ctx context.Context, id string) error

	createCalls  []mockCreateArgs
	getCalls     []mockGetArgs
	getByEmailCalls []mockGetByEmailArgs
	getByPublicIDCalls []mockGetByPublicIDArgs
	listCalls    []mockListArgs
	updateCalls  []mockUpdateArgs
	deleteCalls  []mockDeleteArgs
}

type mockUpdateArgs struct {
	ctx context.Context
	u   *ent.User
}

type mockDeleteArgs struct {
	ctx context.Context
	id  string
}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	m.createCalls = append(m.createCalls, mockCreateArgs{ctx: ctx, u: u})
	if m.CreateFn != nil {
		return m.CreateFn(ctx, u)
	}
	return nil, nil
}

func (m *MockRepository) Get(ctx context.Context, id string) (*ent.User, error) {
	m.getCalls = append(m.getCalls, mockGetArgs{ctx: ctx, id: id})
	if m.GetFn != nil {
		return m.GetFn(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	m.getByEmailCalls = append(m.getByEmailCalls, mockGetByEmailArgs{ctx: ctx, email: email})
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *MockRepository) GetByPublicID(ctx context.Context, publicID string) (*ent.User, error) {
	m.getByPublicIDCalls = append(m.getByPublicIDCalls, mockGetByPublicIDArgs{ctx: ctx, publicID: publicID})
	if m.GetByPublicIDFn != nil {
		return m.GetByPublicIDFn(ctx, publicID)
	}
	return nil, nil
}

func (m *MockRepository) List(ctx context.Context) ([]*ent.User, error) {
	m.listCalls = append(m.listCalls, mockListArgs{ctx: ctx})
	if m.ListFn != nil {
		return m.ListFn(ctx)
	}
	return nil, nil
}

func (m *MockRepository) Update(ctx context.Context, u *ent.User) (*ent.User, error) {
	m.updateCalls = append(m.updateCalls, mockUpdateArgs{ctx: ctx, u: u})
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, u)
	}
	return nil, nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	m.deleteCalls = append(m.deleteCalls, mockDeleteArgs{ctx: ctx, id: id})
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockRepository) CreateCallCount() int  { return len(m.createCalls) }
func (m *MockRepository) GetCallCount() int     { return len(m.getCalls) }
func (m *MockRepository) GetByEmailCallCount() int { return len(m.getByEmailCalls) }
func (m *MockRepository) GetByPublicIDCallCount() int { return len(m.getByPublicIDCalls) }
func (m *MockRepository) ListCallCount() int    { return len(m.listCalls) }
func (m *MockRepository) UpdateCallCount() int  { return len(m.updateCalls) }
func (m *MockRepository) DeleteCallCount() int  { return len(m.deleteCalls) }

func (m *MockRepository) CreateArgs(index int) (context.Context, *ent.User) {
	if index >= len(m.createCalls) {
		return nil, nil
	}
	return m.createCalls[index].ctx, m.createCalls[index].u
}

func (m *MockRepository) GetArgs(index int) (context.Context, string) {
	if index >= len(m.getCalls) {
		return nil, ""
	}
	return m.getCalls[index].ctx, m.getCalls[index].id
}

func (m *MockRepository) GetByEmailArgs(index int) (context.Context, string) {
	if index >= len(m.getByEmailCalls) {
		return nil, ""
	}
	return m.getByEmailCalls[index].ctx, m.getByEmailCalls[index].email
}

func (m *MockRepository) GetByPublicIDArgs(index int) (context.Context, string) {
	if index >= len(m.getByPublicIDCalls) {
		return nil, ""
	}
	return m.getByPublicIDCalls[index].ctx, m.getByPublicIDCalls[index].publicID
}

func (m *MockRepository) UpdateArgs(index int) (context.Context, *ent.User) {
	if index >= len(m.updateCalls) {
		return nil, nil
	}
	return m.updateCalls[index].ctx, m.updateCalls[index].u
}

func (m *MockRepository) DeleteArgs(index int) (context.Context, string) {
	if index >= len(m.deleteCalls) {
		return nil, ""
	}
	return m.deleteCalls[index].ctx, m.deleteCalls[index].id
}
