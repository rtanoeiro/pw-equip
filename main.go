package main

import (
	"fmt"
	"log"
	"net/http"
	"pw-equip-change/api"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	equipCfg := api.EquipConfig{}
	equipCfg.LoadEquipConfig()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/download/*", http.StripPrefix("/download/", http.FileServer(http.Dir("./download"))))

	// Pw Equip Changer Endpoints
	r.Post("/equip/payment-events", equipCfg.HandlePaymentEvents)
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
