package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Catalog struct {
	ent.Schema
}

func (Catalog) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Annotations(entproto.Field(2)),
		field.String("description").Optional().Annotations(entproto.Field(3)),
		field.Float("price").Annotations(entproto.Field(4)),
	}
}

func (Catalog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
