package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func RegisterEmailWithHWID(email string, hwid string) User {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	user := User{Email: email, Hwid: hwid, Active: false, Error: ""}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s?email=%s&hwid=%s", registerUserURL, email, hwid), nil)
	if err != nil {
		user.Error = fmt.Sprintf("falha ao criar requisição: %v", err)
		return user
	}

	log.Printf("Making %s request to register User %s with HWID %s", request.Method, email, hwid)
	request.Header.Set("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		return user
	}

	switch response.StatusCode {
	case http.StatusOK:
		log.Printf("User %s registered successfully with HWID %s", email, hwid)
		user.Active = true
		return user
	case http.StatusCreated:
		log.Printf("User %s registered successfully with HWID %s", email, hwid)
		user.Active = true
		return user
	case http.StatusForbidden:
		log.Printf("User %s is not active", email)
		user.Active = false
		user.Error = "HWID diferente do registrado.\nDesvincule o HWID atual e registre novamente nesse PC"
		return user
	case http.StatusInternalServerError:
		log.Printf("Error registering user %s with HWID %s", email, hwid)
		user.Active = false
		user.Error = "Erro ao receber dados do usuário.\nPor favor, tente novamente ou contate o suporte"
		return user
	case http.StatusPreconditionFailed:
		log.Printf("User %s does not have active subscription", email)
		user.Active = false
		user.Error = "Assinatura não ativa.\nPor favor, ative sua assinatura, para comprar acesse https://gamedevforge.ovh"
		return user
	default:
		log.Printf("Error registering user %s with HWID %s. No valid response code", email, hwid)
		user.Active = false
		user.Error = "Erro no registro do usuario.\nPor favor, tente novamente ou contate o suporte"
		return user
	}
}
