package controller

import (
	"context"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/ent/series"
	"github.com/clustlight/animatrix-api/internal/types"
	"github.com/clustlight/animatrix-api/internal/utils"
)

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
