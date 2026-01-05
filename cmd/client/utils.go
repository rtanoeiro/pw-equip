package main

import (
	"log"
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
		log.Printf("Current set is 1")
		ClickButton(equipSetup.KeyChange)
		log.Printf("Clicked key change. Clicked %s", equipSetup.KeyChange)
		ClickButton(equipSetup.KeyChange)
		log.Printf("Clicked key change. Clicked %s", equipSetup.KeyChange)
		for index, itemToPress := range equipSetup.ItemKeys {
			ClickButton(itemToPress)
			log.Printf("Clicked item number %d. Clicked %s", index+1, itemToPress)
			time.Sleep(time.Duration(equipSetup.InBetweenTimeClicks) * time.Millisecond)
		}
		ClickButton(equipSetup.KeyChange)
		log.Printf("Clicked key change. Clicked %s", equipSetup.KeyChange)
		equipSetup.CurrentSet = 2
		log.Printf("Changed to set 2")
	case 2:
		log.Printf("Current set is 2")
		ClickButton(equipSetup.KeyChange)
		log.Printf("Clicked key change. Clicked %s", equipSetup.KeyChange)
		for index, itemToPress := range equipSetup.ItemKeys {
			ClickButton(itemToPress)
			log.Printf("Clicked item number %d. Clicked %s", index+1, itemToPress)
			time.Sleep(time.Duration(equipSetup.InBetweenTimeClicks) * time.Millisecond)
		}
		ClickButton(equipSetup.KeyChange)
		log.Printf("Clicked key change. Clicked %s", equipSetup.KeyChange)
		ClickButton(equipSetup.KeyChange)
		log.Printf("Clicked key change. Clicked %s", equipSetup.KeyChange)
		equipSetup.CurrentSet = 1
		log.Printf("Changed to set 1")
	}
}

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	log.Printf("Validating email: %s", email)
	email = strings.TrimSpace(email)
	if email == "" {
		log.Printf("Email is empty")
		return false
	}

	// Simple email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	log.Printf("Email regex result: %v", emailRegex.MatchString(email))
	return emailRegex.MatchString(email)
}
