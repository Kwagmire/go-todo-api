package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"go-todo-api/internal/pkg/auth"
	"go-todo-api/internal/pkg/db"
	"go-todo-api/internal/pkg/models"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unaccepted method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var thisRequest models.RegisterRequest
	err = json.Unmarshal(body, &thisRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if thisRequest.Email == "" || thisRequest.Name == "" || thisRequest.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	if len(thisRequest.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(thisRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	query := `
		INSERT INTO users (
			email,
			name,
			password_hash
		) VALUES ($1, $2, $3
		) RETURNING id`
	var userID int
	err = db.DB.QueryRow(query, thisRequest.Email, thisRequest.Name, string(hashedPassword)).Scan(&userID)
	if err != nil {
		if dbError, ok := err.(*pq.Error); ok && dbError.Code.Name() == "unique_violation" {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateJWT(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unaccepted method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var thisRequest models.LoginRequest
	err = json.Unmarshal(body, &thisRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if thisRequest.Email == "" || thisRequest.Password == "" {
		http.Error(w, "Input all fields to login", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, password_hash
		FROM users
		WHERE email = $1`
	var userID int
	var hashedPassword string
	err := db.DB.QueryRow(query, thisRequest.Email).Scan(&userID, &hashedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if err := bvrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(thisRequest.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
