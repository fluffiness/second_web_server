module fluffiness/second_web_server

go 1.26.2

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1 // For loading environmental variables from the .env file
	github.com/lib/pq v1.12.3 // PostgreSQL driver for the sqlc generated queries
)
