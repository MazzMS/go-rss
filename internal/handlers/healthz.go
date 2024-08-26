package handlers

import (
	"log"
	"net/http"

	"github.com/MazzMS/go-rss/internal/utils"
)

func (cfg *ApiConfig) Healhtz(w http.ResponseWriter, r *http.Request) {
	if cfg.Debug {
		log.Println("CALLING HEALTHZ")
		log.Println()
		defer log.Println("END OF HEALTHZ")
		defer log.Println()
	}

	type response struct {
		Status string `json:"status"`
	}

	res := response{}
	res.Status = "ok"

	utils.RespondWithJSON(w, http.StatusOK, res)
}
