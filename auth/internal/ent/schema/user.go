package schema

import (
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Immutable().
			StorageKey("id"),

		field.String("email").
			Unique().
			Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)),

		field.String("password_hash").
			Sensitive().
			Optional(),

		field.String("first_name").
			MaxLen(50).Default("").
			StructTag(`json:"firstName"`),

		field.String("last_name").
			MaxLen(50).Default("").
			StructTag(`json:"lastName"`),

		field.Enum("role").
			Values("ADMIN", "USER").
			Default("USER"),

		field.Bool("is_email_verified").
			Default(false).
			StructTag(`json:"isEmailVerified"`),
	}

}

func (User) Edges() []ent.Edge {
	return nil
}

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			StructTag(`json:"createdAt"`),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updatedAt"`),

		field.Time("deleted_at").
			Optional().
			Nillable().
			StructTag(`json:"deletedAt"`),
	}
}
