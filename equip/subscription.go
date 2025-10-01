package equip

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

var MaxRetries = 3

// GetHWID generates a unique hardware ID for the current machine
func GetHWID() (string, error) {
	// Get system information
	hostInfo, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("failed to get host info: %v", err)
	}

	// Create a unique identifier based on system information
	identifier := fmt.Sprintf("%s-%s-%s-%s",
		hostInfo.HostID,
		hostInfo.Platform,
		hostInfo.PlatformFamily,
		runtime.GOARCH,
	)

	// Generate MD5 hash of the identifier
	hash := md5.Sum([]byte(identifier))
	hwid := fmt.Sprintf("%x", hash)

	return hwid, nil
}

func RegisterEmailWithHWID(email string, hwid string) User {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	user := User{Email: email, Hwid: hwid, Active: false, Error: ""}

	request, err := http.NewRequest("POST", fmt.Sprintf("http://gamedevforge.ovh/register-user?email=%s&hwid=%s", email, hwid), nil)
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
		user.Error = "Erro no registro do usuário.\nPor favor, tente novamente ou contate o suporte"
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

func ResetHWID(email, hwid string) User {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	user := User{Email: email, Hwid: "", Active: false, Error: ""}

	request, err := http.NewRequest("PATCH", fmt.Sprintf("http://gamedevforge.ovh/reset-hwid?email=%s&hwid=%s", email, hwid), nil)
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
	case http.StatusBadRequest:
		user.Active = false
		user.Error = "Erro na requisição.\nPor favor, tente novamente"
		return user
	case http.StatusInternalServerError:
		user.Active = false
		user.Error = "Erro no reset do HWID.\nPor favor, tente novamente ou contate o suporte"
		return user
	case http.StatusNotFound:
		user.Active = false
		user.Error = "Usuário não encontrado ou não tem assinatura.\nPor favor, contate o suporte"
		return user
	default:
		user.Active = false
		user.Error = "Erro ao resetar o HWID.\nPor favor, tente novamente ou contate o suporte"
		return user
	}

}

// CheckSubscription verifies if the current machine has an active subscription
func CheckSubscription(email, hwid string) User {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	user := User{Email: email, Hwid: hwid, Active: false, Error: ""}

	// Make GET request to subscription API
	url := fmt.Sprintf("http://gamedevforge.ovh/validate-user?email=%s&hwid=%s", email, hwid)
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
		user.Error = "Usuário existe, mas não foram encontrados pagamentos.\nPor favor, contate o suporte, ou ative sua assinatura.\nPara comprar acesse https://gamedevforge.ovh"
		return user
	case http.StatusInternalServerError:
		user.Active = false
		user.Error = "Erro ao receber dados do usuário.\nPor favor, tente novamente ou contate o suporte"
		return user
	case http.StatusNotFound:
		user.Active = false
		user.Error = "Erro ao receber dados do usuário.\nPor favor, tente novamente ou contate o suporte"
		return user
	default:
		user.Active = false
		user.Error = "Erro ao receber dados do usuário.\nPor favor, tente novamente ou contate o suporte"
		return user
	}
}

// DisplayHWID shows the current machine's HWID for debugging/registration purposes
func DisplayHWID() {
	hwid, err := GetHWID()
	if err != nil {
		fmt.Printf("Erro ao obter HWID: %v\n", err)
		return
	}
	fmt.Printf("HWID da máquina: %s\n", hwid)
}
