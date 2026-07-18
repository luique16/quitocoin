package utxo_test

import (
	"context"
	"testing"
	"log"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	utxo "github.com/luique16/quitocoin/internal/domain/utxo"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

func newTestRepo() (utxo.Repository, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("miniredis: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return utxo.NewRepository(rdb), mr
}

func TestRealRepo_SetAndGetBalance(t *testing.T) {
	repo, mr := newTestRepo()
	defer mr.Close()
	ctx := context.Background()

	err := repo.SetBalance(ctx, "user-1", 100.0)
	assertNoError(t, err)

	balance, err := repo.GetBalance(ctx, "user-1")
	assertNoError(t, err)
	assertEqual(t, float32(100.0), balance)
}

func TestRealRepo_GetBalanceNotFound(t *testing.T) {
	repo, mr := newTestRepo()
	defer mr.Close()
	ctx := context.Background()

	balance, err := repo.GetBalance(ctx, "nonexistent")

	assertErrorIs(t, err, errorpkg.ErrUTXONotFound)
	assertEqual(t, float32(0), balance)
}

func TestRealRepo_UpdateBalance(t *testing.T) {
	repo, mr := newTestRepo()
	defer mr.Close()
	ctx := context.Background()

	_ = repo.SetBalance(ctx, "user-1", 100.0)

	err := repo.SetBalance(ctx, "user-1", 200.0)
	assertNoError(t, err)

	balance, err := repo.GetBalance(ctx, "user-1")
	assertNoError(t, err)
	assertEqual(t, float32(200.0), balance)
}

func TestRealRepo_MultipleUsers(t *testing.T) {
	repo, mr := newTestRepo()
	defer mr.Close()
	ctx := context.Background()

	_ = repo.SetBalance(ctx, "alice", 100.0)
	_ = repo.SetBalance(ctx, "bob", 50.0)

	aliceBal, _ := repo.GetBalance(ctx, "alice")
	bobBal, _ := repo.GetBalance(ctx, "bob")

	assertEqual(t, float32(100.0), aliceBal)
	assertEqual(t, float32(50.0), bobBal)
}

func TestRealRepo_ConcurrentAccess(t *testing.T) {
	repo, mr := newTestRepo()
	defer mr.Close()
	ctx := context.Background()

	_ = repo.SetBalance(ctx, "user-1", 0)

	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			_ = repo.SetBalance(ctx, "user-1", 50.0)
			_, _ = repo.GetBalance(ctx, "user-1")
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}
