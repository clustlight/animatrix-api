package handler

import (
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
)

var validEpisodeID = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
var validFilename = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

func GetContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		episodeID := chi.URLParam(r, "episode_id")
		filename := chi.URLParam(r, "filename")

		if !validEpisodeID.MatchString(episodeID) {
			http.Error(w, "invalid episode_id", http.StatusBadRequest)
			return
		}
		if !validFilename.MatchString(filename) || strings.Contains(filename, "..") || strings.Contains(filename, "/") {
			http.Error(w, "invalid filename", http.StatusBadRequest)
			return
		}
		pathStr := strings.Replace(episodeID, "_", "/", 1)

		cleanFilename := path.Clean(filename)
		if cleanFilename != filename {
			http.Error(w, "invalid filename", http.StatusBadRequest)
			return
		}

		objectUrl := os.Getenv("OBJECT_STORAGE_URL") + pathStr + "/" + cleanFilename
		http.Redirect(w, r, objectUrl, http.StatusFound)
	}
}
