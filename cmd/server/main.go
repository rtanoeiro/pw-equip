package main

import (
	"fmt"
	"log"
	"net/http"
	"pw-equip-change/internal/api"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	// Optional: load .env for local dev; in production env is set by Docker/orchestrator
	_ = godotenv.Load()

	equipCfg := &api.EquipConfig{}
	equipCfg.LoadEquipConfig()

	log.Printf("EquipConfig: %+v", equipCfg)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/equip/download/*", http.StripPrefix("/equip/download/", http.FileServer(http.Dir("./download"))))

	// Pw Equip Changer Endpoints
	r.Post("/equip/register-user", equipCfg.HandleRegisterUser)
	r.Get("/equip/validate-user", equipCfg.HandleValidateUser)
	r.Patch("/equip/reset-hwid", equipCfg.HandleResetHWID)
	r.Get("/equip/health", equipCfg.HandleHealth)

	// Start web server in a goroutine
	addr := fmt.Sprintf(":%s", equipCfg.ApiPort)
	log.Printf("Starting web server on port %s", equipCfg.ApiPort)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
