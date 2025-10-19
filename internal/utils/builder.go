package utils

import (
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/internal/types"
)

func JoinURL(base, relPath string) string {
	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	ref, err := url.Parse(relPath)
	if err != nil {
		return ""
	}
	return u.ResolveReference(ref).String()
}

func getObjectStorageURL() string {
	return os.Getenv("OBJECT_STORAGE_URL")
}

func buildThumbnailURL(baseURL, id, suffix, ext string) string {
	if suffix != "" {
		return JoinURL(baseURL, id+"/thumbnail_"+suffix+"."+ext)
	}
	return JoinURL(baseURL, id+"/thumbnail."+ext)
}

func BuildSeriesResponse(series *ent.Series, withSeasons, withEpisodes bool) types.SeriesResponse {
	baseURL := getObjectStorageURL()
	resp := types.SeriesResponse{
		SeriesID:     series.SeriesID,
		Title:        series.Title,
		TitleYomi:    series.TitleYomi,
		TitleEn:      series.TitleEn,
		ThumbnailURL: JoinURL(baseURL, series.SeriesID+"/thumbnail.png"),
		PortraitURL:  JoinURL(baseURL, series.SeriesID+"/portrait.png"),
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

func extractSeriesIDAndSuffix(seasonID string) (seriesID, suffix string) {
	idx := strings.Index(seasonID, "_")
	if idx == -1 {
		return seasonID, ""
	}
	seriesID = seasonID[:idx]
	rest := seasonID[idx+1:]
	re := regexp.MustCompile(`^s\d+`)
	suffix = re.FindString(rest)
	return
}

func BuildSeasonResponse(season *ent.Season, withEpisodes bool) types.SeasonResponse {
	baseURL := getObjectStorageURL()
	seriesID, suffix := extractSeriesIDAndSuffix(season.SeasonID)

	thumbURL := ""
	if suffix != "" {
		thumbURL = buildThumbnailURL(baseURL, seriesID, suffix, "png")
	}

	resp := types.SeasonResponse{
		SeriesID:        season.Edges.Series.SeriesID,
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
		ThumbnailURL:    thumbURL,
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
	baseURL := getObjectStorageURL()
	epPath := strings.Replace(ep.EpisodeID, "_", "/", 1)
	VideoUrl := JoinURL(baseURL, epPath+"/video.mp4")
	ThumbnailUrl := JoinURL(baseURL, epPath+"/thumbnail.png")

	return types.EpisodeResponse{
		Title:          ep.Title,
		EpisodeID:      ep.EpisodeID,
		EpisodeNumber:  ep.EpisodeNumber,
		Duration:       ep.Duration,
		DurationString: ep.DurationString,
		Timestamp:      ep.Timestamp,
		FormatID:       ep.FormatID,
		Width:          ep.Width,
		Height:         ep.Height,
		DynamicRange:   ep.DynamicRange,
		VideoURL:       VideoUrl,
		ThumbnailURL:   ThumbnailUrl,
	}
}
