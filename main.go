package main

import (
	"database/sql"
	"fluffiness/second_web_server/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32      // Use atomic values to prevent race conditions, since the net/http package handles requests concurrently.
	db             *database.Queries // Queries compiled by sqlc
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// =================== Databse ===================
	// The code in internal/database is automatically compiled from SQL to Go from the sql/queries directory.
	// The dbURL argument is the postgres connection string loaded from the .env file.
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	dbQueries := database.New(db)

	cfg := apiConfig{db: dbQueries, platform: platform}
	mux := http.NewServeMux()

	// =================== Routing ===================
	// Method specific routing is done by specifying a method in the pattern explicitly in the form [METHOD ][HOST]/[PATH]
	// When accessed by the wrong method, the http/net package automatically responds with status code 405 Method Not Allowed
	// Matching is done by checking if the path segment of a url has the pattern as its prefix.
	// If there are multiple matches, the longest one wins.

	// =================== File Server ===================
	// The handler checks if the path of the segment of the url has the prefix "/app/"
	// If so, http.StripPrefix strips "/app" from the path, and concatenates it to filepathRoot
	// The file at the resulting path is served.
	// If the path points to a directory, the index.html file in the directory is served.
	fileSrv := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app", fileSrv)))

	mux.HandleFunc("GET /admin/metrics/", cfg.handlerMetrics)    // Returns the number of accesses to the fileserver
	mux.HandleFunc("POST /admin/reset/", cfg.handleMetricsReset) // Resets the access count and deletes all users. Only authorized on dev platforms.

	mux.HandleFunc("GET /api/healthz/", handlerReadiness)            // Health check to see if the server is online
	mux.HandleFunc("POST /api/validate_chirp/", handlerValidate)     // Checks character count and removes profanity.
	mux.HandleFunc("POST /api/users/", cfg.handlerCreateUser)        // Create a new user
	mux.HandleFunc("POST /api/chirps/", cfg.handlerCreateChirp)      // Create a new chirp
	mux.HandleFunc("GET /api/chirps/", cfg.handlerGetChirps)         // Get all chirps ordered by created_at
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp) // Get chirp with id matching ChirpID

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
