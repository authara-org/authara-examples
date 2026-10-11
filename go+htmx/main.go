package main

import (
	"log"
	"net/http"

	"github.com/authara-org/authara-go/authara"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"testapp/handlers"
)

func main() {
	cfg, err := authara.ConfigFromEnv()
	if err != nil {
		log.Fatalf("authara config failed: %v", err)
	}

	appSDK, err := authara.New(cfg)
	if err != nil {
		log.Fatalf("authara sdk init failed: %v", err)
	}
	defer appSDK.Close()

	h := handlers.New(cfg.AutharaBaseURL)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/", h.Home)

	r.Group(func(r chi.Router) {
		r.Use(appSDK.RequireAuthWithRefresh)
		r.Get("/private", h.Private)
		r.Get("/private/current-user", h.CurrentUser)
	})

	log.Println("testapp listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
