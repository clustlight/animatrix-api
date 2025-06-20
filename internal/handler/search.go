package handler

import (
	"encoding/json"
	"net/http"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/clustlight/animatrix-api/internal/controller"
)

func SearchHandler(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		result, err := controller.SearchSeries(ctx, client, query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
