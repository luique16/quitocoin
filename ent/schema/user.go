package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("name"),
		field.String("email").Unique(),
		field.String("password"),
		field.String("public_id").Unique(),
		field.Time("created_at"),
	}
}
