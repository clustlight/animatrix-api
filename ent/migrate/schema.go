// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// EpisodesColumns holds the columns for the "episodes" table.
	EpisodesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "title", Type: field.TypeString},
		{Name: "episode_id", Type: field.TypeString, Unique: true},
		{Name: "episode_number", Type: field.TypeInt},
		{Name: "duration", Type: field.TypeFloat64},
		{Name: "duration_string", Type: field.TypeString},
		{Name: "timestamp", Type: field.TypeTime},
		{Name: "format_id", Type: field.TypeString},
		{Name: "width", Type: field.TypeInt},
		{Name: "height", Type: field.TypeInt},
		{Name: "dynamic_range", Type: field.TypeString},
		{Name: "metadata", Type: field.TypeString, Size: 2147483647},
		{Name: "season_episodes", Type: field.TypeInt},
	}
	// EpisodesTable holds the schema information for the "episodes" table.
	EpisodesTable = &schema.Table{
		Name:       "episodes",
		Columns:    EpisodesColumns,
		PrimaryKey: []*schema.Column{EpisodesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "episodes_seasons_episodes",
				Columns:    []*schema.Column{EpisodesColumns[12]},
				RefColumns: []*schema.Column{SeasonsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// SeasonsColumns holds the columns for the "seasons" table.
	SeasonsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "season_id", Type: field.TypeString, Unique: true},
		{Name: "season_title", Type: field.TypeString},
		{Name: "season_title_yomi", Type: field.TypeString, Nullable: true},
		{Name: "season_number", Type: field.TypeInt},
		{Name: "shoboi_tid", Type: field.TypeInt, Nullable: true},
		{Name: "description", Type: field.TypeString, Nullable: true, Size: 2147483647},
		{Name: "first_year", Type: field.TypeInt, Nullable: true},
		{Name: "first_month", Type: field.TypeInt, Nullable: true},
		{Name: "first_end_year", Type: field.TypeInt, Nullable: true},
		{Name: "first_end_month", Type: field.TypeInt, Nullable: true},
		{Name: "series_seasons", Type: field.TypeInt},
	}
	// SeasonsTable holds the schema information for the "seasons" table.
	SeasonsTable = &schema.Table{
		Name:       "seasons",
		Columns:    SeasonsColumns,
		PrimaryKey: []*schema.Column{SeasonsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "seasons_series_seasons",
				Columns:    []*schema.Column{SeasonsColumns[11]},
				RefColumns: []*schema.Column{SeriesColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// SeriesColumns holds the columns for the "series" table.
	SeriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "series_id", Type: field.TypeString, Unique: true},
		{Name: "title", Type: field.TypeString},
		{Name: "title_yomi", Type: field.TypeString, Nullable: true},
		{Name: "title_en", Type: field.TypeString, Nullable: true},
	}
	// SeriesTable holds the schema information for the "series" table.
	SeriesTable = &schema.Table{
		Name:       "series",
		Columns:    SeriesColumns,
		PrimaryKey: []*schema.Column{SeriesColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		EpisodesTable,
		SeasonsTable,
		SeriesTable,
	}
)

func init() {
	EpisodesTable.ForeignKeys[0].RefTable = SeasonsTable
	SeasonsTable.ForeignKeys[0].RefTable = SeriesTable
}
