package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/geneowak/go-expense-tracker/handlers"
	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL env variable has not been set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM env variable has not been set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET env variable has not been set")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to open the DB")
	}
	const filePathRoot string = "."
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	/** we use a single validate instance of this because of
	* 1. heavy initialization cost
	* 2. it is thread safe and so can be used simultaneously
	* 3. when we add custom registrations, they'll be carried through to all users
	**/
	validate := handlers.CreateValidator()

	cfg := handlers.ApiConfig{
		DB:        database.New(db),
		Platform:  platform,
		JwtSecret: jwtSecret,
		Validate:  validate,
	}
	server := handlers.SetupServer(&cfg, filePathRoot, port)

	log.Println("Listening on port:", port)
	log.Fatal(server.ListenAndServe())
}
