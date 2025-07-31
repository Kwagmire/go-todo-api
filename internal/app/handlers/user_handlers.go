package handlers

import (
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
	var UserID int
	err = db.DB.QueryRow(query, thisRequest.Email, thisRequest.Name, string(hashedPassword)).Scan(&UserID)
	if err != nil {
		if dbError, ok := err.(*pq.Error); ok && dbError.Code.Name() == "unique_violation" {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateJWT(UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"token": token})
}
