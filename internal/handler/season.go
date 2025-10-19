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

		if err := seasonData.ValidateRequired(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

func BulkCreateSeasonHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var seasonList []types.CreateSeasonRequest
		if err := json.NewDecoder(r.Body).Decode(&seasonList); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		newSeasonList, err := controller.BulkCreateSeason(r.Context(), client, seasonList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newSeasonList)
	}
}

func DeleteSeason(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seasonID := chi.URLParam(r, "season_id")
		if seasonID == "" {
			http.Error(w, "season_id required", http.StatusBadRequest)
			return
		}
		if err := controller.DeleteSeason(r.Context(), client, seasonID); err != nil {
			if err == controller.ErrHasChildren {
				http.Error(w, "Season has episodes; cannot delete", http.StatusConflict)
			} else if _, ok := err.(*ent.NotFoundError); ok {
				http.Error(w, "Season not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
