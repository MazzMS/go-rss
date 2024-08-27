package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MazzMS/go-rss/internal/database"
	"github.com/MazzMS/go-rss/internal/models"
	"github.com/MazzMS/go-rss/internal/utils"
)

func (cfg *ApiConfig) GetPosts(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	if cfg.Debug {
		log.Println("CALLING POST RECOVER")
		log.Println()
		defer log.Println("END OF POST RECOVER")
		defer log.Println()
	}

	// check if limit
	limitString := r.URL.Query().Get("limit")
	limit := 100
	if limitString != "" {
		n, err := strconv.Atoi(limitString)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprint("failed to parse string to int:", err))
			return
		}
		limit = n
	}

	// get posts
	dbPosts, err := cfg.DB.GetPostsByUser(
		r.Context(),
		database.GetPostsByUserParams{
			ID:    dbUser.ID,
			Limit: int32(limit),
		},
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("failed to fetch posts from database:", err))
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DBPostsToPosts(dbPosts))
}
