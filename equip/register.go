package main

import (
	"fmt"
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

	request.Header.Set("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		return user
	}

	switch response.StatusCode {
	case http.StatusOK:
		user.Active = true
		return user
	case http.StatusCreated:
		user.Active = true
		return user
	case http.StatusForbidden:
		user.Active = false
		user.Error = "HWID diferente do registrado.\nDesvincule o HWID atual e registre novamente nesse PC"
		return user
	case http.StatusInternalServerError:
		user.Active = false
		user.Error = "Erro ao receber dados do usuário.\nPor favor, tente novamente ou contate o suporte"
		return user
	case http.StatusPreconditionFailed:
		user.Active = false
		user.Error = "Assinatura não ativa.\nPor favor, ative sua assinatura, para comprar acesse https://gamedevforge.ovh"
		return user
	default:
		user.Active = false
		user.Error = "Erro no registro do usuario.\nPor favor, tente novamente ou contate o suporte"
		return user
	}
}
