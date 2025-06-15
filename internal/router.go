package internal

import (
	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/internal/handler"

	"github.com/go-chi/chi/v5"
)

func NewRouter(client *ent.Client) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v1", func(api chi.Router) {
		api.Get("/series", handler.GetAllSeries(client))
		api.Post("/series", handler.CreateSeries(client))
		api.Get("/series/{series_id}", handler.GetSeriesDetail(client))
		api.Patch("/series/{series_id}", handler.UpdateSeries(client))

		api.Post("/series/bulk", handler.BulkCreateSeriesHandler(client))

		api.Get("/season", handler.GetAllSeasons(client))
		api.Post("/season", handler.CreateSeason(client))
		api.Get("/season/{season_id}", handler.GetSeasonDetail(client))
		api.Patch("/season/{season_id}", handler.UpdateSeason(client))

		api.Post("/season/bulk", handler.BulkCreateSeasonHandler(client))

		api.Get("/episode", handler.GetAllEpisodes(client))
		api.Post("/episode", handler.CreateEpisode(client))
		api.Get("/episode/{episode_id}", handler.GetEpisodeDetail(client))
		api.Patch("/episode/{episode_id}", handler.UpdateEpisode(client))

		api.Post("/episode/bulk", handler.BulkCreateEpisodeHandler(client))

		api.Get("/content/{episode_id}/{filename}", handler.GetContent())
	})
	return r
}
