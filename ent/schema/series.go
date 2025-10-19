package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Series holds the schema definition for the Series entity.
type Series struct {
	ent.Schema
}

// Fields of the Series.
func (Series) Fields() []ent.Field {
	return []ent.Field{
		field.String("series_id").Unique().NotEmpty(),
		field.String("title").NotEmpty(),
		field.String("title_yomi").Optional(),
		field.String("title_en").Optional(),
		field.Text("description").Optional(),
	}
}

// Edges of the Series.
func (Series) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("seasons", Season.Type),
	}
}
