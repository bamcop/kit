package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"skia/pkg/ent_util"
)

type Todos struct {
	ent.Schema
}

func (Todos) Fields() []ent.Field {
	return []ent.Field{
		field.String("created_at").
			StructTag("json:\"created_at\"").
			SchemaType(map[string]string{"mysql": "datetime"}),
		field.String("updated_at").
			StructTag("json:\"updated_at\"").
			SchemaType(map[string]string{"mysql": "datetime"}),
		field.String("deleted_at").
			StructTag("json:\"deleted_at\"").
			SchemaType(map[string]string{"mysql": "datetime"}),
		field.String("task").
			StructTag("json:\"task\"").
			SchemaType(map[string]string{"mysql": "TEXT"}),
		field.Bool("completed").
			StructTag("json:\"completed\"").
			Default(false),
	}
}

func (Todos) Edges() []ent.Edge {
	return nil
}

func (Todos) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ent_util.BaseMixin{},
	}
}

func (Todos) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "todos"},
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
