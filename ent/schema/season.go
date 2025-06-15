package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Season holds the schema definition for the Season entity.
type Season struct {
	ent.Schema
}

// Fields of the Season.
func (Season) Fields() []ent.Field {
	return []ent.Field{
		field.String("season_id").Unique().NotEmpty(),
		field.String("season_title").NotEmpty(),
		field.String("season_title_yomi").Optional(),
		field.Int("season_number"),
		field.Int("shoboi_tid").Optional(),
		field.Text("description").Optional(), // Comment
		field.Int("first_year").Optional(),
		field.Int("first_month").Optional(),
		field.Int("first_end_year").Optional(),
		field.Int("first_end_month").Optional(),
	}
}

// Edges of the Season.
func (Season) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("series", Series.Type).Ref("seasons").Unique().Required(),
		edge.To("episodes", Episode.Type),
	}
}
