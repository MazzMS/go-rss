package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/MazzMS/go-rss/internal/config"
	"github.com/MazzMS/go-rss/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// initialize vars
	var cfg config.ApiConfig
	// wrapper function to pass config
	wrapper := func(handler func(http.ResponseWriter, *http.Request, *config.ApiConfig)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, &cfg)
		}
	}

	// check if debug
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	cfg.Debug = *debug

	// load .env
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	// get port
	port := os.Getenv("PORT")

	if cfg.Debug {
		log.Println("Starting go-rss")
		log.Printf("%s will be used as port", port)
	}

	// === HANDLERS ===
	// config mux
	mux := http.NewServeMux()

	// server status
	mux.HandleFunc("GET /v1/healthz", wrapper(handlers.Healthz))
	mux.HandleFunc("GET /v1/err", wrapper(handlers.Err))

	// server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("Starting server")
	log.Fatal(srv.ListenAndServe())
}
