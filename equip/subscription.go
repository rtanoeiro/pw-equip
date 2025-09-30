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
		user.Error = "HWID diferente do registrado"
		return user
	case http.StatusInternalServerError:
		user.Active = false
		user.Error = "Erro ao criar usuário, tente novamente"
		return user
	default:
		user.Active = false
		user.Error = "Falha no registro do usuario, tente novamente"
		return user
	}
}

// CheckSubscription verifies if the current machine has an active subscription
func CheckSubscription(email, hwid string) (bool, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make GET request to subscription API
	url := fmt.Sprintf("http://gamedevforge.ovh/validate-user?email=%s&hwid=%s", email, hwid)
	resp, err := client.Get(url)
	if err != nil {
		return false, fmt.Errorf("falha ao checar usuario: %v", err)
	}
	defer resp.Body.Close()
	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("validacao de usuario API retornou status: %d", resp.StatusCode)
	}

	return true, nil
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
