package api

import (
	"log"
	"net/http"
)

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
