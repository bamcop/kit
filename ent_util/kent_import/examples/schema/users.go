package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"skia/pkg/ent_util"
)

type Users struct {
	ent.Schema
}

func (Users) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			StructTag("json:\"name\"").
			SchemaType(map[string]string{"mysql": "TEXT"}),
		field.String("channel").
			StructTag("json:\"channel\"").
			SchemaType(map[string]string{"mysql": "TEXT"}),
		field.String("channel_uid").
			StructTag("json:\"channel_uid\"").
			SchemaType(map[string]string{"mysql": "TEXT"}),
	}
}

func (Users) Edges() []ent.Edge {
	return nil
}

func (Users) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ent_util.BaseMixin{},
	}
}

func (Users) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "users"},
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
