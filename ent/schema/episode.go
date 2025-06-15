package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Episode holds the schema definition for the Episode entity.
type Episode struct {
	ent.Schema
}

// Fields of the Episode.
func (Episode) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("episode_id").Unique(),
		field.Int("episode_number"),
		field.Float("duration"),
		field.String("duration_string"),
		field.Time("timestamp"),
		field.String("thumbnail"),

		field.String("format_id"),
		field.Int("width"),
		field.Int("height"),
		field.String("dynamic_range"),
		field.Text("metadata"),
	}
}

// Edges of the Episode.
func (Episode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("season", Season.Type).Ref("episodes").Unique().Required(),
	}
}
