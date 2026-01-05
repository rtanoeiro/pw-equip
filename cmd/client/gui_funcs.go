package main

import (
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	hook "github.com/robotn/gohook"
)

// UpdateButtonState returns a closure that reads the current values of isEmailValid and isSubscriptionValid
// We use pointers so the function always reads the latest values when called
func UpdateButtonState(startButton **widget.Button, resetHWIDButton **widget.Button, isEmailValid *bool, isSubscriptionValid *bool) func() {
	return func() {
		if startButton != nil && *startButton != nil {
			if *isEmailValid && *isSubscriptionValid {
				(*startButton).Enable()
			} else {
				(*startButton).Disable()
			}
		}
		if resetHWIDButton != nil && *resetHWIDButton != nil {
			if *isEmailValid {
				(*resetHWIDButton).Enable()
			} else {
				(*resetHWIDButton).Disable()
			}
		}
	}
}

func OnEmailChanged(emailStatusLabel *widget.Label, isEmailValid *bool, isSubscriptionValid *bool, updateStateButton func()) func(string) {
	return func(email string) {
		*isEmailValid = false
		*isSubscriptionValid = false
		updateStateButton()

		if email == "" {
			emailStatusLabel.SetText("")
			return
		}

		if !IsValidEmail(email) {
			emailStatusLabel.SetText("❌ Email inválido")
			return
		}

		emailStatusLabel.SetText("✅ Email válido - Pressione Enter para registrar")
		*isEmailValid = true

		updateStateButton()
	}
}

func OnEmailSubmitted(hwid string, emailStatusLabel *widget.Label, isEmailValid *bool, isSubscriptionValid *bool, updateStateButton func()) func(string) {
	return func(email string) {
		if !IsValidEmail(email) {
			emailStatusLabel.SetText("❌ Email inválido")
			*isEmailValid = false
			*isSubscriptionValid = false
			updateStateButton()
			return
		}

		emailStatusLabel.SetText("🔄 Registrando email...")
		*isSubscriptionValid = false
		updateStateButton()

		// Run registration and subscription check in goroutine
		go func() {
			fyne.Do(func() {
				userReg := RegisterEmailWithHWID(email, hwid)

				// Update UI in main thread
				emailStatusLabel.SetText("🔄 Verificando assinatura...")

				if userReg.Active {
					// Check subscription
					userSub := ValidadeUser(email, hwid)
					if userSub.Error != "" {
						emailStatusLabel.SetText("⚠️ Erro ao verificar assinatura: \n" + userSub.Error)
						*isSubscriptionValid = false
					} else if userSub.Active {
						emailStatusLabel.SetText("✅ Email registrado e assinatura ativa")
						*isEmailValid = true
						*isSubscriptionValid = true
					} else {
						emailStatusLabel.SetText("⚠️ Assinatura inválida: \n" + userSub.Error)
						*isEmailValid = true
						*isSubscriptionValid = false
					}
				} else {
					emailStatusLabel.SetText("⚠️ " + userReg.Email + ": " + userReg.Error)
					*isEmailValid = false
					*isSubscriptionValid = false
				}

				updateStateButton()
			})
		}()
	}
}

func OnChangeSetKeyButtonClicked(changeSetKeyButton **widget.Button, config *Config, isMonitoring *bool) func() {
	return func() {
		// Check if button exists and is not nil
		if changeSetKeyButton == nil || *changeSetKeyButton == nil {
			log.Printf("ERROR: changeSetKeyButton is nil!")
			return
		}

		// Prevent capturing keys while monitoring is active to avoid hook conflicts
		if *isMonitoring {
			log.Printf("Cannot capture key while monitoring is active")
			fyne.Do(func() {
				(*changeSetKeyButton).SetText("Pare o monitoramento primeiro!")
			})
			return
		}

		// Show user that we're waiting for key press
		(*changeSetKeyButton).SetText("Pressione uma tecla...")
		(*changeSetKeyButton).Disable() // Prevent multiple clicks

		// Run hook listener in a goroutine to avoid blocking the UI
		go func() {
			// Start listening for keyboard events
			evChan := hook.Start()

			// Wait for the FIRST key press event
			for ev := range evChan {
				// Only capture KeyDown events (ignore KeyUp)
				if ev.Kind == hook.KeyDown {
					// Store both keycode and character representation
					config.ChangeSetKeyCode = ev.Keycode
					config.ChangeSetKeyChar = string(ev.Keychar)
					log.Printf("Captured key: %s (code: %d)", config.ChangeSetKeyChar, config.ChangeSetKeyCode)

					// Stop the hook listener immediately after capturing one key
					hook.End()

					// Update UI - must use fyne.Do() to update from goroutine
					fyne.Do(func() {
						// Show the captured key on the button
						(*changeSetKeyButton).SetText(fmt.Sprintf("Tecla: %s (código: %d)", string(ev.Keychar), ev.Keycode))
						(*changeSetKeyButton).Enable() // Re-enable the button
					})

					// Exit the goroutine
					return
				}
			}
		}()
	}
}

func UpdateItemKeys(config *Config, numItems int, itemKeysContainer *fyne.Container, itemKeyEntries *[]*widget.Entry) func(int) {
	return func(numItems int) {
		itemKeysContainer.RemoveAll()
		*itemKeyEntries = make([]*widget.Entry, numItems)

		for i := 0; i < numItems; i++ {
			entry := widget.NewEntry()

			label := widget.NewLabel(fmt.Sprintf("Item %d:", i+1))
			// In case we have saved keys, auto populate them
			if i < len(config.Keys) {
				if config.Keys[i] != "" {
					entry.SetText(config.Keys[i])
				}
			} else {
				// If we don't have them saved, set this placeholder
				entry.SetPlaceHolder(fmt.Sprintf("Tecla do item %d", i+1))
			}
			(*itemKeyEntries)[i] = entry
			itemKeysContainer.Add(container.NewHBox(label, entry))
		}
		itemKeysContainer.Refresh()
	}
}

// Whenever we update the number of items that need change, run this function.
// It calls updateItemKeys which will reuse saved config, whenever available.
func UpdateItemKeysOnChanged(numItemsEntry *widget.Entry, updateItemKeys func(int)) func(string) {
	return func(text string) {
		if num, err := strconv.Atoi(text); err == nil && num >= 1 && num <= 11 {
			updateItemKeys(num)
		}
	}
}

func GetHWIDLabel(hwidLabel *widget.Label) {
	go func() {
		fyne.Do(func() {
			hwid, err := GetHWID()
			if err == nil {
				hwidLabel.SetText(fmt.Sprintf("HWID: %s", hwid))
			} else {
				hwidLabel.SetText("Erro ao obter HWID")
			}
		})
	}()
}

func (g *GuiApp) StartButton(
	emailEntry *widget.Entry,
	hwid string,
	keyBarShiftEntry *widget.Entry,
	InBetweenTimeClicksEntry *widget.Entry,
	isEmailValid *bool,
	isSubscriptionValid *bool,
	config *Config,
	itemKeyEntries *[]*widget.Entry,
	numItemsEntry *widget.Entry,
	statusLabel *widget.Label,
	startButton **widget.Button,
	stopMonitoring *chan bool,
	isMonitoring *bool,
) func() {
	return func() {
		if (*startButton).Text == "Parar Monitoramento" {
			if *stopMonitoring != nil {
				close(*stopMonitoring)
			}
			*isMonitoring = false
			statusLabel.SetText("Monitoramento parado.")
			(*startButton).SetText("Iniciar Monitoramento")
			return
		}

		// Read email when button is clicked, not when created
		email := emailEntry.Text
		if !IsValidEmail(email) {
			log.Printf("Invalid email. Email used: %s", email)
			statusLabel.SetText("Erro: Digite um email válido")
			return
		}

		// Check if email and subscription are valid (already verified)
		if !*isEmailValid || !*isSubscriptionValid {
			log.Printf("Email is either invalid or usubscribed. Invalid: %v. Unsubscribed %v", isEmailValid, isSubscriptionValid)
			statusLabel.SetText("Erro: Email deve estar registrado e assinatura ativa")
			return
		}

		// Get the captured keycode from config
		// We already stored this when the user clicked the capture button
		changeSetKeyCode := config.ChangeSetKeyCode
		if changeSetKeyCode == 0 {
			log.Printf("User did not configure change set key")
			statusLabel.SetText("Erro: Defina a tecla para trocar de set")
			return
		}

		// Validate inputs
		numItems, err := strconv.Atoi(numItemsEntry.Text)
		if err != nil || numItems < 1 || numItems > 11 {
			log.Printf("User entered wrong number of items")
			statusLabel.SetText("Erro: Número de items deve ser entre 1 e 11")
			return
		}

		keyBarShift := keyBarShiftEntry.Text
		if keyBarShift != "v" && keyBarShift != "`" && keyBarShift != "'" {
			log.Printf("User select wrong key to change bar")
			statusLabel.SetText("Erro: Tecla deve ser *v* ou *`* ou *'*")
			return
		}

		inBetweenTimeClicks, err := strconv.Atoi(InBetweenTimeClicksEntry.Text)
		if err != nil || inBetweenTimeClicks < 0 {
			log.Printf("User selected wrong timing in between changes")
			statusLabel.SetText("Erro: Tempo deve ser um número válido")
			return
		}

		// Collect item keys
		itemKeys := make([]string, numItems)
		for i := 0; i < numItems; i++ {
			if i < len(*itemKeyEntries) && (*itemKeyEntries)[i].Text != "" {
				itemKeys[i] = (*itemKeyEntries)[i].Text
			} else {
				log.Printf("User attempted to started before filling all keys. Missing key %d", i+1)
				statusLabel.SetText(fmt.Sprintf("Erro: Digite a tecla para o item %d", i+1))
				return
			}
		}
		// Check subscription before starting monitoring
		statusLabel.SetText("Verificando assinatura...")
		(*startButton).SetText("Verificando...")

		// Check subscription with retry
		user := ValidadeUser(email, hwid)
		if !user.Active {
			statusLabel.SetText(user.Error)
			(*startButton).SetText("Iniciar Monitoramento")
			(*startButton).Disable()
			return
		}

		// Setup configuration
		g.setup = &SetupEquip{
			NumberItems:         numItems,
			KeyChange:           keyBarShift,
			InBetweenTimeClicks: inBetweenTimeClicks,
			ItemKeys:            itemKeys,
			CurrentSet:          1,
		}

		(*startButton).Enable()
		statusLabel.SetText(fmt.Sprintf("Assinatura ativa! Monitoramento iniciado. Pressione %s para trocar de set.", config.ChangeSetKeyChar))
		(*startButton).SetText("Parar Monitoramento")

		// Start monitoring in a goroutine
		*stopMonitoring = make(chan bool)
		*isMonitoring = true
		log.Printf("Saving configuration into file")
		errConfig := SaveConfig(email, hwid, keyBarShift, config.ChangeSetKeyChar, inBetweenTimeClicks, itemKeys, changeSetKeyCode)
		if errConfig != nil {
			log.Printf("Failed to save config. Error %s", errConfig)
			return
		}
		go g.startMonitoring(statusLabel, *stopMonitoring, changeSetKeyCode)
	}
}

func ResetHWIDButton(emailEntry *widget.Entry, hwid string, emailWidget *widget.Label) func() {
	return func() {
		email := emailEntry.Text
		user := ResetHWID(email, hwid)
		if user.Active {
			emailWidget.SetText("✅ HWID resetado")
		} else {
			emailWidget.SetText("⚠️ " + user.Error)
		}
	}
}
