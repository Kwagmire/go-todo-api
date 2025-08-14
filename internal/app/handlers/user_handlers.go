package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Kwagmire/go-todo-api/internal/pkg/auth"
	"github.com/Kwagmire/go-todo-api/internal/pkg/db"
	"github.com/Kwagmire/go-todo-api/internal/pkg/models"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Creates a new user with the provided credentials
// @Security ApiKeyAuth
// @Accept  json
// @Produce json,plain
// @Param   user  body  models.RegisterRequest  true  "Credentials for new user"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Invalid request payload"
// @Failure 403 {string} string "User exists"
// @Router /register [post]
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

	token, err := auth.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"token": token})
}

// @Summary Log a user in
// @Description Authenticate a user with provided credentials
// @Security ApiKeyAuth
// @Accept  json
// @Produce json,plain
// @Param   user  body  models.LoginRequest  true  "User login credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
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
	err = db.DB.QueryRow(query, thisRequest.Email).Scan(&userID, &hashedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(thisRequest.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
