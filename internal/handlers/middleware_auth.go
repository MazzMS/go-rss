package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MazzMS/go-rss/internal/auth"
	"github.com/MazzMS/go-rss/internal/database"
	"github.com/MazzMS/go-rss/internal/utils"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get apiKey
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %v", err))
			return
		}
		if cfg.Debug {
			log.Printf("ApiKey: %s", apiKey)
		}

		// get user
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("could not find user: %v", err))
			return
		}

		if cfg.Debug {
			log.Printf("find user %v, with uuid %v", user.Name, user.ID)
		}

		handler(w, r, user)
	}
}
