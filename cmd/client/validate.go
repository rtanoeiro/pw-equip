package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

var MaxRetries = 3

// GetHWID generates a unique hardware ID for the current machine
func GetHWID() (string, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("failed to get host info: %v", err)
	}

	identifier := fmt.Sprintf("%s-%s-%s-%s",
		hostInfo.HostID,
		hostInfo.Platform,
		hostInfo.PlatformFamily,
		runtime.GOARCH,
	)

	hash := md5.Sum([]byte(identifier))
	hwid := fmt.Sprintf("%x", hash)

	return hwid, nil
}

func ValidadeUser(email, hwid string) User {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	user := User{Email: email, Hwid: hwid, Active: false, Error: ""}

	// Make GET request to subscription API
	url := fmt.Sprintf("%s?email=%s&hwid=%s", validateUserURL, email, hwid)
	log.Printf("Validating user: %s with HWID: %s", email, hwid)
	response, err := client.Get(url)
	if err != nil {
		user.Error = fmt.Sprintf("falha ao checar usuario: %v", err)
		return user
	}

	switch response.StatusCode {
	case http.StatusOK:
		user.Active = true
		return user
	case http.StatusForbidden:
		user.Active = false
		user.Error = "HWID diferente do registrado.\nDesvincule o HWID atual e tente novamente"
		return user
	case http.StatusPreconditionFailed:
		user.Active = false
		user.Error = "Assinatura não ativa.\nPor favor, ative sua assinatura, para comprar acesse https://painelguildpw.com.br"
		return user
	case http.StatusInternalServerError:
		user.Active = false
		user.Error = "Erro ao receber dados do usuário.\nPor favor, tente novamente ou contate o suporte"
		return user
	case http.StatusNotFound:
		user.Active = false
		user.Error = "Email não encontrado.\nPor favor, realize a compra do programa ou entre em contato com o suporte caso já tenha realizado"
		return user
	default:
		user.Active = false
		user.Error = "Erro ao receber dados do usuário.\nPor favor, tente novamente ou contate o suporte"
		return user
	}
}

func ValidateEmailWithHWID(email string, hwid string) {
	log.Printf("Validating user: %s with HWID: %s", email, hwid)
	request, errorRequest := http.NewRequest("GET", fmt.Sprintf("%s?email=%s&hwid=%s", validateUserURL, email, hwid), nil)
	if errorRequest != nil {
		log.Fatal(errorRequest)
	}
	request.Header.Set("Content-Type", "text/plain")
	response, errorResponse := http.DefaultClient.Do(request)
	if errorResponse != nil {
		log.Fatal(errorResponse)
	}
	defer response.Body.Close()

	// read body and log it
	if response.StatusCode != http.StatusOK {
		body, errorBody := io.ReadAll(response.Body)
		if errorBody != nil {
			log.Fatal(errorBody)
		}
		log.Fatalf("response.StatusCode: %d, response.Body: %s", response.StatusCode, string(body))
	}
}
