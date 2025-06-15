package types

type SeriesResponse struct {
	SeriesID  string           `json:"series_id"`
	Title     string           `json:"title"`
	TitleYomi string           `json:"title_yomi"`
	Seasons   []SeasonResponse `json:"seasons,omitempty"`
}

type CreateSeriesRequest struct {
	SeriesID  string `json:"series_id" validate:"required"`
	Title     string `json:"title" validate:"required"`
	TitleYomi string `json:"title_yomi"`
}

type UpdateSeriesRequest struct {
	Title     *string `json:"title,omitempty"`
	TitleYomi *string `json:"title_yomi,omitempty"`
}
