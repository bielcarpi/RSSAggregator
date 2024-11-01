package main

import (
	"database/sql"
	"fmt"
	"github.com/bielcarpi/RSSAggregator/internal/db"
	"log"
	"net/http"
	"os"
	"time"

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
	dbConn := db.New(conn)
	apiConfig := apiConfig{
		DB: dbConn,
	}

	fmt.Println("Server starting on port " + port)
	router := chi.NewRouter()

	// Set up CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Set up a new router under /v1
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)

	v1Router.Post("/user", apiConfig.handlerCreateUser)
	v1Router.Get("/user", apiConfig.middlewareAuth(apiConfig.handlerGetUser))

	v1Router.Post("/feed", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feeds", apiConfig.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollows))

	// Mount the v1Router under /v1
	router.Mount("/v1", v1Router)

	// Start scrapping for RSS Feeds
	go startScraping(dbConn, 10, time.Minute)

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
