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

func ClickButton(button string) {
	robotgo.KeyPress(button)
	// Note: Errors are silently ignored to prevent console output in GUI mode
}

func ChangeItems(equipSetup *SetupEquip) {
	if equipSetup.CurrentSet == 1 {
		ClickButton(equipSetup.KeyChange)
		ClickButton(equipSetup.KeyChange)
		for _, itemToPress := range equipSetup.ItemKeys {
			ClickButton(itemToPress)
			time.Sleep(time.Duration(equipSetup.TimeClicks) * time.Millisecond)
		}
		ClickButton(equipSetup.KeyChange)
		equipSetup.CurrentSet = 2
	}

	if equipSetup.CurrentSet == 2 {
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

func RegisterEmailWithHWID(email string, hwid string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("http://gamedevforge.ovh/register-user?email=%s&hwid=%s", email, hwid), nil)
	if err != nil {
		return fmt.Errorf("falha ao criar requisição: %v", err)
	}

	request.Header.Set("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("falha ao enviar requisição: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("falha ao ler o corpo da resposta: %v", err)
		}
		return fmt.Errorf("falha no registro com Status %d: %s", response.StatusCode, string(body))
	}

	return nil
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
