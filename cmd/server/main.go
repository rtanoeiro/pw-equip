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
	envFiles := []string{".env-dev", ".env-test", ".env-prod"}
	var envFile string
	for _, envFile := range envFiles {
		err := godotenv.Load(envFile)
		if err == nil {
			break
		}
	}
	log.Printf("Loaded environment variables from %s", envFile)

	equipCfg := &api.EquipConfig{}
	equipCfg.LoadEquipConfig()

	log.Printf("EquipConfig: %+v", equipCfg)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/equip/download/*", http.StripPrefix("/equip/download/", http.FileServer(http.Dir("./download"))))
	r.Get("/equip/equipment-changer", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/equipment-changer.html")
	})
	r.Get("/equip/equipment-changer-en", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/equipment-changer-en.html")
	})

	// Pw Equip Changer Endpoints
	r.Post("/equip/payment-events", equipCfg.HandlePaymentEvents)
	r.Post("/equip/register-user", equipCfg.HandleRegisterUser)
	r.Get("/equip/validate-user", equipCfg.HandleValidateUser)
	r.Patch("/equip/reset-hwid", equipCfg.HandleResetHWID)
	r.Get("/equip/health", equipCfg.HandleHealth)

	// Start web server in a goroutine
	addr := fmt.Sprintf(":%s", equipCfg.ApiPort)
	log.Printf("Starting web server on port %s", equipCfg.Config.ApiPort)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
