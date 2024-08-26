package handlers

import (
	"log"
	"net/http"

	"github.com/MazzMS/go-rss/internal/utils"
)

func (cfg *ApiConfig) Err(w http.ResponseWriter, r *http.Request) {
	if cfg.Debug {
		log.Println("CALLING ERROR")
		log.Println()
		defer log.Println("END OF ERROR")
		defer log.Println()
	}

	msg := "Internal Server Error"

	utils.RespondWithError(w, http.StatusInternalServerError, msg)
}
