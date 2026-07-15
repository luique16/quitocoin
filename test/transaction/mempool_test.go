package transaction_test

import (
	"testing"

	"github.com/luique16/quitocoin/internal/domain/transaction"
)

func newMP() transaction.MemPool {
	return transaction.NewMemPool()
}

func TestCreateTransaction(t *testing.T) {
	mp := newMP()

	result := mp.CreateTransaction("alice", "bob", 100.0)

	if result.From != "alice" {
		t.Errorf("expected From=alice, got %s", result.From)
	}
	if result.To != "bob" {
		t.Errorf("expected To=bob, got %s", result.To)
	}
	if result.Amount != 100.0 {
		t.Errorf("expected Amount=100.0, got %f", result.Amount)
	}
}

func TestPushTransaction(t *testing.T) {
	t.Run("push single transaction", func(t *testing.T) {
		mp := newMP()

		tt := mp.CreateTransaction("alice", "bob", 50.0)
		mp.PushTransaction(tt)

		if mp.Count() != 1 {
			t.Errorf("expected count=1, got %d", mp.Count())
		}
	})

	t.Run("push multiple transactions", func(t *testing.T) {
		mp := newMP()

		for i := 0; i < 3; i++ {
			mp.PushTransaction(mp.CreateTransaction("alice", "bob", float32(i+1)*10.0))
		}

		if mp.Count() != 3 {
			t.Errorf("expected count=3, got %d", mp.Count())
		}
	})
}

func TestPullFirstTransactions(t *testing.T) {
	t.Run("pull one transaction", func(t *testing.T) {
		mp := newMP()
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))

		results := mp.PullFirstTransactions(1)

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].Amount != 100.0 {
			t.Errorf("expected Amount=100.0, got %f", results[0].Amount)
		}
		if mp.Count() != 1 {
			t.Errorf("expected count=1 (pull should not remove), got %d", mp.Count())
		}
	})

	t.Run("pull multiple transactions", func(t *testing.T) {
		mp := newMP()
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))
		mp.PushTransaction(mp.CreateTransaction("alice", "carol", 200.0))
		mp.PushTransaction(mp.CreateTransaction("alice", "david", 300.0))

		results := mp.PullFirstTransactions(2)

		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
		if results[0].Amount != 100.0 {
			t.Errorf("expected first Amount=100.0, got %f", results[0].Amount)
		}
		if results[1].Amount != 200.0 {
			t.Errorf("expected second Amount=200.0, got %f", results[1].Amount)
		}
		if mp.Count() != 3 {
			t.Errorf("expected count=3 (pull should not remove), got %d", mp.Count())
		}
	})

	t.Run("pull from empty pool", func(t *testing.T) {
		mp := newMP()
		results := mp.PullFirstTransactions(5)

		if len(results) != 0 {
			t.Errorf("expected 0 results from empty pool, got %d", len(results))
		}
	})

	t.Run("pull zero transactions", func(t *testing.T) {
		mp := newMP()
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))

		results := mp.PullFirstTransactions(0)

		if len(results) != 0 {
			t.Errorf("expected 0 results for pull 0, got %d", len(results))
		}
		if mp.Count() != 1 {
			t.Errorf("expected count=1 after pull 0, got %d", mp.Count())
		}
	})

	t.Run("pull more than available", func(t *testing.T) {
		mp := newMP()
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))

		results := mp.PullFirstTransactions(5)

		if len(results) != 1 {
			t.Errorf("expected 1 result (limited by available), got %d", len(results))
		}
		if mp.Count() != 1 {
			t.Errorf("expected count=1 (pull should not remove), got %d", mp.Count())
		}
	})
}

func TestDeleteFirstTransactions(t *testing.T) {
	t.Run("delete one transaction", func(t *testing.T) {
		mp := newMP()
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))

		mp.DeleteFirstTransactions(1)

		if mp.Count() != 0 {
			t.Errorf("expected count=0 after delete 1, got %d", mp.Count())
		}
	})

	t.Run("delete multiple transactions", func(t *testing.T) {
		mp := newMP()
		for i := 0; i < 5; i++ {
			mp.PushTransaction(mp.CreateTransaction("alice", "bob", float32(i)*10.0))
		}

		mp.DeleteFirstTransactions(3)

		if mp.Count() != 2 {
			t.Errorf("expected count=2 after delete 3, got %d", mp.Count())
		}
	})

	t.Run("delete from empty pool", func(t *testing.T) {
		mp := newMP()
		mp.DeleteFirstTransactions(5)

		if mp.Count() != 0 {
			t.Errorf("expected count=0 from empty pool, got %d", mp.Count())
		}
	})

	t.Run("delete zero transactions", func(t *testing.T) {
		mp := newMP()
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))

		mp.DeleteFirstTransactions(0)

		if mp.Count() != 1 {
			t.Errorf("expected count=1 after delete 0, got %d", mp.Count())
		}
	})

	t.Run("delete negative transactions", func(t *testing.T) {
		mp := newMP()
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))

		mp.DeleteFirstTransactions(-1)

		if mp.Count() != 1 {
			t.Errorf("expected count=1 after delete negative, got %d", mp.Count())
		}
	})
}

func TestCount(t *testing.T) {
	mp := newMP()

	if mp.Count() != 0 {
		t.Errorf("expected count=0 initially, got %d", mp.Count())
	}

	mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))
	if mp.Count() != 1 {
		t.Errorf("expected count=1 after 1 push, got %d", mp.Count())
	}

	mp.PushTransaction(mp.CreateTransaction("alice", "carol", 200.0))
	if mp.Count() != 2 {
		t.Errorf("expected count=2 after 2 pushes, got %d", mp.Count())
	}

	mp.PullFirstTransactions(1)
	if mp.Count() != 2 {
		t.Errorf("expected count=2 (pull should not remove), got %d", mp.Count())
	}

	mp.DeleteFirstTransactions(1)
	if mp.Count() != 1 {
		t.Errorf("expected count=1 after delete 1, got %d", mp.Count())
	}
}

func TestClear(t *testing.T) {
	mp := newMP()

	for i := 0; i < 5; i++ {
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", float32(i)))
	}

	if mp.Count() != 5 {
		t.Errorf("expected count=5 before clear, got %d", mp.Count())
	}

	mp.Clear()

	if mp.Count() != 0 {
		t.Errorf("expected count=0 after clear, got %d", mp.Count())
	}
}

func TestPullReturnsCopy(t *testing.T) {
	mp := newMP()
	mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))
	mp.PushTransaction(mp.CreateTransaction("alice", "carol", 200.0))

	results := mp.PullFirstTransactions(1)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].From != "alice" {
		t.Errorf("expected From=alice, got %s", results[0].From)
	}
	if results[0].To != "bob" {
		t.Errorf("expected To=bob, got %s", results[0].To)
	}
	if results[0].Amount != 100.0 {
		t.Errorf("expected Amount=100.0, got %f", results[0].Amount)
	}
	if mp.Count() != 2 {
		t.Errorf("expected count=2 (pull should not remove), got %d", mp.Count())
	}
}

func TestTransactionOrder(t *testing.T) {
	t.Run("pull returns first n in order", func(t *testing.T) {
		mp := newMP()
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))
		mp.PushTransaction(mp.CreateTransaction("alice", "carol", 200.0))
		mp.PushTransaction(mp.CreateTransaction("alice", "david", 300.0))

		results := mp.PullFirstTransactions(3)

		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}
		if results[0].To != "bob" {
			t.Errorf("expected first To=bob, got %s", results[0].To)
		}
		if results[1].To != "carol" {
			t.Errorf("expected second To=carol, got %s", results[1].To)
		}
		if results[2].To != "david" {
			t.Errorf("expected third To=david, got %s", results[2].To)
		}
	})

	t.Run("pull does not modify pool", func(t *testing.T) {
		mp := newMP()
		mp.PushTransaction(mp.CreateTransaction("alice", "bob", 100.0))
		mp.PushTransaction(mp.CreateTransaction("alice", "carol", 200.0))

		mp.PullFirstTransactions(1)
		mp.PullFirstTransactions(1)

		results := mp.PullFirstTransactions(2)

		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
		if results[0].To != "bob" {
			t.Errorf("expected first To=bob, got %s", results[0].To)
		}
		if results[1].To != "carol" {
			t.Errorf("expected second To=carol, got %s", results[1].To)
		}
		if mp.Count() != 2 {
			t.Errorf("expected count=2 (pull should not remove), got %d", mp.Count())
		}
	})
}
