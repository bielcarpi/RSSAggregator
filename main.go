package main

import (
	"database/sql"
	"fmt"
	"github.com/bielcarpi/RSSAggregator/internal/db"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the port and dbUrl from the environment
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("$DB_URL must be set")
	}

	// Connect to the database
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new apiConfig
	apiConfig := apiConfig{
		DB: db.New(conn),
	}

	fmt.Println("Server starting on port " + port)
	router := chi.NewRouter()

	// Set up CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Set up a new router under /v1
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/user", apiConfig.handlerCreateUser)                           // Create User
	v1Router.Get("/user", apiConfig.middlewareAuth(apiConfig.handlerGetUser))     //Get User (auth)
	v1Router.Post("/feed", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed)) //Get User (auth)
	v1Router.Get("/feeds", apiConfig.handlerGetFeeds)                             //Get Feeds

	// Mount the v1Router under /v1
	router.Mount("/v1", v1Router)

	// Start the server
	server := &http.Server{
		Handler: router,
		Addr:    "localhost:" + port,
	}
	err = server.ListenAndServe() // Blocks until the server is closed
	if err != nil {
		log.Fatal(err)
	}
}
