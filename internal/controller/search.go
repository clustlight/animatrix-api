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

// Convert Katakana to Hiragana
func katakanaToHiragana(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 0x30A1 && r <= 0x30F6 {
			return r - 0x60
		}
		return r
	}, s)
}

func hiraganaToKatakana(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 0x3041 && r <= 0x3096 {
			return r + 0x60
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

// Search series by title or yomi, and also by related seasons, using kagome tokens for DB search
func SearchSeries(ctx context.Context, client *ent.Client, query string) ([]types.SeriesResponse, error) {
	queryHira := normalizeJapanese(query)
	queryKana := hiraganaToKatakana(query)
	tokens := tokenizeJapanese(query)
	tokenSet := make(map[string]struct{})
	for _, token := range tokens {
		tokenSet[token] = struct{}{}
		tokenSet[normalizeJapanese(token)] = struct{}{}
		tokenSet[hiraganaToKatakana(token)] = struct{}{}
	}

	seriesPredicates := []predicate.Series{
		series.TitleContainsFold(query),
		series.TitleYomiContainsFold(query),
		series.TitleYomiContainsFold(queryHira),
		series.TitleYomiContainsFold(queryKana),
		series.TitleEnContainsFold(query),
	}
	for token := range tokenSet {
		if token == "" {
			continue
		}
		seriesPredicates = append(seriesPredicates,
			series.TitleContainsFold(token),
			series.TitleYomiContainsFold(token),
			series.TitleYomiContainsFold(normalizeJapanese(token)),
			series.TitleYomiContainsFold(hiraganaToKatakana(token)), // 追加
			series.TitleEnContainsFold(token),
		)
	}

	andSeriesPredicates := []predicate.Series{}
	for token := range tokenSet {
		if token == "" {
			continue
		}
		andSeriesPredicates = append(andSeriesPredicates,
			series.Or(
				series.TitleContainsFold(token),
				series.TitleYomiContainsFold(token),
				series.TitleYomiContainsFold(normalizeJapanese(token)),
				series.TitleYomiContainsFold(hiraganaToKatakana(token)),
				series.TitleEnContainsFold(token),
			),
		)
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
		andSeasonPredicates = append(andSeasonPredicates,
			season.Or(
				season.SeasonTitleContainsFold(token),
				season.SeasonTitleYomiContainsFold(token),
				season.SeasonTitleYomiContainsFold(normalizeJapanese(token)),
				season.SeasonTitleYomiContainsFold(hiraganaToKatakana(token)),
			),
		)
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
