// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/clustlight/animatrix-api/ent/season"
	"github.com/clustlight/animatrix-api/ent/series"
)

// SeriesCreate is the builder for creating a Series entity.
type SeriesCreate struct {
	config
	mutation *SeriesMutation
	hooks    []Hook
}

// SetSeriesID sets the "series_id" field.
func (sc *SeriesCreate) SetSeriesID(s string) *SeriesCreate {
	sc.mutation.SetSeriesID(s)
	return sc
}

// SetTitle sets the "title" field.
func (sc *SeriesCreate) SetTitle(s string) *SeriesCreate {
	sc.mutation.SetTitle(s)
	return sc
}

// SetTitleYomi sets the "title_yomi" field.
func (sc *SeriesCreate) SetTitleYomi(s string) *SeriesCreate {
	sc.mutation.SetTitleYomi(s)
	return sc
}

// SetNillableTitleYomi sets the "title_yomi" field if the given value is not nil.
func (sc *SeriesCreate) SetNillableTitleYomi(s *string) *SeriesCreate {
	if s != nil {
		sc.SetTitleYomi(*s)
	}
	return sc
}

// SetTitleEn sets the "title_en" field.
func (sc *SeriesCreate) SetTitleEn(s string) *SeriesCreate {
	sc.mutation.SetTitleEn(s)
	return sc
}

// SetNillableTitleEn sets the "title_en" field if the given value is not nil.
func (sc *SeriesCreate) SetNillableTitleEn(s *string) *SeriesCreate {
	if s != nil {
		sc.SetTitleEn(*s)
	}
	return sc
}

// AddSeasonIDs adds the "seasons" edge to the Season entity by IDs.
func (sc *SeriesCreate) AddSeasonIDs(ids ...int) *SeriesCreate {
	sc.mutation.AddSeasonIDs(ids...)
	return sc
}

// AddSeasons adds the "seasons" edges to the Season entity.
func (sc *SeriesCreate) AddSeasons(s ...*Season) *SeriesCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sc.AddSeasonIDs(ids...)
}

// Mutation returns the SeriesMutation object of the builder.
func (sc *SeriesCreate) Mutation() *SeriesMutation {
	return sc.mutation
}

// Save creates the Series in the database.
func (sc *SeriesCreate) Save(ctx context.Context) (*Series, error) {
	return withHooks(ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SeriesCreate) SaveX(ctx context.Context) *Series {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SeriesCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SeriesCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SeriesCreate) check() error {
	if _, ok := sc.mutation.SeriesID(); !ok {
		return &ValidationError{Name: "series_id", err: errors.New(`ent: missing required field "Series.series_id"`)}
	}
	if _, ok := sc.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`ent: missing required field "Series.title"`)}
	}
	return nil
}

func (sc *SeriesCreate) sqlSave(ctx context.Context) (*Series, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *SeriesCreate) createSpec() (*Series, *sqlgraph.CreateSpec) {
	var (
		_node = &Series{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(series.Table, sqlgraph.NewFieldSpec(series.FieldID, field.TypeInt))
	)
	if value, ok := sc.mutation.SeriesID(); ok {
		_spec.SetField(series.FieldSeriesID, field.TypeString, value)
		_node.SeriesID = value
	}
	if value, ok := sc.mutation.Title(); ok {
		_spec.SetField(series.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if value, ok := sc.mutation.TitleYomi(); ok {
		_spec.SetField(series.FieldTitleYomi, field.TypeString, value)
		_node.TitleYomi = value
	}
	if value, ok := sc.mutation.TitleEn(); ok {
		_spec.SetField(series.FieldTitleEn, field.TypeString, value)
		_node.TitleEn = value
	}
	if nodes := sc.mutation.SeasonsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   series.SeasonsTable,
			Columns: []string{series.SeasonsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(season.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// SeriesCreateBulk is the builder for creating many Series entities in bulk.
type SeriesCreateBulk struct {
	config
	err      error
	builders []*SeriesCreate
}

// Save creates the Series entities in the database.
func (scb *SeriesCreateBulk) Save(ctx context.Context) ([]*Series, error) {
	if scb.err != nil {
		return nil, scb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Series, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SeriesMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SeriesCreateBulk) SaveX(ctx context.Context) []*Series {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SeriesCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SeriesCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}
