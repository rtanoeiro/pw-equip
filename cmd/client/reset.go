package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func ResetHWID(email, hwid string) User {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	user := User{Email: email, Hwid: "", Active: false, Error: ""}
	log.Printf("Resetting HWID for user %s with HWID %s", email, hwid)
	request, err := http.NewRequest("PATCH", fmt.Sprintf("%s?email=%s&hwid=%s", resetHWIDURL, email, hwid), nil)
	if err != nil {
		user.Error = fmt.Sprintf("falha ao criar requisição: %v", err)
		return user
	}

	request.Header.Set("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		return user
	}

	switch response.StatusCode {
	case http.StatusOK:
		log.Printf("HWID reset successfully for user %s", email)
		user.Active = true
		return user
	case http.StatusBadRequest:
		log.Printf("Error resetting HWID for user %s", email)
		user.Active = false
		user.Error = "Erro na requisição.\nPor favor, tente novamente"
		return user
	case http.StatusInternalServerError:
		log.Printf("Error resetting HWID for user %s", email)
		user.Active = false
		user.Error = "Erro no reset do HWID.\nPor favor, tente novamente ou contate o suporte"
		return user
	case http.StatusNotFound:
		log.Printf("User %s not found", email)
		user.Active = false
		user.Error = "Usuário não encontrado ou não tem assinatura.\nPor favor, contate o suporte"
		return user
	default:
		log.Printf("Error resetting HWID for user %s. No valid response code", email)
		user.Active = false
		user.Error = "Erro ao resetar o HWID.\nPor favor, tente novamente ou contate o suporte"
		return user
	}

}
