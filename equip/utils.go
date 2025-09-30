package equip

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

type User struct {
	Email  string
	Hwid   string
	Active bool
	Error  string
}

func ClickButton(button string) {
	robotgo.KeyPress(button)
	time.Sleep(10 * time.Millisecond)
	// Note: Errors are silently ignored to prevent console output in GUI mode
}

func ChangeItems(equipSetup *SetupEquip) {
	switch equipSetup.CurrentSet {
	case 1:
		ClickButton(equipSetup.KeyChange)
		ClickButton(equipSetup.KeyChange)
		for _, itemToPress := range equipSetup.ItemKeys {
			ClickButton(itemToPress)
			time.Sleep(time.Duration(equipSetup.TimeClicks) * time.Millisecond)
		}
		ClickButton(equipSetup.KeyChange)
		equipSetup.CurrentSet = 2
	case 2:
		ClickButton(equipSetup.KeyChange)
		for _, itemToPress := range equipSetup.ItemKeys {
			ClickButton(itemToPress)
			time.Sleep(time.Duration(equipSetup.TimeClicks) * time.Millisecond)
		}
		ClickButton(equipSetup.KeyChange)
		ClickButton(equipSetup.KeyChange)
		equipSetup.CurrentSet = 1
	}
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

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Simple email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateEmailWithHWID(email string, hwid string) {
	request, errorRequest := http.NewRequest("GET", fmt.Sprintf("http://gamedevforge.ovh/validate-user?email=%s&hwid=%s", email, hwid), nil)
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
