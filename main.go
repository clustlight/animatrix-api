package main

import (
	"log"
	"net/http"

	"github.com/clustlight/animatrix-api/internal"
	"github.com/clustlight/animatrix-api/internal/utils"
)

func main() {
	client := utils.NewDBClient()
	defer client.Close()
	log.Println("server started at :8080")
	http.ListenAndServe(":8080", internal.NewRouter(client))
}
