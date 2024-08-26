package handlers

import (
	"log"
	"net/http"

	"github.com/MazzMS/go-rss/internal/config"
	"github.com/MazzMS/go-rss/internal/utils"
)

func Err(w http.ResponseWriter, r *http.Request, config *config.ApiConfig) {
	if config.Debug {
		log.Println("CALLING ERROR")
		log.Println()
		defer log.Println("END OF ERROR")
		defer log.Println()
	}

	msg := "Internal Server Error"

	utils.RespondWithError(w, http.StatusInternalServerError, msg)
}
