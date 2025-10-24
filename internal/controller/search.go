package controller

import (
	"context"
	"strings"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/ent/predicate"
	"github.com/clustlight/animatrix-api/ent/season"
	"github.com/clustlight/animatrix-api/ent/series"
	"github.com/clustlight/animatrix-api/internal/types"
	"github.com/clustlight/animatrix-api/internal/utils"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

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
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Comprehensive normalization: lowercase, normalize symbols, convert Katakana to Hiragana
func normalizeJapanese(s string) string {
	s = strings.ToLower(s)
	s = normalizeSymbols(s)
	return s
}

// Tokenize Japanese text using kagome
func tokenizeJapanese(text string) []string {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}
	tokens := t.Tokenize(text)
	words := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}
		surface := token.Surface
		if surface != "" {
			words = append(words, surface)
		}
	}
	return words
}

func hiraToKata(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 0x3041 && r <= 0x3096 {
			b.WriteRune(r + 0x60)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func kataToHira(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 0x30A1 && r <= 0x30F6 {
			b.WriteRune(r - 0x60)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func isHiragana(s string) bool {
	for _, r := range s {
		if r < 0x3041 || r > 0x3096 {
			return false
		}
	}
	return len(s) > 0
}

func isKatakana(s string) bool {
	for _, r := range s {
		if r < 0x30A1 || r > 0x30F6 {
			return false
		}
	}
	return len(s) > 0
}

// Search series by title or yomi, and also by related seasons, using kagome tokens for DB search
func SearchSeries(ctx context.Context, client *ent.Client, query string) ([]types.SeriesResponse, error) {
	queryHira := normalizeJapanese(query)
	tokens := tokenizeJapanese(query)
	tokenSet := make(map[string]struct{})
	for _, token := range tokens {
		tokenSet[token] = struct{}{}
		tokenSet[normalizeJapanese(token)] = struct{}{}
	}

	var hiraOrKata string
	if isHiragana(queryHira) {
		hiraOrKata = hiraToKata(queryHira)
	} else if isKatakana(queryHira) {
		hiraOrKata = kataToHira(queryHira)
	}

	seriesPredicates := []predicate.Series{
		series.TitleContainsFold(query),
		series.TitleYomiContainsFold(query),
		series.TitleYomiContainsFold(queryHira),
		series.TitleEnContainsFold(query),
	}
	if hiraOrKata != "" {
		seriesPredicates = append(seriesPredicates,
			series.TitleContainsFold(hiraOrKata),
			series.TitleYomiContainsFold(hiraOrKata),
		)
	}
	for token := range tokenSet {
		if token == "" {
			continue
		}
		seriesPredicates = append(seriesPredicates,
			series.TitleContainsFold(token),
			series.TitleYomiContainsFold(token),
			series.TitleYomiContainsFold(normalizeJapanese(token)),
			series.TitleEnContainsFold(token),
		)

		if isHiragana(token) {
			kata := hiraToKata(token)
			seriesPredicates = append(seriesPredicates,
				series.TitleContainsFold(kata),
				series.TitleYomiContainsFold(kata),
			)
		} else if isKatakana(token) {
			hira := kataToHira(token)
			seriesPredicates = append(seriesPredicates,
				series.TitleContainsFold(hira),
				series.TitleYomiContainsFold(hira),
			)
		}
	}

	// Build AND condition for all tokens
	andSeriesPredicates := []predicate.Series{}
	for token := range tokenSet {
		if token == "" {
			continue
		}
		orPreds := []predicate.Series{
			series.TitleContainsFold(token),
			series.TitleYomiContainsFold(token),
			series.TitleYomiContainsFold(normalizeJapanese(token)),
			series.TitleEnContainsFold(token),
		}

		if isHiragana(token) {
			kata := hiraToKata(token)
			orPreds = append(orPreds,
				series.TitleContainsFold(kata),
				series.TitleYomiContainsFold(kata),
			)
		} else if isKatakana(token) {
			hira := kataToHira(token)
			orPreds = append(orPreds,
				series.TitleContainsFold(hira),
				series.TitleYomiContainsFold(hira),
			)
		}
		andSeriesPredicates = append(andSeriesPredicates, series.Or(orPreds...))
	}

	seriesList, err := client.Series.Query().
		Where(
			series.And(andSeriesPredicates...),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	seriesMap := make(map[int]*ent.Series)
	for _, s := range seriesList {
		seriesMap[s.ID] = s
	}

	// Similarly, build AND condition for seasons
	andSeasonPredicates := []predicate.Season{}
	for token := range tokenSet {
		if token == "" {
			continue
		}
		orPreds := []predicate.Season{
			season.SeasonTitleContainsFold(token),
			season.SeasonTitleYomiContainsFold(token),
			season.SeasonTitleYomiContainsFold(normalizeJapanese(token)),
		}

		if isHiragana(token) {
			kata := hiraToKata(token)
			orPreds = append(orPreds,
				season.SeasonTitleContainsFold(kata),
				season.SeasonTitleYomiContainsFold(kata),
			)
		} else if isKatakana(token) {
			hira := kataToHira(token)
			orPreds = append(orPreds,
				season.SeasonTitleContainsFold(hira),
				season.SeasonTitleYomiContainsFold(hira),
			)
		}
		andSeasonPredicates = append(andSeasonPredicates, season.Or(orPreds...))
	}

	seasons, err := client.Season.Query().
		Where(
			season.And(andSeasonPredicates...),
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

	seriesRes := make([]types.SeriesResponse, 0, len(seriesMap))
	for _, s := range seriesMap {
		seriesRes = append(seriesRes, utils.BuildSeriesResponse(s, false, false))
	}

	return seriesRes, nil
}
