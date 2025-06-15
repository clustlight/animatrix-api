package controller

import (
	"context"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/ent/episode"
	"github.com/clustlight/animatrix-api/ent/season"
	"github.com/clustlight/animatrix-api/internal/types"
	"github.com/clustlight/animatrix-api/internal/utils"
)

func GetAllEpisodes(ctx context.Context, client *ent.Client) (*[]types.EpisodeResponse, error) {
	episodes, err := client.Episode.Query().
		Order(ent.Asc("episode_number")).
		All(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]types.EpisodeResponse, 0, len(episodes))
	for _, e := range episodes {
		resp := utils.BuildEpisodeResponse(e)
		responses = append(responses, resp)
	}

	return &responses, nil
}

func GetEpisode(ctx context.Context, client *ent.Client, episodeID string) (*types.EpisodeResponse, error) {
	episode, err := client.Episode.
		Query().
		Where(episode.EpisodeIDEQ(episodeID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	resp := utils.BuildEpisodeResponse(episode)
	return &resp, nil
}

func CreateEpisode(ctx context.Context, client *ent.Client, req *types.CreateEpisodeRequest) (*types.EpisodeResponse, error) {
	season, err := client.Season.
		Query().
		Where(season.SeasonIDEQ(req.SeasonID)).
		Only(ctx)
	if err != nil {
		return nil, err // Seasonが見つからない場合はエラー
	}

	newEpisode, err := client.Episode.Create().
		SetEpisodeID(req.EpisodeID).
		SetTitle(req.Title).
		SetEpisodeNumber(req.EpisodeNumber).
		SetDuration(req.Duration).
		SetDurationString(req.DurationString).
		SetTimestamp(req.Timestamp).
		SetFormatID(req.FormatID).
		SetWidth(req.Width).
		SetHeight(req.Height).
		SetDynamicRange(req.DynamicRange).
		SetMetadata(req.Metadata).
		SetSeason(season).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	resp := utils.BuildEpisodeResponse(newEpisode)
	return &resp, nil
}

func UpdateEpisode(ctx context.Context, client *ent.Client, episodeID string, req *types.UpdateEpisodeRequest) (*types.EpisodeResponse, error) {
	episodeObj, err := client.Episode.
		Query().
		Where(episode.EpisodeIDEQ(episodeID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	update := episodeObj.Update()

	if req.Title != nil {
		update.SetTitle(*req.Title)
	}
	if req.EpisodeNumber != nil {
		update.SetEpisodeNumber(*req.EpisodeNumber)
	}
	if req.Duration != nil {
		update.SetDuration(*req.Duration)
	}
	if req.DurationString != nil {
		update.SetDurationString(*req.DurationString)
	}
	if req.Timestamp != nil {
		update.SetTimestamp(*req.Timestamp)
	}
	if req.FormatID != nil {
		update.SetFormatID(*req.FormatID)
	}
	if req.Width != nil {
		update.SetWidth(*req.Width)
	}
	if req.Height != nil {
		update.SetHeight(*req.Height)
	}
	if req.DynamicRange != nil {
		update.SetDynamicRange(*req.DynamicRange)
	}
	if req.Metadata != nil {
		update.SetMetadata(*req.Metadata)
	}

	updatedEpisode, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	resp := utils.BuildEpisodeResponse(updatedEpisode)
	return &resp, nil
}
