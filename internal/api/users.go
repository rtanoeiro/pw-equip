package api

import (
	"database/sql"
	"log"
	"net/http"
	"pw-equip-change/internal/database"
)

var paymentValueTreshold = 2000

func (config *Config) HandleRegisterUser(w http.ResponseWriter, request *http.Request) {
	credentials, valid := parseEmailAndHWID(request)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if config.UserExists(credentials.Email) {
		log.Printf("HandleRegisterUser: User '%s' already exists in database", credentials.Email)
		user, errUser := config.DB.GetUserByEmail(request.Context(), credentials.Email)
		if errUser != nil {
			log.Printf("HandleRegisterUser: Error getting user '%s': %s", credentials.Email, errUser.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("HandleRegisterUser: Found user '%s' with HWID '%s'", credentials.Email, user.Hwid.String)

		// User exists with different HWID
		if user.Hwid.Valid && user.Hwid.String != credentials.Hwid {
			log.Printf("HandleRegisterUser: '%s' exists with different HWID '%s'", credentials.Email, credentials.Hwid)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// User exists but hasn't paid
		if !config.UserPaid(credentials.Email) {
			log.Printf("HandleRegisterUser: '%s' has not paid", credentials.Email)
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}

		// User exists but doesn't have a HWID, so we trigger an HWID update
		if user.Email != "" && !user.Hwid.Valid {
			config.UpdateHWID(credentials, request, w)
			return
		}

		// User exists with correct HWID
		log.Printf("HandleRegisterUser: '%s' already registered with correct HWID", credentials.Email)
		w.WriteHeader(http.StatusOK)
		return
	}

	// User doesn't exist, but has paid, it's rare, but I guess it could happen
	if !config.UserPaid(credentials.Email) {
		log.Printf("HandleRegisterUser: Failed to get charges for user '%s", credentials.Email)
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	// User doesn't exist, and has paid, as the check above failed, so we create a new user
	createUser(credentials, config, request, w)
	w.WriteHeader(http.StatusCreated)
}

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

func (config *Config) HandleValidateUser(w http.ResponseWriter, request *http.Request) {
	credentials, valid := parseEmailAndHWID(request)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if config.UserExists(credentials.Email) {
		user, errUser := config.DB.GetUserByEmail(request.Context(), credentials.Email)
		if errUser != nil {
			log.Printf("HandleValidateUser: Error getting user %s: %v", credentials.Email, errUser)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("HandleValidateUser: Found user %s with HWID %s", credentials.Email, user.Hwid.String)

		// Exists but with different HWID
		if user.Hwid.Valid && user.Hwid.String != credentials.Hwid {
			log.Printf("HandleValidateUser: User %s already exists with different HWID %s", credentials.Email, credentials.Hwid)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Exists but hasn't paid
		if !config.UserPaid(credentials.Email) {
			log.Printf("HandleValidateUser: '%s' has not paid", credentials.Email)
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}

		// If all goes well, user exists, has paid and has correct HWID
		log.Printf("HandleValidateUser: '%s' already exists with correct HWID %s", credentials.Email, credentials.Hwid)
		w.WriteHeader(http.StatusOK)
		return
	}

	// If we reach this point, user wasn't found
	log.Printf("HandleValidateUser: '%s' not found", credentials.Email)
	w.WriteHeader(http.StatusNotFound)
}

func (config *Config) HandleResetHWID(w http.ResponseWriter, request *http.Request) {
	credentials, valid := parseEmailAndHWID(request)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if config.UserExists(credentials.Email) && config.UserPaid(credentials.Email) {
		errReset := config.DB.ResetUserHWID(request.Context(), credentials.Email)
		if errReset != nil {
			log.Printf("HandleResetHWID: Error resetting HWID for email %s: %v", credentials.Email, errReset)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("HandleResetHWID: '%s' HWID reset", credentials.Email)
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Printf("HandleResetHWID: '%s' not found or hasn't paid", credentials.Email)
	w.WriteHeader(http.StatusNotFound)
}
