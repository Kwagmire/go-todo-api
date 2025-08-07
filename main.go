package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/Kwagmire/go-todo-api/internal/app/handlers"
	"github.com/Kwagmire/go-todo-api/internal/pkg/auth"
	"github.com/Kwagmire/go-todo-api/internal/pkg/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file. Assuming environment variables are set in the environment.")
	}

	db.InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", handlers.RegisterUser)
	mux.HandleFunc("POST /login", handlers.LoginUser)

	mux.HandleFunc("POST /todos", auth.AuthMiddleware(handlers.AddTodo))
	mux.HandleFunc("GET /todos", auth.AuthMiddleware(handlers.GetTodos))

	mux.HandleFunc("PUT /todos/", auth.AuthMiddleware(handlers.UpdateTodo))
	mux.HandleFunc("DELETE /todos/", auth.AuthMiddleware(handlers.DeleteTodo))

	serverPort := ":8080"

	fmt.Printf("Todo API server starting on http://localhost%s...", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, mux))
}
