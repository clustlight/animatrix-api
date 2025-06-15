package utils

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/clustlight/animatrix-api/ent"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewDBClient() *ent.Client {
	_ = godotenv.Load()
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := cmp.Or(os.Getenv("DATABASE_PORT"), "5432")
	dbname := os.Getenv("DATABASE_NAME")
	sslmode := cmp.Or(os.Getenv("DATABASE_SSLMODE"), "disable")

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
