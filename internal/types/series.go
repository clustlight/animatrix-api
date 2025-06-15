package types

type SeriesResponse struct {
	SeriesID  string           `json:"series_id"`
	Title     string           `json:"title"`
	TitleYomi string           `json:"title_yomi"`
	TitleEn   string           `json:"title_en"`
	Seasons   []SeasonResponse `json:"seasons,omitempty"`
}

type CreateSeriesRequest struct {
	SeriesID  string `json:"series_id" validate:"required"`
	Title     string `json:"title" validate:"required"`
	TitleYomi string `json:"title_yomi,omitempty"`
	TitleEn   string `json:"title_en,omitempty"`
}

type UpdateSeriesRequest struct {
	Title     *string `json:"title,omitempty"`
	TitleYomi *string `json:"title_yomi,omitempty"`
	TitleEn   *string `json:"title_en,omitempty"`
}
