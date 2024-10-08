package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/MazzMS/go-rss/internal/database"
	"github.com/MazzMS/go-rss/internal/handlers"
	"github.com/MazzMS/go-rss/internal/scraper"
	"github.com/joho/godotenv"
)

func main() {
	// === API CONFIG STRUCT === 
	// load .env
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	// check if debug
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// get port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// get db
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// get queries
	dbQueries := database.New(db)

	cfg := handlers.ApiConfig{
		DB: dbQueries,
		Debug: *debug,
	}

	cfg.DB = dbQueries


	if cfg.Debug {
		log.Println("Starting go-rss")
		log.Printf("%s will be used as port", port)
	}

	// start scrapping
	go scraper.StartScraping(cfg.DB, 10, time.Minute)

	// === HANDLERS ===

	// config mux
	mux := http.NewServeMux()

	// server status
	mux.HandleFunc("GET /v1/healthz", cfg.Healhtz)
	mux.HandleFunc("GET /v1/err", cfg.Err)
	// users
	mux.HandleFunc("POST /v1/users", cfg.CreateUser)
	mux.HandleFunc("GET /v1/users", cfg.MiddlewareAuth(cfg.GetUser))
	// feeds
	mux.HandleFunc("POST /v1/feeds", cfg.MiddlewareAuth(cfg.CreateFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.GetAllFeeds)
	// feeds follows
	mux.HandleFunc("POST /v1/feed_follows", cfg.MiddlewareAuth(cfg.CreateFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.MiddlewareAuth(cfg.DeleteFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", cfg.MiddlewareAuth(cfg.GetFollowsFeeds))
	// posts
	mux.HandleFunc("GET /v1/posts", cfg.MiddlewareAuth(cfg.GetPosts))

	// server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("Starting server")
	log.Fatal(srv.ListenAndServe())
}
