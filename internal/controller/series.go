package controller

import (
	"context"
	"sort"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/ent/episode"
	"github.com/clustlight/animatrix-api/ent/series"
	"github.com/clustlight/animatrix-api/internal/types"
	"github.com/clustlight/animatrix-api/internal/utils"
)

type seriesWithTimestamp struct {
	series    *ent.Series
	timestamp int64
}

func GetAllSeries(ctx context.Context, client *ent.Client) (*[]types.SeriesResponse, error) {
	series, err := client.Series.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]types.SeriesResponse, 0, len(series))
	for _, s := range series {
		resp := utils.BuildSeriesResponse(s, false, false)
		responses = append(responses, resp)
	}

	return &responses, nil
}

func GetSeries(ctx context.Context, client *ent.Client, seriesID string) (*types.SeriesResponse, error) {
	series, err := client.Series.
		Query().
		Where(series.SeriesIDEQ(seriesID)).
		WithSeasons(func(q *ent.SeasonQuery) {
			q.WithEpisodes(func(eq *ent.EpisodeQuery) {
				eq.Order(ent.Asc("episode_number"))
			}).
				Order(ent.Asc("season_number"))
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	resp := utils.BuildSeriesResponse(series, true, true)
	return &resp, nil
}

func CreateSeries(ctx context.Context, client *ent.Client, req *types.CreateSeriesRequest) (*types.SeriesResponse, error) {
	newSeries, err := client.Series.Create().
		SetSeriesID(req.SeriesID).
		SetTitle(req.Title).
		SetTitleYomi(req.TitleYomi).
		SetTitleEn(req.TitleEn).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	resp := utils.BuildSeriesResponse(newSeries, false, false)
	return &resp, nil
}

func UpdateSeries(ctx context.Context, client *ent.Client, seriesID string, req *types.UpdateSeriesRequest) (*types.SeriesResponse, error) {
	seriesToUpdate, err := client.Series.
		Query().
		Where(series.SeriesIDEQ(seriesID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	upd := seriesToUpdate.Update()
	if req.Title != nil {
		upd = upd.SetTitle(*req.Title)
	}
	if req.TitleYomi != nil {
		upd = upd.SetTitleYomi(*req.TitleYomi)
	}
	if req.TitleEn != nil {
		upd = upd.SetTitleEn(*req.TitleEn)
	}
	updatedSeries, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}

	resp := utils.BuildSeriesResponse(updatedSeries, false, false)
	return &resp, nil
}

func BulkCreateSeries(ctx context.Context, client *ent.Client, seriesList []types.CreateSeriesRequest) ([]types.SeriesResponse, error) {
	bulk := make([]*ent.SeriesCreate, 0, len(seriesList))
	for _, req := range seriesList {
		sc := client.Series.Create().
			SetSeriesID(req.SeriesID).
			SetTitle(req.Title).
			SetTitleYomi(req.TitleYomi).
			SetTitleEn(req.TitleEn)
		bulk = append(bulk, sc)
	}
	created, err := client.Series.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, err
	}
	resps := make([]types.SeriesResponse, 0, len(created))
	for _, s := range created {
		resps = append(resps, utils.BuildSeriesResponse(s, false, false))
	}
	return resps, nil
}

func GetRecentlyUpdatedSeries(ctx context.Context, client *ent.Client) ([]types.SeriesResponse, error) {
	episodes, err := client.Episode.
		Query().
		Order(ent.Desc(episode.FieldTimestamp)).
		WithSeason(func(sq *ent.SeasonQuery) {
			sq.WithSeries()
		}).
		Limit(30).
		All(ctx)
	if err != nil {
		return nil, err
	}

	seriesMap := make(map[string]seriesWithTimestamp)
	for _, ep := range episodes {
		season := ep.Edges.Season
		if season == nil || season.Edges.Series == nil {
			continue
		}
		seriesID := season.Edges.Series.SeriesID
		if _, exists := seriesMap[seriesID]; !exists {
			seriesMap[seriesID] = seriesWithTimestamp{
				series:    season.Edges.Series,
				timestamp: ep.Timestamp.Unix(),
			}
		}
	}

	// Convert map to slice and sort by timestamp
	list := make([]seriesWithTimestamp, 0, len(seriesMap))
	for _, v := range seriesMap {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].timestamp > list[j].timestamp
	})

	// Build response
	responses := make([]types.SeriesResponse, 0, len(list))
	for _, s := range list {
		resp := utils.BuildSeriesResponse(s.series, false, false)
		responses = append(responses, resp)
	}
	return responses, nil
}

func DeleteSeries(ctx context.Context, client *ent.Client, seriesID string) error {
	s, err := client.Series.
		Query().
		Where(series.SeriesIDEQ(seriesID)).
		WithSeasons().
		Only(ctx)
	if err != nil {
		return err
	}
	if len(s.Edges.Seasons) > 0 {
		return ErrHasChildren
	}
	return client.Series.DeleteOneID(s.ID).Exec(ctx)
}
