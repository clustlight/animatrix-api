package controller

import (
	"context"
	"strings"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/ent/season"
	"github.com/clustlight/animatrix-api/ent/series"
	"github.com/clustlight/animatrix-api/internal/types"
	"github.com/clustlight/animatrix-api/internal/utils"
)

// Convert Katakana to Hiragana
func katakanaToHiragana(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 0x30A1 && r <= 0x30F6 {
			return r - 0x60
		}
		return r
	}, s)
}

// Normalize symbols: unify long vowel marks, middle dots, and remove spaces
func normalizeSymbols(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case 'ー', '-', 'ｰ':
			b.WriteRune('ー') // unify to full-width long vowel mark
		case '・', '･':
			b.WriteRune('・') // unify to full-width middle dot
		case ' ', '　', '\t':
			// remove all types of spaces
		default:
			// add other symbols if needed
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Comprehensive normalization: lowercase, normalize symbols, convert Katakana to Hiragana
func normalizeJapanese(s string) string {
	s = strings.ToLower(s)
	s = normalizeSymbols(s)
	s = katakanaToHiragana(s)
	return s
}

// Search series by title or yomi, and also by related seasons
func SearchSeries(ctx context.Context, client *ent.Client, query string) ([]types.SeriesResponse, error) {
	queryHira := normalizeJapanese(query)

	// Search by series title and yomi
	seriesList, err := client.Series.Query().
		Where(
			series.Or(
				series.TitleContainsFold(query),
				series.TitleYomiContainsFold(query),
				series.TitleYomiContainsFold(queryHira),
				series.TitleEnContainsFold(query),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	seriesMap := make(map[int]*ent.Series)
	for _, s := range seriesList {
		seriesMap[s.ID] = s
	}

	// Search by season title and yomi, add related series
	seasons, err := client.Season.Query().
		Where(
			season.Or(
				season.SeasonTitleContainsFold(query),
				season.SeasonTitleYomiContainsFold(query),
				season.SeasonTitleYomiContainsFold(queryHira),
			),
		).
		WithSeries().
		All(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range seasons {
		if s.Edges.Series != nil {
			seriesMap[s.Edges.Series.ID] = s.Edges.Series
		}
	}

	// Remove duplicates and build response
	seriesRes := make([]types.SeriesResponse, 0, len(seriesMap))
	for _, s := range seriesMap {
		seriesRes = append(seriesRes, utils.BuildSeriesResponse(s, false, false))
	}

	return seriesRes, nil
}
