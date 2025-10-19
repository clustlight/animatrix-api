package handler

import (
	"encoding/json"
	"net/http"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/internal/controller"
	"github.com/clustlight/animatrix-api/internal/types"

	"github.com/go-chi/chi/v5"
)

func GetAllEpisodes(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		episodes, err := controller.GetAllEpisodes(ctx, client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(episodes)
	}
}

func GetEpisodeDetail(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		episodeID := chi.URLParam(r, "episode_id")
		ctx := r.Context()

		episode, err := controller.GetEpisode(ctx, client, episodeID)
		if err != nil {
			switch err.(type) {
			case *ent.NotFoundError:
				http.Error(w, "Episode not found", http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(episode)
	}
}

func CreateEpisode(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var episodeData types.CreateEpisodeRequest
		if err := json.NewDecoder(r.Body).Decode(&episodeData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := episodeData.ValidateRequired(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		newEpisode, err := controller.CreateEpisode(ctx, client, &episodeData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newEpisode)
	}
}

func UpdateEpisode(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		episodeID := chi.URLParam(r, "episode_id")
		var episodeData types.UpdateEpisodeRequest
		if err := json.NewDecoder(r.Body).Decode(&episodeData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		updatedEpisode, err := controller.UpdateEpisode(ctx, client, episodeID, &episodeData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedEpisode)
	}
}

func BulkCreateEpisodeHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var EpisodeList []types.CreateEpisodeRequest
		if err := json.NewDecoder(r.Body).Decode(&EpisodeList); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		newEpisodeList, err := controller.BulkCreateEpisode(r.Context(), client, EpisodeList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newEpisodeList)
	}
}

func DeleteEpisode(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		episodeID := chi.URLParam(r, "episode_id")
		if episodeID == "" {
			http.Error(w, "episode_id required", http.StatusBadRequest)
			return
		}
		if err := controller.DeleteEpisode(r.Context(), client, episodeID); err != nil {
			switch _, ok := err.(*ent.NotFoundError); {
			case ok:
				http.Error(w, "Episode not found", http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
