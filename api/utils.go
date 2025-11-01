package api

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type UserCredentials struct {
	Email string
	Hwid  string
}

func parseEmailAndHWID(request *http.Request) (UserCredentials, bool) {
	log.Printf("Received Parsing request from Query %s", request.URL.RawQuery)
	data, errParse := url.ParseQuery(request.URL.RawQuery)
	if errParse != nil {
		return UserCredentials{}, false
	}

	email := strings.Trim(data.Get("email"), " ")
	hwid := strings.Trim(data.Get("hwid"), " ")

	if email == "" || hwid == "" {
		log.Printf("Email or HWID is empty")
		return UserCredentials{Email: email, Hwid: hwid}, false
	}
	log.Printf("Received request for user '%s' with HWID '%s'", email, hwid)

	return UserCredentials{Email: email, Hwid: hwid}, true
}

func GetEnvVar(key, defaultValue string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return defaultValue
}
