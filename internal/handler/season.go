package handler

import (
	"encoding/json"
	"net/http"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/internal/controller"
	"github.com/clustlight/animatrix-api/internal/types"
	"github.com/go-chi/chi/v5"
)

func GetAllSeasons(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		seasons, err := controller.GetAllSeasons(ctx, client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(seasons)
	}
}

func GetSeasonDetail(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seasonID := chi.URLParam(r, "season_id")
		ctx := r.Context()

		season, err := controller.GetSeason(ctx, client, seasonID)
		if err != nil {
			switch err.(type) {
			case *ent.NotFoundError:
				http.Error(w, "Season not found", http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(season)
	}
}

func CreateSeason(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var seasonData types.CreateSeasonRequest
		if err := json.NewDecoder(r.Body).Decode(&seasonData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		newSeason, err := controller.CreateSeason(ctx, client, &seasonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newSeason)
	}
}

func UpdateSeason(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seasonID := chi.URLParam(r, "season_id")
		var seasonData types.UpdateSeasonRequest
		if err := json.NewDecoder(r.Body).Decode(&seasonData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		updatedSeason, err := controller.UpdateSeason(ctx, client, seasonID, &seasonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedSeason)
	}
}
