package block_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/luique16/quitocoin/ent"
	block "github.com/luique16/quitocoin/internal/domain/block"
	"github.com/luique16/quitocoin/internal/domain/transaction"
	errorpkg "github.com/luique16/quitocoin/internal/error"
)

func newService(repo *MockRepository) block.Service {
	return block.NewService(10, repo)
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

func assertNotEqual[T comparable](t *testing.T, a, b T) {
	t.Helper()
	if a == b {
		t.Fatalf("expected different values, both are %v", a)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	t.Helper()
	if v == nil {
		t.Fatal("expected non-nil, got nil")
	}
}

func assertNil(t *testing.T, v interface{}) {
	t.Helper()
	if v == nil {
		return
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		if rv.IsNil() {
			return
		}
	}
	t.Fatalf("expected nil, got %v", v)
}

func assertTrue(t *testing.T, v bool, msg string) {
	t.Helper()
	if !v {
		t.Fatalf("expected true: %s", msg)
	}
}

func assertLen(t *testing.T, items interface{}, expected int) {
	t.Helper()
	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Slice {
		t.Fatalf("expected slice, got %T", items)
	}
	if v.Len() != expected {
		t.Fatalf("expected length %d, got %d", expected, v.Len())
	}
}

func assertDeepEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func calculateHashHelper(index int, miner string, reward float32, previousHash string, nonce int64, tx []transaction.Transaction) string {
	transactionsData := ""
	for i := range tx {
		transactionsData += fmt.Sprintf("%s:%f:%s;", tx[i].From, tx[i].Amount, tx[i].To)
	}

	data := fmt.Sprintf("%d:%s:%f:%s:%d:%s", index, miner, reward, previousHash, nonce, transactionsData)

	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}

// -- test data helpers ---------------------------------------------------

func now() time.Time {
	return time.Now().Truncate(time.Second)
}

func makeBlock(index int, overrides ...func(*ent.Block)) *ent.Block {
	t := now()
	b := &ent.Block{
		ID:           index,
		Hash:         "",
		Index:        index,
		PreviousHash: "",
		Nonce:        0,
		Miner:        "system",
		Reward:       0,
		Transactions: nil,
		CreatedAt:    t,
	}

	hash := calculateHashHelper(index, "system", 0, "", 0, nil)
	b.Hash = hash

	for _, o := range overrides {
		o(b)
	}
	return b
}

func makeValidChain(n int) []*ent.Block {
	blocks := make([]*ent.Block, n)
	for i := 0; i < n; i++ {
		prevHash := ""
		if i > 0 {
			prevHash = blocks[i-1].Hash
		}
		t := now().Add(time.Duration(i) * time.Second)
		blocks[i] = &ent.Block{
			ID:           i,
			Hash:         "",
			Index:        i,
			PreviousHash: prevHash,
			Nonce:        0,
			Miner:        "miner",
			Reward:       1.0,
			Transactions: nil,
			CreatedAt:    t,
		}
		blocks[i].Hash = calculateHashHelper(i, "miner", 1.0, prevHash, 0, nil)
	}
	return blocks
}

// -- CreateGenesisBlock tests -------------------------------------------

func TestCreateGenesisBlock_Success(t *testing.T) {
	oldPrefix := block.DifficultyPrefix
	block.DifficultyPrefix = ""
	defer func() { block.DifficultyPrefix = oldPrefix }()

	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.CreateFn = func(_ context.Context, b *ent.Block) (*ent.Block, error) {
		assertEqual(t, 0, b.Index)
		assertEqual(t, "system", b.Miner)
		assertEqual(t, float64(0), b.Reward)
		assertEqual(t, "", b.PreviousHash)
		assertTrue(t, b.Nonce >= 0, "nonce must be non-negative")
		assertTrue(t, b.Hash != "", "hash must be set")
		assertTrue(t, !b.CreatedAt.IsZero(), "created_at must be set")
		return b, nil
	}

	result, err := svc.CreateGenesisBlock(ctx)

	assertNoError(t, err)
	assertNotNil(t, result)
	assertEqual(t, 0, result.Index)
	assertEqual(t, "system", result.Miner)
	assertEqual(t, "", result.PreviousHash)
	assertTrue(t, result.Nonce >= 0, "nonce must be non-negative")
	assertTrue(t, result.Hash != "", "hash must be returned")
	assertEqual(t, 1, repo.CreateCallCount())
}

func TestCreateGenesisBlock_RepositoryError(t *testing.T) {
	oldPrefix := block.DifficultyPrefix
	block.DifficultyPrefix = ""
	defer func() { block.DifficultyPrefix = oldPrefix }()

	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.CreateFn = func(_ context.Context, _ *ent.Block) (*ent.Block, error) {
		return nil, errorpkg.ErrInternal
	}

	result, err := svc.CreateGenesisBlock(ctx)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertNil(t, result)
	assertEqual(t, 1, repo.CreateCallCount())
}

func TestCreateGenesisBlock_BlockExists(t *testing.T) {
	oldPrefix := block.DifficultyPrefix
	block.DifficultyPrefix = ""
	defer func() { block.DifficultyPrefix = oldPrefix }()

	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.CreateFn = func(_ context.Context, _ *ent.Block) (*ent.Block, error) {
		return nil, errorpkg.ErrBlockExists
	}

	result, err := svc.CreateGenesisBlock(ctx)

	assertErrorIs(t, err, errorpkg.ErrBlockExists)
	assertNil(t, result)
	assertEqual(t, 1, repo.CreateCallCount())
}

// -- TryToMineBlock tests ------------------------------------------------

func TestTryToMineBlock_Success(t *testing.T) {
	oldPrefix := block.DifficultyPrefix
	block.DifficultyPrefix = ""
	defer func() { block.DifficultyPrefix = oldPrefix }()

	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	genesis := makeBlock(0)
	genesis.Hash = "000abc123"

	repo.GetLastFn = func(_ context.Context) (*ent.Block, error) {
		return genesis, nil
	}

	repo.CreateFn = func(_ context.Context, b *ent.Block) (*ent.Block, error) {
		assertEqual(t, 1, b.Index)
		assertEqual(t, "alice", b.Miner)
		assertEqual(t, genesis.Hash, b.PreviousHash)
		assertTrue(t, b.Hash != "", "hash must be set")
		assertTrue(t, b.Hash != genesis.Hash, "hash must be different from genesis")
		return b, nil
	}

	result, err := svc.TryToMineBlock(ctx, "alice", 0, 1.0)

	assertNoError(t, err)
	assertNotNil(t, result)
	assertEqual(t, 1, result.Index)
	assertEqual(t, "alice", result.Miner)
	assertEqual(t, genesis.Hash, result.PreviousHash)
	assertEqual(t, 1, repo.GetLastCallCount())
	assertEqual(t, 1, repo.CreateCallCount())
}

func TestTryToMineBlock_InvalidNonce(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	genesis := makeBlock(0)
	genesis.Hash = "ffffffff"

	repo.GetLastFn = func(_ context.Context) (*ent.Block, error) {
		return genesis, nil
	}

	result, err := svc.TryToMineBlock(ctx, "alice", 0, 1.0)

	assertErrorIs(t, err, errorpkg.ErrInvalidNonce)
	assertNil(t, result)
	assertEqual(t, 1, repo.GetLastCallCount())
	assertEqual(t, 0, repo.CreateCallCount())
}

func TestTryToMineBlock_NoLastBlock(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetLastFn = func(_ context.Context) (*ent.Block, error) {
		return nil, errorpkg.ErrBlockNotFound
	}

	result, err := svc.TryToMineBlock(ctx, "alice", 0, 1.0)

	assertErrorIs(t, err, errorpkg.ErrBlockNotFound)
	assertNil(t, result)
	assertEqual(t, 1, repo.GetLastCallCount())
	assertEqual(t, 0, repo.CreateCallCount())
}

func TestTryToMineBlock_CreateError(t *testing.T) {
	oldPrefix := block.DifficultyPrefix
	block.DifficultyPrefix = ""
	defer func() { block.DifficultyPrefix = oldPrefix }()

	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	genesis := makeBlock(0)
	genesis.Hash = "000abc123"

	repo.GetLastFn = func(_ context.Context) (*ent.Block, error) {
		return genesis, nil
	}

	repo.CreateFn = func(_ context.Context, _ *ent.Block) (*ent.Block, error) {
		return nil, errorpkg.ErrInternal
	}

	result, err := svc.TryToMineBlock(ctx, "alice", 0, 1.0)
	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertNil(t, result)
	assertEqual(t, 1, repo.GetLastCallCount())
	assertEqual(t, 1, repo.CreateCallCount())
}

// -- GetBlockByHash tests ------------------------------------------------

func TestGetBlockByHash_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	expected := makeBlock(0)

	repo.GetByHashFn = func(_ context.Context, hash string) (*ent.Block, error) {
		assertEqual(t, expected.Hash, hash)
		return expected, nil
	}

	result, err := svc.GetBlockByHash(ctx, expected.Hash)

	assertNoError(t, err)
	assertDeepEqual(t, expected, result)
	assertEqual(t, 1, repo.GetByHashCallCount())
}

func TestGetBlockByHash_EmptyHash(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	result, err := svc.GetBlockByHash(context.Background(), "")

	assertErrorIs(t, err, errorpkg.ErrInvalidBlockHash)
	assertNil(t, result)
	assertEqual(t, 0, repo.GetByHashCallCount())
}

func TestGetBlockByHash_NotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetByHashFn = func(_ context.Context, _ string) (*ent.Block, error) {
		return nil, errorpkg.ErrBlockNotFound
	}

	result, err := svc.GetBlockByHash(ctx, "nonexistent")

	assertErrorIs(t, err, errorpkg.ErrBlockNotFound)
	assertNil(t, result)
	assertEqual(t, 1, repo.GetByHashCallCount())
}

// -- GetBlockByIndex tests ----------------------------------------------

func TestGetBlockByIndex_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	expected := makeBlock(0)

	repo.GetByIndexFn = func(_ context.Context, index int) (*ent.Block, error) {
		assertEqual(t, 0, index)
		return expected, nil
	}

	result, err := svc.GetBlockByIndex(ctx, 0)

	assertNoError(t, err)
	assertDeepEqual(t, expected, result)
	assertEqual(t, 1, repo.GetByIndexCallCount())
}

func TestGetBlockByIndex_Negative(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	result, err := svc.GetBlockByIndex(context.Background(), -1)

	assertErrorIs(t, err, errorpkg.ErrInvalidBlockIndex)
	assertNil(t, result)
	assertEqual(t, 0, repo.GetByIndexCallCount())
}

func TestGetBlockByIndex_NotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetByIndexFn = func(_ context.Context, _ int) (*ent.Block, error) {
		return nil, errorpkg.ErrBlockNotFound
	}

	result, err := svc.GetBlockByIndex(ctx, 999)

	assertErrorIs(t, err, errorpkg.ErrBlockNotFound)
	assertNil(t, result)
	assertEqual(t, 1, repo.GetByIndexCallCount())
}

// -- GetLastBlock tests -------------------------------------------------

func TestGetLastBlock_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	expected := makeBlock(5)

	repo.GetLastFn = func(_ context.Context) (*ent.Block, error) {
		return expected, nil
	}

	result, err := svc.GetLastBlock(ctx)

	assertNoError(t, err)
	assertDeepEqual(t, expected, result)
	assertEqual(t, 1, repo.GetLastCallCount())
}

func TestGetLastBlock_NotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetLastFn = func(_ context.Context) (*ent.Block, error) {
		return nil, errorpkg.ErrBlockNotFound
	}

	result, err := svc.GetLastBlock(ctx)

	assertErrorIs(t, err, errorpkg.ErrBlockNotFound)
	assertNil(t, result)
	assertEqual(t, 1, repo.GetLastCallCount())
}

// -- GetLastBlocks tests ------------------------------------------------

func TestGetLastBlocks_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	expected := makeValidChain(3)

	repo.GetLastNFn = func(_ context.Context, n int) ([]*ent.Block, error) {
		assertEqual(t, 2, n)
		return expected[1:], nil
	}

	result, err := svc.GetLastBlocks(ctx, 2)

	assertNoError(t, err)
	assertLen(t, result, 2)
	assertEqual(t, 1, repo.GetLastNCallCount())
}

func TestGetLastBlocks_Zero(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	result, err := svc.GetLastBlocks(context.Background(), 0)

	assertNoError(t, err)
	assertLen(t, result, 0)
	assertEqual(t, 0, repo.GetLastNCallCount())
}

func TestGetLastBlocks_Negative(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	result, err := svc.GetLastBlocks(context.Background(), -1)

	assertNoError(t, err)
	assertLen(t, result, 0)
	assertEqual(t, 0, repo.GetLastNCallCount())
}

func TestGetLastBlocks_Empty(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetLastNFn = func(_ context.Context, n int) ([]*ent.Block, error) {
		return []*ent.Block{}, nil
	}

	result, err := svc.GetLastBlocks(ctx, 5)

	assertNoError(t, err)
	assertLen(t, result, 0)
	assertEqual(t, 1, repo.GetLastNCallCount())
}

// -- GetChainLength tests -----------------------------------------------

func TestGetChainLength_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.CountFn = func(_ context.Context) (int, error) {
		return 10, nil
	}

	length, err := svc.GetChainLength(ctx)

	assertNoError(t, err)
	assertEqual(t, 10, length)
	assertEqual(t, 1, repo.CountCallCount())
}

func TestGetChainLength_Zero(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.CountFn = func(_ context.Context) (int, error) {
		return 0, nil
	}

	length, err := svc.GetChainLength(ctx)

	assertNoError(t, err)
	assertEqual(t, 0, length)
	assertEqual(t, 1, repo.CountCallCount())
}

func TestGetChainLength_Error(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.CountFn = func(_ context.Context) (int, error) {
		return 0, errorpkg.ErrInternal
	}

	length, err := svc.GetChainLength(ctx)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, 0, length)
	assertEqual(t, 1, repo.CountCallCount())
}

// -- ValidateChain tests ------------------------------------------------

func TestValidateChain_Empty(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.ListFn = func(_ context.Context) ([]*ent.Block, error) {
		return []*ent.Block{}, nil
	}

	valid, err := svc.ValidateChain(ctx)

	assertNoError(t, err)
	assertTrue(t, valid, "empty chain should be valid")
	assertEqual(t, 1, repo.ListCallCount())
}

func TestValidateChain_SingleBlock(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	blocks := makeValidChain(1)

	repo.ListFn = func(_ context.Context) ([]*ent.Block, error) {
		return blocks, nil
	}

	valid, err := svc.ValidateChain(ctx)

	assertNoError(t, err)
	assertTrue(t, valid, "single block chain should be valid")
	assertEqual(t, 1, repo.ListCallCount())
}

func TestValidateChain_MultipleBlocks(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	blocks := makeValidChain(5)

	repo.ListFn = func(_ context.Context) ([]*ent.Block, error) {
		return blocks, nil
	}

	valid, err := svc.ValidateChain(ctx)

	assertNoError(t, err)
	assertTrue(t, valid, "valid chain of 5 blocks should be valid")
	assertEqual(t, 1, repo.ListCallCount())
}

func TestValidateChain_InvalidHash(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	blocks := makeValidChain(3)
	blocks[1].Hash = "tamperedhash"

	repo.ListFn = func(_ context.Context) ([]*ent.Block, error) {
		return blocks, nil
	}

	valid, err := svc.ValidateChain(ctx)

	assertNoError(t, err)
	assertEqual(t, false, valid)
	assertEqual(t, 1, repo.ListCallCount())
}

func TestValidateChain_BrokenLink(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	blocks := makeValidChain(3)
	blocks[2].PreviousHash = "wronghash"

	repo.ListFn = func(_ context.Context) ([]*ent.Block, error) {
		return blocks, nil
	}

	valid, err := svc.ValidateChain(ctx)

	assertNoError(t, err)
	assertEqual(t, false, valid)
	assertEqual(t, 1, repo.ListCallCount())
}

func TestValidateChain_RepositoryError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.ListFn = func(_ context.Context) ([]*ent.Block, error) {
		return nil, errorpkg.ErrInternal
	}

	valid, err := svc.ValidateChain(ctx)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, false, valid)
	assertEqual(t, 1, repo.ListCallCount())
}

func TestValidateChain_HashMismatch(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	blocks := makeValidChain(2)
	blocks[0].Nonce = 999

	repo.ListFn = func(_ context.Context) ([]*ent.Block, error) {
		return blocks, nil
	}

	valid, err := svc.ValidateChain(ctx)

	assertNoError(t, err)
	assertEqual(t, false, valid)
	assertEqual(t, 1, repo.ListCallCount())
}
