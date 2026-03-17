package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/authara-org/authara-go/authara"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"testapp/handlers"
)

func main() {
	autharaBaseURL := envOrDefault("AUTHARA_BASE_URL", "http://authara:8080")

	keys, err := parseJWTKeys(os.Getenv("AUTHARA_JWT_KEYS"))
	if err != nil {
		log.Fatalf("parse AUTHARA_JWT_KEYS: %v", err)
	}

	appSDK, err := authara.New(authara.Config{
		Audience:       envOrDefault("AUTHARA_AUDIENCE", "app"),
		Issuer:         envOrDefault("AUTHARA_ISSUER", "authara"),
		Keys:           keys,
		AutharaBaseURL: autharaBaseURL,
	})
	if err != nil {
		log.Fatalf("authara sdk init failed: %v", err)
	}

	h := handlers.New(autharaBaseURL)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/", h.Home)

	r.Group(func(r chi.Router) {
		r.Use(appSDK.RequireAuthWithRefresh)
		r.Get("/private", h.Private)
		r.Get("/private/pulse", h.PrivatePulse)
	})

	log.Println("testapp listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func parseJWTKeys(raw string) (map[string][]byte, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errEnv("AUTHARA_JWT_KEYS is empty")
	}

	out := make(map[string][]byte)

	entries := strings.Split(raw, ",")
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 {
			return nil, errEnv("invalid AUTHARA_JWT_KEYS entry: " + entry)
		}

		keyID := strings.TrimSpace(parts[0])
		b64 := strings.TrimSpace(parts[1])

		if keyID == "" {
			return nil, errEnv("empty key id in AUTHARA_JWT_KEYS")
		}
		if b64 == "" {
			return nil, errEnv("empty key value for key id " + keyID)
		}

		key, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, errEnv("invalid base64 for key id " + keyID + ": " + err.Error())
		}

		out[keyID] = key
	}

	if len(out) == 0 {
		return nil, errEnv("no valid keys found in AUTHARA_JWT_KEYS")
	}

	return out, nil
}

type errEnv string

func (e errEnv) Error() string { return string(e) }

func envOrDefault(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}
