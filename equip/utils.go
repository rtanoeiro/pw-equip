package main

import (
	"fmt"
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

// DisplayHWID shows the current machine's HWID for debugging/registration purposes
func DisplayHWID() {
	hwid, err := GetHWID()
	if err != nil {
		fmt.Printf("Erro ao obter HWID: %v\n", err)
		return
	}
	fmt.Printf("HWID da máquina: %s\n", hwid)
}
