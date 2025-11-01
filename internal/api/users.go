package api

import (
	"database/sql"
	"log"
	"net/http"
	"pw-equip-change/internal/database"
)

var paymentValueTreshold = 2000

func createUser(credentials UserCredentials, config *Config, request *http.Request, w http.ResponseWriter) bool {
	createUser := database.CreateUserParams{
		Email: credentials.Email,
		Hwid:  sql.NullString{String: credentials.Hwid, Valid: true},
	}
	errCreate := config.DB.CreateUser(request.Context(), createUser)
	if errCreate != nil {
		log.Printf("HandleRegisterUser: Error creating user '%s': '%s'", credentials.Email, errCreate.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}
	log.Printf("HandleRegisterUser: Successfully created new user '%s'", credentials.Email)
	return false
}

func (config *Config) UpdateHWID(credentials UserCredentials, request *http.Request, w http.ResponseWriter) bool {
	updateParams := database.UpdateUserHWIDParams{
		Hwid:  sql.NullString{String: credentials.Hwid, Valid: true},
		Email: credentials.Email,
	}
	errUpdate := config.DB.UpdateUserHWID(request.Context(), updateParams)
	if errUpdate != nil {
		log.Printf("HandleRegisterUser: Error updating user '%s': '%s'", credentials.Email, errUpdate.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}
	log.Printf("HandleRegisterUser: Successfully updated HWID for user '%s'", credentials.Email)
	w.WriteHeader(http.StatusOK)
	return false
}
