# Animatrix API

Animatrix API is a RESTful server for managing ABM archives.

## Features

- CRUD operations for series, seasons, and episodes
- Bulk registration API
- Recently updated series endpoint
- Search endpoint
- Supports PostgreSQL database
- Object storage integration (for thumbnails, etc.)

## Directory Structure

```
.
├── main.go                # Entry point
├── Dockerfile             # Docker build file
├── compose.yaml           # Docker Compose
├── ent/                   # ent ORM definitions
├── internal/              # Routers, utilities, etc.
```

## API Endpoints

### Series
- `GET    /v1/series`                 - List all series
- `POST   /v1/series`                 - Create a new series
- `GET    /v1/series/{series_id}`     - Get a specific series
- `PATCH  /v1/series/{series_id}`     - Update a series
- `POST   /v1/series/bulk`            - Bulk create series
- `GET    /v1/series/recent`          - List recently updated series

### Season
- `GET    /v1/season`                 - List all seasons
- `POST   /v1/season`                 - Create a new season
- `GET    /v1/season/{season_id}`     - Get a specific season
- `PATCH  /v1/season/{season_id}`     - Update a season
- `POST   /v1/season/bulk`            - Bulk create seasons

### Episode
- `GET    /v1/episode`                - List all episodes
- `POST   /v1/episode`                - Create a new episode
- `GET    /v1/episode/{episode_id}`   - Get a specific episode
- `PATCH  /v1/episode/{episode_id}`   - Update an episode
- `POST   /v1/episode/bulk`           - Bulk create episodes

### Search
- `GET    /v1/search`                 - Search series, seasons