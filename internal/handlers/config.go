package handlers

import "github.com/MazzMS/go-rss/internal/database"

type ApiConfig struct {
	DB *database.Queries
	Debug bool
}
