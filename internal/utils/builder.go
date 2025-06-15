package utils

import (
	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/internal/types"
)

func BuildSeriesResponse(series *ent.Series, withSeasons, withEpisodes bool) types.SeriesResponse {
	resp := types.SeriesResponse{
		SeriesID:  series.SeriesID,
		Title:     series.Title,
		TitleYomi: series.TitleYomi,
		TitleEn:   series.TitleEn,
	}
	if withSeasons && series.Edges.Seasons != nil {
		seasons := make([]types.SeasonResponse, 0, len(series.Edges.Seasons))
		for _, season := range series.Edges.Seasons {
			seasons = append(seasons, BuildSeasonResponse(season, withEpisodes))
		}
		resp.Seasons = seasons
	}
	return resp
}

func BuildSeasonResponse(season *ent.Season, withEpisodes bool) types.SeasonResponse {
	resp := types.SeasonResponse{
		SeasonID:        season.SeasonID,
		SeasonTitle:     season.SeasonTitle,
		SeasonTitleYomi: season.SeasonTitleYomi,
		SeasonNumber:    season.SeasonNumber,
		ShoboiTID:       season.ShoboiTid,
		Description:     season.Description,
		FirstYear:       season.FirstYear,
		FirstMonth:      season.FirstMonth,
		FirstEndYear:    season.FirstEndYear,
		FirstEndMonth:   season.FirstEndMonth,
	}
	if withEpisodes && season.Edges.Episodes != nil {
		episodes := make([]types.EpisodeResponse, 0, len(season.Edges.Episodes))
		for _, ep := range season.Edges.Episodes {
			episodes = append(episodes, BuildEpisodeResponse(ep))
		}
		resp.Episodes = episodes
	}
	return resp
}

func BuildEpisodeResponse(ep *ent.Episode) types.EpisodeResponse {
	return types.EpisodeResponse{
		Title:          ep.Title,
		EpisodeID:      ep.EpisodeID,
		EpisodeNumber:  ep.EpisodeNumber,
		Duration:       ep.Duration,
		DurationString: ep.DurationString,
		Timestamp:      ep.Timestamp,
		Thumbnail:      ep.Thumbnail,
		FormatID:       ep.FormatID,
		Width:          ep.Width,
		Height:         ep.Height,
		DynamicRange:   ep.DynamicRange,
	}
}
