package types

import (
	"fmt"
	"time"
)

type EpisodeResponse struct {
	EpisodeID      string    `json:"episode_id"`
	Title          string    `json:"title"`
	EpisodeNumber  int       `json:"episode_number"`
	Duration       float64   `json:"duration"`
	DurationString string    `json:"duration_string"`
	Timestamp      time.Time `json:"timestamp"`
	FormatID       string    `json:"format_id"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
	DynamicRange   string    `json:"dynamic_range"`
	VideoURL       string    `json:"video_url"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	Description    string    `json:"description"`
}

type CreateEpisodeRequest struct {
	SeasonID       string    `json:"season_id" validate:"required"`
	EpisodeID      string    `json:"episode_id" validate:"required"`
	Title          string    `json:"title" validate:"required"`
	EpisodeNumber  int       `json:"episode_number" validate:"required"`
	Duration       float64   `json:"duration" validate:"required"`
	DurationString string    `json:"duration_string" validate:"required"`
	Timestamp      time.Time `json:"timestamp" validate:"required"` // ISO 8601 format
	FormatID       string    `json:"format_id" validate:"required"`
	Width          int       `json:"width" validate:"required"`
	Height         int       `json:"height" validate:"required"`
	DynamicRange   string    `json:"dynamic_range" validate:"required"`
	Metadata       string    `json:"metadata,omitempty"` // Optional field for additional metadata
	Description    string    `json:"description,omitempty"`
}

type UpdateEpisodeRequest struct {
	Title          *string    `json:"title,omitempty"`
	EpisodeNumber  *int       `json:"episode_number,omitempty"`
	Duration       *float64   `json:"duration,omitempty"`
	DurationString *string    `json:"duration_string,omitempty"`
	Timestamp      *time.Time `json:"timestamp,omitempty"` // ISO 8601 format`
	FormatID       *string    `json:"format_id,omitempty"`
	Width          *int       `json:"width,omitempty"`
	Height         *int       `json:"height,omitempty"`
	DynamicRange   *string    `json:"dynamic_range,omitempty"`
	Metadata       *string    `json:"metadata,omitempty"`
	Description    *string    `json:"description,omitempty"`
}

func (r *CreateEpisodeRequest) ValidateRequired() error {
	if r.SeasonID == "" {
		return fmt.Errorf("season_id is required")
	}
	if r.EpisodeID == "" {
		return fmt.Errorf("episode_id is required")
	}
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if r.EpisodeNumber < 0 {
		return fmt.Errorf("episode_number must be greater than -1")
	}
	if r.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	if r.DurationString == "" {
		return fmt.Errorf("duration_string is required")
	}
	if r.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}
	if r.FormatID == "" {
		return fmt.Errorf("format_id is required")
	}
	if r.Width <= 0 {
		return fmt.Errorf("width must be greater than 0")
	}
	if r.Height <= 0 {
		return fmt.Errorf("height must be greater than 0")
	}
	if r.DynamicRange == "" {
		return fmt.Errorf("dynamic_range is required")
	}
	return nil
}
