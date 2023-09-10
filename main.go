package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/vadimistar/hackathon1/adapters/ydb"
	"github.com/vadimistar/hackathon1/handlers"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	godotenv.Load()

	db, err := ydb.New()
	if err != nil {
		log.Fatalf("cannot create db connection: %s", err)
	}

	// serviceEndpoint := os.Getenv("SERVICE_ENDPOINT")
	// if serviceEndpoint == "" {
	// 	log.Fatalln("env variable SERVICE_ENDPOINT is empty")
	// }

	jwtKey := []byte(os.Getenv("JWT_KEY"))

	r.Use(cors.AllowAll().Handler)

	register, err := handlers.Register(db, db, db, jwtKey /*, serviceEndpoint */)
	if err != nil {
		log.Fatalf("create Register handler: %s", err)
	}

	r.Post("/register", register)

	r.Post("/login", handlers.Login(
		db,
		jwtKey,
	))

	r.Get("/tags", handlers.GetTags(db, db, jwtKey))

	r.Post("/tags", handlers.PostTags(db, db, jwtKey))

	r.Get("/findMatch", handlers.FindMatch(db, db, db, db, jwtKey))

	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
