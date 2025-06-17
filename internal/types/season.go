package types

import (
	"fmt"
)

type SeasonResponse struct {
	SeasonID        string            `json:"season_id"`
	SeasonTitle     string            `json:"season_title"`
	SeasonTitleYomi string            `json:"season_title_yomi"`
	SeasonNumber    int               `json:"season_number"`
	ShoboiTID       int               `json:"shoboi_tid"`
	Description     string            `json:"description"`
	FirstYear       int               `json:"first_year"`
	FirstMonth      int               `json:"first_month"`
	FirstEndYear    int               `json:"first_end_year"`
	FirstEndMonth   int               `json:"first_end_month"`
	ThumbnailURL    string            `json:"thumbnail_url"`
	Episodes        []EpisodeResponse `json:"episodes,omitempty"`
}

type CreateSeasonRequest struct {
	SeriesID        string  `json:"series_id" validate:"required"`
	SeasonID        string  `json:"season_id" validate:"required"`
	SeasonTitle     string  `json:"season_title" validate:"required"`
	SeasonTitleYomi *string `json:"season_title_yomi,omitempty"`
	SeasonNumber    int     `json:"season_number" validate:"required"`
	ShoboiTID       *int    `json:"shoboi_tid,omitempty"`
	Description     *string `json:"description,omitempty"`
	FirstYear       *int    `json:"first_year,omitempty"`
	FirstMonth      *int    `json:"first_month,omitempty"`
	FirstEndYear    *int    `json:"first_end_year,omitempty"`
	FirstEndMonth   *int    `json:"first_end_month,omitempty"`
}

type UpdateSeasonRequest struct {
	SeasonTitle     *string `json:"season_title,omitempty"`
	SeasonTitleYomi *string `json:"season_title_yomi,omitempty"`
	SeasonNumber    *int    `json:"season_number,omitempty"`
	ShoboiTID       *int    `json:"shoboi_tid,omitempty"`
	Description     *string `json:"description,omitempty"`
	FirstYear       *int    `json:"first_year,omitempty"`
	FirstMonth      *int    `json:"first_month,omitempty"`
	FirstEndYear    *int    `json:"first_end_year,omitempty"`
	FirstEndMonth   *int    `json:"first_end_month,omitempty"`
}

func (r *CreateSeasonRequest) ValidateRequired() error {
	if r.SeriesID == "" {
		return fmt.Errorf("series_id")
	}
	if r.SeasonID == "" {
		return fmt.Errorf("season_id")
	}
	if r.SeasonTitle == "" {
		return fmt.Errorf("season_title")
	}
	if r.SeasonNumber < 0 {
		return fmt.Errorf("season_number must be greater than -1")
	}
	return nil
}
