package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"workspace/github.com/Benjysparks/daily/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Hours     int32     `json:"Hours"`
	Minutes   int32     `json:"Minutes"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Hours    int32  `json:"hourNumber"`
		Minutes  int32  `json:"minuteNumber"`
	}

	decoder := json.NewDecoder(r.Body) // Fix: use r.Body instead of r.Email
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	fmt.Println(params.Hours)
	fmt.Println(params.Minutes)

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:       params.Email,
		Pword:       params.Password,
		FullName:    params.Name,
		UserHours:   params.Hours,
		UserMinutes: params.Minutes,
	})

	err = cfg.scheduleUserEmail(params.Email, params.Hours, params.Minutes)
	if err != nil {
		log.Printf("Failed to schedule email for %s: %v", user.Email, err)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// func (cfg *apiConfig) scheduleUserEmail(email string, hour, minute int) error {
	// // Validate hour and minute
	// if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
	// 	return fmt.Errorf("invalid hour or minute")
	// }

	// // Format cron spec: "30 8 * * *"
	// spec := fmt.Sprintf("%d %d * * *", minute, hour)

	// // Add function to cron
	// return cfg.cron.AddFunc(spec, func() {
	// 	cfg.emailerPrefs(email)
	// })
	// }

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Username:  user.FullName,
	})
}

func (cfg *apiConfig) handlerClearUsers(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.ClearUserTable(r.Context())
	if err != nil {
		fmt.Println("error")
	}
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
			Hours:     user.UserHours,
			Minutes:   user.UserMinutes,
		}
	}

	respondWithJSON(w, http.StatusOK, responseUsers)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	// Look up user in database
	user, err := cfg.db.GetUserByEmail(r.Context(), creds.Email)
	if err != nil || user.Pword != creds.Password {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Create token with user_id and expiry
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(tokenExpiry).Unix(),
	})

	// Sign the token using []byte for jwtSecret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate token", err)
		return
	}

	err = cfg.db.SaveUserToken(r.Context(), database.SaveUserTokenParams{
		ID:       user.ID,
		JwtToken: tokenString,
	})
	if err != nil {
		log.Println("SaveUserToken DB error:", err)
		respondWithError(w, http.StatusBadRequest, "Could not save token to database", err)
		return
	}

	// Return the token as JSON
	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Ensure jwtSecret is of type []byte
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract user_id from claims and store in context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store the user_id in the request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next(w, r.WithContext(ctx))
	}
}

func getUserIDFromContextOrToken(r *http.Request) (uuid.UUID, error) {

	if userIDVal := r.Context().Value("userID"); userIDVal != nil {
		if userID, ok := userIDVal.(uuid.UUID); ok {
			return userID, nil
		}
	}

	// If not in context, parse Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return uuid.UUID{}, errors.New("missing Authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return uuid.UUID{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.UUID{}, errors.New("user_id not a string")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid user_id UUID format")
	}
	return userID, nil

}

func (cfg *apiConfig) handlerSendInfoToFront(w http.ResponseWriter, r *http.Request) {
	var user database.GetUserByIDRow
	userId, err := getUserIDFromContextOrToken(r)
	if err != nil {
		fmt.Println(err)
	}
	user, err = cfg.db.GetUserByID(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	respondWithJSON(w, http.StatusOK, database.User{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	})
}

func (cfg *apiConfig) handlerUpdatePreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Preferences []string `json:"preferences"`
		ExtraData   []string `json:"extraData"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContextOrToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	prefsJSON, err := json.Marshal(req.Preferences)
	if err != nil {
		http.Error(w, "Failed to marshal preferences", http.StatusInternalServerError)
		return
	}

	extraDataJSON, err := json.Marshal(req.ExtraData)
	if err != nil {
		http.Error(w, "Failed to marshal extra data", http.StatusInternalServerError)
		return
	}

	err = cfg.db.SavePreferences(r.Context(), database.SavePreferencesParams{
		UserID:              userID,
		Preferences:         prefsJSON,
		PreferenceVariables: extraDataJSON,
	})
	if err != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getUserPreferences(db *sql.DB, userID int) ([]string, []string, error) {
	var prefsJSON string
	var extraDataJSON string
	query := `SELECT preferences, preference_variables FROM user_preferences WHERE user_id = $1`

	err := db.QueryRow(query, userID).Scan(&prefsJSON, &extraDataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return []string{}, []string{}, nil
		}
		return nil, nil, err
	}

	var preferences []string
	if err := json.Unmarshal([]byte(prefsJSON), &preferences); err != nil {
		return nil, nil, err
	}

	var extraData []string
	if err := json.Unmarshal([]byte(extraDataJSON), &extraData); err != nil {
		return nil, nil, err
	}

	return preferences, extraData, nil
}

func (cfg *apiConfig) handlerShowUserPreferences(w http.ResponseWriter, r *http.Request) {
	userPrefs := []database.ShowAllUserPreferencesRow{}
	users, _ := cfg.db.ShowAllUserPreferences(r.Context())
	for _, user := range users {
		userPrefs = append(userPrefs, database.ShowAllUserPreferencesRow{
			ID:          user.ID,
			Email:       user.Email,
			UserHours:   user.UserHours,
			UserMinutes: user.UserMinutes,
			Preferences: user.Preferences,
		})
	}
	respondWithJSON(w, http.StatusOK, userPrefs)
}

func (cfg *apiConfig) handlerGetUserByToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error": "Missing Authorization header"}`, http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, `{"error": "Invalid Authorization header format"}`, http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	user, err := cfg.db.GetUserByToken(r.Context(), tokenString)
	if err != nil {
		http.Error(w, `{"error": "Error getting user information"}`, http.StatusInternalServerError)
		fmt.Println("DB error:", err)
		return
	}

	response := struct {
		FullName string `json:"fullName"`
	}{
		FullName: user.FullName,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}
