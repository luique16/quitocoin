package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/luique16/quitocoin/internal/domain/transaction"
)

type Block struct {
	ent.Schema
}

func (Block) Fields() []ent.Field {
	return []ent.Field{
		field.String("hash").Unique().Immutable(),
		field.Int("index"),
		field.String("previous_hash"),
		field.Int64("nonce"),
		field.String("miner"),
		field.Float("reward"),
		field.JSON("transactions", []transaction.Transaction{}),
		field.Time("created_at"),
	}
}
