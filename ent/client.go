// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/clustlight/animatrix-api/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/clustlight/animatrix-api/ent/episode"
	"github.com/clustlight/animatrix-api/ent/season"
	"github.com/clustlight/animatrix-api/ent/series"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Episode is the client for interacting with the Episode builders.
	Episode *EpisodeClient
	// Season is the client for interacting with the Season builders.
	Season *SeasonClient
	// Series is the client for interacting with the Series builders.
	Series *SeriesClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Episode = NewEpisodeClient(c.config)
	c.Season = NewSeasonClient(c.config)
	c.Series = NewSeriesClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:     ctx,
		config:  cfg,
		Episode: NewEpisodeClient(cfg),
		Season:  NewSeasonClient(cfg),
		Series:  NewSeriesClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:     ctx,
		config:  cfg,
		Episode: NewEpisodeClient(cfg),
		Season:  NewSeasonClient(cfg),
		Series:  NewSeriesClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Episode.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Episode.Use(hooks...)
	c.Season.Use(hooks...)
	c.Series.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Episode.Intercept(interceptors...)
	c.Season.Intercept(interceptors...)
	c.Series.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *EpisodeMutation:
		return c.Episode.mutate(ctx, m)
	case *SeasonMutation:
		return c.Season.mutate(ctx, m)
	case *SeriesMutation:
		return c.Series.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// EpisodeClient is a client for the Episode schema.
type EpisodeClient struct {
	config
}

// NewEpisodeClient returns a client for the Episode from the given config.
func NewEpisodeClient(c config) *EpisodeClient {
	return &EpisodeClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `episode.Hooks(f(g(h())))`.
func (c *EpisodeClient) Use(hooks ...Hook) {
	c.hooks.Episode = append(c.hooks.Episode, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `episode.Intercept(f(g(h())))`.
func (c *EpisodeClient) Intercept(interceptors ...Interceptor) {
	c.inters.Episode = append(c.inters.Episode, interceptors...)
}

// Create returns a builder for creating a Episode entity.
func (c *EpisodeClient) Create() *EpisodeCreate {
	mutation := newEpisodeMutation(c.config, OpCreate)
	return &EpisodeCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Episode entities.
func (c *EpisodeClient) CreateBulk(builders ...*EpisodeCreate) *EpisodeCreateBulk {
	return &EpisodeCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *EpisodeClient) MapCreateBulk(slice any, setFunc func(*EpisodeCreate, int)) *EpisodeCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &EpisodeCreateBulk{err: fmt.Errorf("calling to EpisodeClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*EpisodeCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &EpisodeCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Episode.
func (c *EpisodeClient) Update() *EpisodeUpdate {
	mutation := newEpisodeMutation(c.config, OpUpdate)
	return &EpisodeUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EpisodeClient) UpdateOne(e *Episode) *EpisodeUpdateOne {
	mutation := newEpisodeMutation(c.config, OpUpdateOne, withEpisode(e))
	return &EpisodeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *EpisodeClient) UpdateOneID(id int) *EpisodeUpdateOne {
	mutation := newEpisodeMutation(c.config, OpUpdateOne, withEpisodeID(id))
	return &EpisodeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Episode.
func (c *EpisodeClient) Delete() *EpisodeDelete {
	mutation := newEpisodeMutation(c.config, OpDelete)
	return &EpisodeDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *EpisodeClient) DeleteOne(e *Episode) *EpisodeDeleteOne {
	return c.DeleteOneID(e.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *EpisodeClient) DeleteOneID(id int) *EpisodeDeleteOne {
	builder := c.Delete().Where(episode.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EpisodeDeleteOne{builder}
}

// Query returns a query builder for Episode.
func (c *EpisodeClient) Query() *EpisodeQuery {
	return &EpisodeQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeEpisode},
		inters: c.Interceptors(),
	}
}

// Get returns a Episode entity by its id.
func (c *EpisodeClient) Get(ctx context.Context, id int) (*Episode, error) {
	return c.Query().Where(episode.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EpisodeClient) GetX(ctx context.Context, id int) *Episode {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QuerySeason queries the season edge of a Episode.
func (c *EpisodeClient) QuerySeason(e *Episode) *SeasonQuery {
	query := (&SeasonClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(episode.Table, episode.FieldID, id),
			sqlgraph.To(season.Table, season.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, episode.SeasonTable, episode.SeasonColumn),
		)
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EpisodeClient) Hooks() []Hook {
	return c.hooks.Episode
}

// Interceptors returns the client interceptors.
func (c *EpisodeClient) Interceptors() []Interceptor {
	return c.inters.Episode
}

func (c *EpisodeClient) mutate(ctx context.Context, m *EpisodeMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&EpisodeCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&EpisodeUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&EpisodeUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&EpisodeDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Episode mutation op: %q", m.Op())
	}
}

// SeasonClient is a client for the Season schema.
type SeasonClient struct {
	config
}

// NewSeasonClient returns a client for the Season from the given config.
func NewSeasonClient(c config) *SeasonClient {
	return &SeasonClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `season.Hooks(f(g(h())))`.
func (c *SeasonClient) Use(hooks ...Hook) {
	c.hooks.Season = append(c.hooks.Season, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `season.Intercept(f(g(h())))`.
func (c *SeasonClient) Intercept(interceptors ...Interceptor) {
	c.inters.Season = append(c.inters.Season, interceptors...)
}

// Create returns a builder for creating a Season entity.
func (c *SeasonClient) Create() *SeasonCreate {
	mutation := newSeasonMutation(c.config, OpCreate)
	return &SeasonCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Season entities.
func (c *SeasonClient) CreateBulk(builders ...*SeasonCreate) *SeasonCreateBulk {
	return &SeasonCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *SeasonClient) MapCreateBulk(slice any, setFunc func(*SeasonCreate, int)) *SeasonCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &SeasonCreateBulk{err: fmt.Errorf("calling to SeasonClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*SeasonCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &SeasonCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Season.
func (c *SeasonClient) Update() *SeasonUpdate {
	mutation := newSeasonMutation(c.config, OpUpdate)
	return &SeasonUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SeasonClient) UpdateOne(s *Season) *SeasonUpdateOne {
	mutation := newSeasonMutation(c.config, OpUpdateOne, withSeason(s))
	return &SeasonUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SeasonClient) UpdateOneID(id int) *SeasonUpdateOne {
	mutation := newSeasonMutation(c.config, OpUpdateOne, withSeasonID(id))
	return &SeasonUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Season.
func (c *SeasonClient) Delete() *SeasonDelete {
	mutation := newSeasonMutation(c.config, OpDelete)
	return &SeasonDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SeasonClient) DeleteOne(s *Season) *SeasonDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SeasonClient) DeleteOneID(id int) *SeasonDeleteOne {
	builder := c.Delete().Where(season.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SeasonDeleteOne{builder}
}

// Query returns a query builder for Season.
func (c *SeasonClient) Query() *SeasonQuery {
	return &SeasonQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSeason},
		inters: c.Interceptors(),
	}
}

// Get returns a Season entity by its id.
func (c *SeasonClient) Get(ctx context.Context, id int) (*Season, error) {
	return c.Query().Where(season.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SeasonClient) GetX(ctx context.Context, id int) *Season {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QuerySeries queries the series edge of a Season.
func (c *SeasonClient) QuerySeries(s *Season) *SeriesQuery {
	query := (&SeriesClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(season.Table, season.FieldID, id),
			sqlgraph.To(series.Table, series.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, season.SeriesTable, season.SeriesColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEpisodes queries the episodes edge of a Season.
func (c *SeasonClient) QueryEpisodes(s *Season) *EpisodeQuery {
	query := (&EpisodeClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(season.Table, season.FieldID, id),
			sqlgraph.To(episode.Table, episode.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, season.EpisodesTable, season.EpisodesColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SeasonClient) Hooks() []Hook {
	return c.hooks.Season
}

// Interceptors returns the client interceptors.
func (c *SeasonClient) Interceptors() []Interceptor {
	return c.inters.Season
}

func (c *SeasonClient) mutate(ctx context.Context, m *SeasonMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SeasonCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SeasonUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SeasonUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SeasonDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Season mutation op: %q", m.Op())
	}
}

// SeriesClient is a client for the Series schema.
type SeriesClient struct {
	config
}

// NewSeriesClient returns a client for the Series from the given config.
func NewSeriesClient(c config) *SeriesClient {
	return &SeriesClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `series.Hooks(f(g(h())))`.
func (c *SeriesClient) Use(hooks ...Hook) {
	c.hooks.Series = append(c.hooks.Series, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `series.Intercept(f(g(h())))`.
func (c *SeriesClient) Intercept(interceptors ...Interceptor) {
	c.inters.Series = append(c.inters.Series, interceptors...)
}

// Create returns a builder for creating a Series entity.
func (c *SeriesClient) Create() *SeriesCreate {
	mutation := newSeriesMutation(c.config, OpCreate)
	return &SeriesCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Series entities.
func (c *SeriesClient) CreateBulk(builders ...*SeriesCreate) *SeriesCreateBulk {
	return &SeriesCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *SeriesClient) MapCreateBulk(slice any, setFunc func(*SeriesCreate, int)) *SeriesCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &SeriesCreateBulk{err: fmt.Errorf("calling to SeriesClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*SeriesCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &SeriesCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Series.
func (c *SeriesClient) Update() *SeriesUpdate {
	mutation := newSeriesMutation(c.config, OpUpdate)
	return &SeriesUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SeriesClient) UpdateOne(s *Series) *SeriesUpdateOne {
	mutation := newSeriesMutation(c.config, OpUpdateOne, withSeries(s))
	return &SeriesUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SeriesClient) UpdateOneID(id int) *SeriesUpdateOne {
	mutation := newSeriesMutation(c.config, OpUpdateOne, withSeriesID(id))
	return &SeriesUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Series.
func (c *SeriesClient) Delete() *SeriesDelete {
	mutation := newSeriesMutation(c.config, OpDelete)
	return &SeriesDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SeriesClient) DeleteOne(s *Series) *SeriesDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SeriesClient) DeleteOneID(id int) *SeriesDeleteOne {
	builder := c.Delete().Where(series.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SeriesDeleteOne{builder}
}

// Query returns a query builder for Series.
func (c *SeriesClient) Query() *SeriesQuery {
	return &SeriesQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSeries},
		inters: c.Interceptors(),
	}
}

// Get returns a Series entity by its id.
func (c *SeriesClient) Get(ctx context.Context, id int) (*Series, error) {
	return c.Query().Where(series.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SeriesClient) GetX(ctx context.Context, id int) *Series {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QuerySeasons queries the seasons edge of a Series.
func (c *SeriesClient) QuerySeasons(s *Series) *SeasonQuery {
	query := (&SeasonClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(series.Table, series.FieldID, id),
			sqlgraph.To(season.Table, season.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, series.SeasonsTable, series.SeasonsColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SeriesClient) Hooks() []Hook {
	return c.hooks.Series
}

// Interceptors returns the client interceptors.
func (c *SeriesClient) Interceptors() []Interceptor {
	return c.inters.Series
}

func (c *SeriesClient) mutate(ctx context.Context, m *SeriesMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SeriesCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SeriesUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SeriesUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SeriesDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Series mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Episode, Season, Series []ent.Hook
	}
	inters struct {
		Episode, Season, Series []ent.Interceptor
	}
)
