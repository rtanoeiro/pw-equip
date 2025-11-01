package api

import (
	"log"
	"net/http"
)

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
