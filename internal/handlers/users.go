package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MazzMS/go-rss/internal/database"
	"github.com/MazzMS/go-rss/internal/models"
	"github.com/MazzMS/go-rss/internal/utils"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	if cfg.Debug {
		log.Println("CALLING USER CREATION")
		log.Println()
		defer log.Println("END OF USER CREATION")
		defer log.Println()
	}

	// types for JSON's input and output
	type parameters struct {
		Name string `json:"name"`
	}

	param := parameters{}

	// decode input
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("failed to decode JSON input:", err))
		return
	}

	// check if name was null
	if param.Name == "" {
		utils.RespondWithError(w, http.StatusInternalServerError, "name cannot be null")
		return
	}

	currentTime := time.Now().UTC()

	// create user
	user, err := cfg.DB.CreateUser(
		r.Context(),
		database.CreateUserParams{
			ID:        uuid.New(),
			Name:      param.Name,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
	)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("failed to create user:", err))
		log.Println(err)
		return
	}

	if cfg.Debug {
		log.Printf("successfully recorded user: %s, with UUID %s at %v", user.Name, user.ID, user.CreatedAt)
	}

	utils.RespondWithJSON(w, http.StatusCreated, user)
}

func (cfg *ApiConfig) GetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	if cfg.Debug {
		log.Println("CALLING USER READING")
		log.Println()
		defer log.Println("END OF USER READING")
		defer log.Println()
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DBUserToUser(user))
}
