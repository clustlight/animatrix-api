package handler

import (
	"encoding/json"
	"net/http"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/internal/controller"
	"github.com/clustlight/animatrix-api/internal/types"

	"github.com/go-chi/chi/v5"
)

func GetAllSeries(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		series, err := controller.GetAllSeries(ctx, client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(series)
	}
}

func GetSeriesDetail(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seriesID := chi.URLParam(r, "series_id")
		ctx := r.Context()

		series, err := controller.GetSeries(ctx, client, seriesID)
		if err != nil {
			switch err.(type) {
			case *ent.NotFoundError:
				http.Error(w, "Series not found", http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(series)
	}
}

func CreateSeries(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var seriesData types.CreateSeriesRequest
		if err := json.NewDecoder(r.Body).Decode(&seriesData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		newSeries, err := controller.CreateSeries(ctx, client, &seriesData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newSeries)
	}
}

func UpdateSeries(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seriesID := chi.URLParam(r, "series_id")
		var seriesData types.UpdateSeriesRequest
		if err := json.NewDecoder(r.Body).Decode(&seriesData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		updatedSeries, err := controller.UpdateSeries(ctx, client, seriesID, &seriesData)
		if err != nil {
			switch err.(type) {
			case *ent.NotFoundError:
				http.Error(w, "Series not found", http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedSeries)
	}
}
