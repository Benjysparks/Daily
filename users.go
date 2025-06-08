package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"workspace/github.com/Benjysparks/daily/internal/database"
)

type User struct {
		ID              uuid.UUID `json:"id"`   
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
		Email           string    `json:"email"`
		Username        string    `json:"username"`

	}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

    type parameters struct {
        Email string `json:"email"`
        Password string `json:"password"`
        Name string `json:"name`
    }   

    decoder := json.NewDecoder(r.Body)  // Fix: use r.Body instead of r.Email
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    params.Email,
		Pword:    params.Password,
		FullName: params.Name,
	})
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
	}

	respondWithJSON(w, http.StatusCreated, User{
			ID:			user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt:	user.UpdatedAt,
			Email:		user.Email,
			Username:	user.FullName,
	})
}

func (cfg *apiConfig) handlerShowAllUser(w http.ResponseWriter, r *http.Request) {
    users, err := cfg.db.GetAllUsers(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve users", err)
        return
    }

    responseUsers := make([]User, len(users))
    for i, user := range users {
        responseUsers[i] = User{
            ID:        user.ID,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
            Email:     user.Email,
            Username:  user.FullName,
        }
    }

    respondWithJSON(w, http.StatusOK, responseUsers)
}
