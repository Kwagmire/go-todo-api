// @title ToDo List API
// @version 1.0
// @description This is a sample ToDo List API built with Go.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/Kwagmire/go-todo-api/internal/app/handlers"
	"github.com/Kwagmire/go-todo-api/internal/pkg/auth"
	"github.com/Kwagmire/go-todo-api/internal/pkg/db"

	_ "github.com/Kwagmire/go-todo-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file. Assuming environment variables are set in the environment.")
	}

	db.InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	mux.HandleFunc("POST /register", handlers.RegisterUser)
	mux.HandleFunc("POST /login", handlers.LoginUser)

	mux.HandleFunc("POST /todos", auth.AuthMiddleware(handlers.AddTodo))
	mux.HandleFunc("GET /todos", auth.AuthMiddleware(handlers.GetTodos))

	mux.HandleFunc("PUT /todos/", auth.AuthMiddleware(handlers.UpdateTodo))
	mux.HandleFunc("DELETE /todos/", auth.AuthMiddleware(handlers.DeleteTodo))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)
	serverPort := ":8080"

	fmt.Printf("Todo API server starting on http://localhost%s...", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, handler))
}
