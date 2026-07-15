package utxo_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	utxo "github.com/luique16/quitocoin/internal/domain/utxo"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

func newService(repo *MockRepository) utxo.Service {
	return utxo.NewService(repo)
}

// -- assertion helpers ---------------------------------------------------

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func assertErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error %q, got %q", target, err)
	}
}

func assertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	t.Helper()
	if v == nil {
		t.Fatal("expected non-nil, got nil")
	}
}

func assertTrue(t *testing.T, v bool, msg string) {
	t.Helper()
	if !v {
		t.Fatalf("expected true: %s", msg)
	}
}

func assertDeepEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

// -- SetBalance tests ----------------------------------------------------

func TestSetBalance_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.SetBalanceFn = func(_ context.Context, userId string, amount float32) error {
		assertEqual(t, "user-1", userId)
		assertEqual(t, float32(100.0), amount)
		return nil
	}

	err := svc.SetBalance(ctx, "user-1", 100.0)

	assertNoError(t, err)
	assertEqual(t, 1, repo.SetBalanceCallCount())
}

func TestSetBalance_ZeroAmount(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.SetBalanceFn = func(_ context.Context, userId string, amount float32) error {
		assertEqual(t, "user-1", userId)
		assertEqual(t, float32(0), amount)
		return nil
	}

	err := svc.SetBalance(ctx, "user-1", 0)

	assertNoError(t, err)
	assertEqual(t, 1, repo.SetBalanceCallCount())
}

func TestSetBalance_RepositoryError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.SetBalanceFn = func(_ context.Context, _ string, _ float32) error {
		return errorpkg.ErrInternal
	}

	err := svc.SetBalance(ctx, "user-1", 100.0)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, 1, repo.SetBalanceCallCount())
}

// -- GetBalance tests ----------------------------------------------------

func TestGetBalance_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, userId string) (float32, error) {
		assertEqual(t, "user-1", userId)
		return 250.0, nil
	}

	balance, err := svc.GetBalance(ctx, "user-1")

	assertNoError(t, err)
	assertEqual(t, float32(250.0), balance)
	assertEqual(t, 1, repo.GetBalanceCallCount())
}

func TestGetBalance_ZeroBalance(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 0, nil
	}

	balance, err := svc.GetBalance(ctx, "user-1")

	assertNoError(t, err)
	assertEqual(t, float32(0), balance)
	assertEqual(t, 1, repo.GetBalanceCallCount())
}

func TestGetBalance_RepositoryError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 0, errorpkg.ErrInternal
	}

	balance, err := svc.GetBalance(ctx, "user-1")

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, float32(0), balance)
	assertEqual(t, 1, repo.GetBalanceCallCount())
}

func TestGetBalance_EmptyUserId(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 0, nil
	}

	balance, err := svc.GetBalance(ctx, "")

	assertNoError(t, err)
	assertEqual(t, float32(0), balance)
	assertEqual(t, 1, repo.GetBalanceCallCount())
}

// -- Credit tests --------------------------------------------------------

func TestCredit_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, userId string) (float32, error) {
		assertEqual(t, "user-1", userId)
		return 100.0, nil
	}

	repo.SetBalanceFn = func(_ context.Context, userId string, amount float32) error {
		assertEqual(t, "user-1", userId)
		assertEqual(t, float32(150.0), amount)
		return nil
	}

	err := svc.Credit(ctx, "user-1", 50.0)

	assertNoError(t, err)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 1, repo.SetBalanceCallCount())
}

func TestCredit_FromZero(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 0, nil
	}

	repo.SetBalanceFn = func(_ context.Context, userId string, amount float32) error {
		assertEqual(t, "user-1", userId)
		assertEqual(t, float32(50.0), amount)
		return nil
	}

	err := svc.Credit(ctx, "user-1", 50.0)

	assertNoError(t, err)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 1, repo.SetBalanceCallCount())
}

func TestCredit_NegativeAmount(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	err := svc.Credit(ctx, "user-1", -50.0)

	assertErrorIs(t, err, errorpkg.ErrNegativeAmount)
	assertEqual(t, 0, repo.GetBalanceCallCount())
	assertEqual(t, 0, repo.SetBalanceCallCount())
}

func TestCredit_ZeroAmount(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	err := svc.Credit(ctx, "user-1", 0)

	assertNoError(t, err)
	assertEqual(t, 0, repo.GetBalanceCallCount())
	assertEqual(t, 0, repo.SetBalanceCallCount())
}

func TestCredit_GetBalanceError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 0, errorpkg.ErrInternal
	}

	err := svc.Credit(ctx, "user-1", 50.0)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 0, repo.SetBalanceCallCount())
}

func TestCredit_NotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 0, errorpkg.ErrUTXONotFound
	}

	repo.SetBalanceFn = func(_ context.Context, userId string, amount float32) error {
		assertEqual(t, "user-1", userId)
		assertEqual(t, float32(50.0), amount)
		return nil
	}

	err := svc.Credit(ctx, "user-1", 50.0)

	assertNoError(t, err)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 1, repo.SetBalanceCallCount())
}

func TestCredit_SetBalanceError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 100.0, nil
	}

	repo.SetBalanceFn = func(_ context.Context, _ string, _ float32) error {
		return errorpkg.ErrInternal
	}

	err := svc.Credit(ctx, "user-1", 50.0)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 1, repo.SetBalanceCallCount())
}

// -- Debit tests ---------------------------------------------------------

func TestDebit_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, userId string) (float32, error) {
		assertEqual(t, "user-1", userId)
		return 100.0, nil
	}

	repo.SetBalanceFn = func(_ context.Context, userId string, amount float32) error {
		assertEqual(t, "user-1", userId)
		assertEqual(t, float32(50.0), amount)
		return nil
	}

	err := svc.Debit(ctx, "user-1", 50.0)

	assertNoError(t, err)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 1, repo.SetBalanceCallCount())
}

func TestDebit_ExactBalance(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 100.0, nil
	}

	repo.SetBalanceFn = func(_ context.Context, _ string, amount float32) error {
		assertEqual(t, float32(0), amount)
		return nil
	}

	err := svc.Debit(ctx, "user-1", 100.0)

	assertNoError(t, err)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 1, repo.SetBalanceCallCount())
}

func TestDebit_InsufficientBalance(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 30.0, nil
	}

	err := svc.Debit(ctx, "user-1", 50.0)

	assertErrorIs(t, err, errorpkg.ErrInsufficientBalance)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 0, repo.SetBalanceCallCount())
}

func TestDebit_NegativeAmount(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	err := svc.Debit(ctx, "user-1", -50.0)

	assertErrorIs(t, err, errorpkg.ErrNegativeAmount)
	assertEqual(t, 0, repo.GetBalanceCallCount())
	assertEqual(t, 0, repo.SetBalanceCallCount())
}

func TestDebit_ZeroAmount(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	err := svc.Debit(ctx, "user-1", 0)

	assertNoError(t, err)
	assertEqual(t, 0, repo.GetBalanceCallCount())
	assertEqual(t, 0, repo.SetBalanceCallCount())
}

func TestDebit_BalanceNotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 0, errorpkg.ErrUTXONotFound
	}

	err := svc.Debit(ctx, "user-1", 50.0)

	assertErrorIs(t, err, errorpkg.ErrUTXONotFound)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 0, repo.SetBalanceCallCount())
}

func TestDebit_GetBalanceError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 0, errorpkg.ErrInternal
	}

	err := svc.Debit(ctx, "user-1", 50.0)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 0, repo.SetBalanceCallCount())
}

func TestDebit_SetBalanceError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetBalanceFn = func(_ context.Context, _ string) (float32, error) {
		return 100.0, nil
	}

	repo.SetBalanceFn = func(_ context.Context, _ string, _ float32) error {
		return errorpkg.ErrInternal
	}

	err := svc.Debit(ctx, "user-1", 50.0)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, 1, repo.GetBalanceCallCount())
	assertEqual(t, 1, repo.SetBalanceCallCount())
}
