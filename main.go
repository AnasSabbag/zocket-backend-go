package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/AnasSabbag/task-manager/handlers"
	"github.com/AnasSabbag/task-manager/middleware"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.Handle("/tasks", middleware.JWTMiddleware(http.HandlerFunc(handlers.GetTasks))).Methods("GET")
	router.Handle("/tasks", middleware.JWTMiddleware(http.HandlerFunc(handlers.CreateTask))).Methods("POST")
	router.Handle("/tasks/{id}", middleware.JWTMiddleware(http.HandlerFunc(handlers.UpdateTask))).Methods("PUT")
	router.Handle("/tasks/{id}", middleware.JWTMiddleware(http.HandlerFunc(handlers.DeleteTask))).Methods("DELETE")



	// WebSocket Route
	router.HandleFunc("/ws", handlers.WebSocketHandler)
	

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Start WebSocket Listener
	go handlers.StartWebSocketListener()

	log.Println("Server started on :5000")
	http.ListenAndServe(":5000", c.Handler(router))
}
