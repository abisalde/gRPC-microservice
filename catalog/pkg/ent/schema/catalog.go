package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Catalog struct {
	ent.Schema
}

func (Catalog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Catalog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			StorageKey("id").
			Annotations(
				entproto.Field(1),
			),
		field.String("name").Annotations(entproto.Field(2)),
		field.String("description").Optional().Annotations(entproto.Field(3)),
		field.Float("price").Annotations(entproto.Field(4)),
	}
}

func (Catalog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(
			entproto.PackageName("catalog_pbuf"),
		),
		entproto.Service(),
	}
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
