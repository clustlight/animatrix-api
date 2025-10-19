package controller

import (
	"context"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/ent/season"
	"github.com/clustlight/animatrix-api/ent/series"
	"github.com/clustlight/animatrix-api/internal/types"
	"github.com/clustlight/animatrix-api/internal/utils"
)

func GetAllSeasons(ctx context.Context, client *ent.Client) (*[]types.SeasonResponse, error) {
	seasons, err := client.Season.Query().
		WithEpisodes().
		Order(ent.Asc("season_number")).
		All(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]types.SeasonResponse, 0, len(seasons))
	for _, s := range seasons {
		resp := utils.BuildSeasonResponse(s, false)
		responses = append(responses, resp)
	}

	return &responses, nil
}

func GetSeason(ctx context.Context, client *ent.Client, seasonID string) (*types.SeasonResponse, error) {
	season, err := client.Season.
		Query().
		Where(season.SeasonIDEQ(seasonID)).
		WithEpisodes(func(q *ent.EpisodeQuery) {
			q.Order(ent.Asc("episode_number"))
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	resp := utils.BuildSeasonResponse(season, true)
	return &resp, nil
}

func CreateSeason(ctx context.Context, client *ent.Client, req *types.CreateSeasonRequest) (*types.SeasonResponse, error) {
	series, err := client.Series.
		Query().
		Where(series.SeriesIDEQ(req.SeriesID)).
		Only(ctx)
	if err != nil {
		return nil, err // Seriesが見つからない場合はエラー
	}

	newSeason := client.Season.Create().
		SetSeries(series).
		SetSeasonID(req.SeasonID).
		SetSeasonTitle(req.SeasonTitle).
		SetSeasonNumber(req.SeasonNumber)

	if req.SeasonTitleYomi != nil {
		newSeason = newSeason.SetSeasonTitleYomi(*req.SeasonTitleYomi)
	}
	if req.ShoboiTID != nil {
		newSeason = newSeason.SetShoboiTid(*req.ShoboiTID)
	}
	if req.Description != nil {
		newSeason = newSeason.SetDescription(*req.Description)
	}
	if req.FirstYear != nil {
		newSeason = newSeason.SetFirstYear(*req.FirstYear)
	}
	if req.FirstMonth != nil {
		newSeason = newSeason.SetFirstMonth(*req.FirstMonth)
	}
	if req.FirstEndYear != nil {
		newSeason = newSeason.SetFirstEndYear(*req.FirstEndYear)
	}
	if req.FirstEndMonth != nil {
		newSeason = newSeason.SetFirstEndMonth(*req.FirstEndMonth)
	}

	saved, err := newSeason.Save(ctx)
	if err != nil {
		return nil, err
	}

	resp := utils.BuildSeasonResponse(saved, true)
	return &resp, nil
}

func UpdateSeason(ctx context.Context, client *ent.Client, seasonID string, req *types.UpdateSeasonRequest) (*types.SeasonResponse, error) {
	season, err := client.Season.
		Query().
		Where(season.SeasonIDEQ(seasonID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	update := season.Update()

	if req.SeasonTitle != nil {
		update.SetSeasonTitle(*req.SeasonTitle)
	}
	if req.SeasonTitleYomi != nil {
		update.SetSeasonTitleYomi(*req.SeasonTitleYomi)
	}
	if req.SeasonNumber != nil {
		update.SetSeasonNumber(*req.SeasonNumber)
	}
	if req.ShoboiTID != nil {
		update.SetShoboiTid(*req.ShoboiTID)
	}
	if req.Description != nil {
		update.SetDescription(*req.Description)
	}
	if req.FirstYear != nil {
		update.SetFirstYear(*req.FirstYear)
	}
	if req.FirstMonth != nil {
		update.SetFirstMonth(*req.FirstMonth)
	}
	if req.FirstEndYear != nil {
		update.SetFirstEndYear(*req.FirstEndYear)
	}
	if req.FirstEndMonth != nil {
		update.SetFirstEndMonth(*req.FirstEndMonth)
	}

	if req.SeriesID != nil {
		srs, err := client.Series.
			Query().
			Where(series.SeriesIDEQ(*req.SeriesID)).
			Only(ctx)
		if err != nil {
			return nil, err
		}
		update.SetSeries(srs)
	}

	saved, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	resp := utils.BuildSeasonResponse(saved, true)
	return &resp, nil
}

func BulkCreateSeason(ctx context.Context, client *ent.Client, seasonList []types.CreateSeasonRequest) ([]types.SeasonResponse, error) {
	bulk := make([]*ent.SeasonCreate, 0, len(seasonList))
	for _, req := range seasonList {
		series, err := client.Series.
			Query().
			Where(series.SeriesIDEQ(req.SeriesID)).
			Only(ctx)
		if err != nil {
			return nil, err
		}

		sc := client.Season.Create().
			SetSeasonID(req.SeasonID).
			SetSeasonTitle(req.SeasonTitle).
			SetSeasonNumber(req.SeasonNumber).
			SetSeries(series)
		if req.SeasonTitleYomi != nil {
			sc = sc.SetSeasonTitleYomi(*req.SeasonTitleYomi)
		}
		if req.ShoboiTID != nil {
			sc = sc.SetShoboiTid(*req.ShoboiTID)
		}
		if req.Description != nil {
			sc = sc.SetDescription(*req.Description)
		}
		if req.FirstYear != nil {
			sc = sc.SetFirstYear(*req.FirstYear)
		}
		if req.FirstMonth != nil {
			sc = sc.SetFirstMonth(*req.FirstMonth)
		}
		if req.FirstEndYear != nil {
			sc = sc.SetFirstEndYear(*req.FirstEndYear)
		}
		if req.FirstEndMonth != nil {
			sc = sc.SetFirstEndMonth(*req.FirstEndMonth)
		}
		bulk = append(bulk, sc)
	}
	created, err := client.Season.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, err
	}
	resps := make([]types.SeasonResponse, 0, len(created))
	for _, s := range created {
		resps = append(resps, utils.BuildSeasonResponse(s, false))
	}
	return resps, nil
}

func DeleteSeason(ctx context.Context, client *ent.Client, seasonID string) error {
	s, err := client.Season.
		Query().
		Where(season.SeasonIDEQ(seasonID)).
		WithEpisodes().
		Only(ctx)
	if err != nil {
		return err
	}
	if len(s.Edges.Episodes) > 0 {
		return ErrHasChildren
	}
	return client.Season.DeleteOneID(s.ID).Exec(ctx)
}
