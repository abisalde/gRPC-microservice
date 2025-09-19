package schema

import (
	"regexp"
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
			StorageKey("id").
			Annotations(
				entproto.Field(1),
			),

		field.String("email").
			Unique().
			Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)).
			Annotations(
				entproto.Field(2),
			),

		field.String("password_hash").
			Sensitive().
			Optional().
			Annotations(
				entproto.Field(3),
			),

		field.String("first_name").
			MaxLen(50).Default("").
			StructTag(`json:"firstName"`).
			Annotations(
				entproto.Field(4),
			),

		field.String("last_name").
			MaxLen(50).Default("").
			StructTag(`json:"lastName"`).
			Annotations(
				entproto.Field(5),
			),

		field.Enum("role").
			Values("USER", "ADMIN").
			Default("USER").
			Comment("User role defining access permissions and privileges").
			Annotations(
				entproto.Enum(map[string]int32{
					"USER":  0,
					"ADMIN": 1,
				}),
				entproto.Field(6),
			),

		field.Bool("is_email_verified").
			Default(false).
			StructTag(`json:"isEmailVerified"`).
			Annotations(
				entproto.Field(7),
			),
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
			StructTag(`json:"createdAt"`).
			Annotations(
				entproto.Field(8),
			),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updatedAt"`).
			Annotations(
				entproto.Field(9),
			),

		field.Time("deleted_at").
			Optional().
			Nillable().
			StructTag(`json:"deletedAt"`).
			Annotations(
				entproto.Field(10),
			),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email"),
		index.Fields("is_email_verified"),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(
			entproto.PackageName("auth_pbuf"),
		),
		entproto.Service(),
	}
}
