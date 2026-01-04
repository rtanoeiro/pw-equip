package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	hook "github.com/robotn/gohook"
)

type GuiApp struct {
	app    fyne.App
	window fyne.Window
	setup  *SetupEquip
}

func NewGuiApp() *GuiApp {
	myApp := app.New()
	myWindow := myApp.NewWindow("PW Equipment Changer")
	myWindow.Resize(fyne.NewSize(900, 600))

	iconPath, err := os.ReadFile("media/icon.jpg")
	if err != nil {
		fyne.LogError("Error reading icon.png", err)
	}
	iconResource := fyne.NewStaticResource("icon.png", iconPath)
	myWindow.SetIcon(iconResource)

	return &GuiApp{
		app:    myApp,
		window: myWindow,
		setup:  &SetupEquip{},
	}
}

func (g *GuiApp) RunGUI() {
	// Title
	title := widget.NewLabel("Bem vindo ao seu auxilio de troca de set")
	title.TextStyle.Bold = true

	hwid, err := GetHWID()
	if err != nil {
		hwid = "Erro ao obter HWID"
	}

	// Instructions
	instructions := widget.NewRichTextFromMarkdown(`
**Instruções:**
- Deixe 3 barras livres para serem rotacionadas
- Em sua barra principal, deixe suas skills/boticarios como deseja usa-los
- Se deseja iniciar com equipamentos de ataque, na segunda barra deixe os Equipamentos de ataque
- Na ultima barra, deixe os Equipamentos de defesa
- Para trocar de set aperte a tecla Q!
	`)

	config, errConfig := LoadConfig()
	if errConfig != nil {
		log.Printf("failed to load config file")
		return
	}
	log.Printf("Successfully loaded configuration! Config: %v", config)

	emailEntry := widget.NewEntry()

	// Load saved email
	if config.Email != "" {
		emailEntry.SetText(config.Email)
	} else {
		emailEntry.SetPlaceHolder("Digite seu email usado na compra do programa")
	}

	// Email status label
	emailStatusLabel := widget.NewLabel("")

	// Button state management
	var startButton *widget.Button
	var resetHWIDButton *widget.Button
	var isEmailValid bool = false
	var isSubscriptionValid bool = false

	// Function to update button state
	updateButtonState := func() {
		if startButton != nil {
			if isEmailValid && isSubscriptionValid {
				startButton.Enable()
			} else {
				startButton.Disable()
			}
		}
		if resetHWIDButton != nil {
			if isEmailValid {
				resetHWIDButton.Enable()
			} else {
				resetHWIDButton.Disable()
			}
		}
	}

	emailEntry.OnChanged = func(email string) {
		isEmailValid = false
		isSubscriptionValid = false
		updateButtonState()

		if email == "" {
			emailStatusLabel.SetText("")
			return
		}

		if !IsValidEmail(email) {
			emailStatusLabel.SetText("❌ Email inválido")
			return
		}

		emailStatusLabel.SetText("✅ Email válido - Pressione Enter para registrar")
		isEmailValid = true

		updateButtonState()
	}

	emailEntry.OnSubmitted = func(email string) {
		if !IsValidEmail(email) {
			emailStatusLabel.SetText("❌ Email inválido")
			isEmailValid = false
			isSubscriptionValid = false
			updateButtonState()
			return
		}

		emailStatusLabel.SetText("🔄 Registrando email...")
		isSubscriptionValid = false
		updateButtonState()

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
						isSubscriptionValid = false
					} else if userSub.Active {
						emailStatusLabel.SetText("✅ Email registrado e assinatura ativa")
						isEmailValid = true
						isSubscriptionValid = true
					} else {
						emailStatusLabel.SetText("⚠️ Assinatura inválida: \n" + userSub.Error)
						isEmailValid = true
						isSubscriptionValid = false
					}
				} else {
					emailStatusLabel.SetText("⚠️ " + userReg.Email + ": " + userReg.Error)
					isEmailValid = false
					isSubscriptionValid = false
				}

				updateButtonState()
			})
		}()
	}

	// Base config for app usage
	numItemsEntry := widget.NewEntry()	
	// In case we have configured keys, auto load them
	if config.Keys != nil {
		numItemsEntry.SetText(fmt.Sprintf("%d", len(config.Keys)))
	} else {
		numItemsEntry.SetPlaceHolder("Digite um número de 1 a 11")
	}

	keyShiftEntry := widget.NewEntry()
	if config.BarChangeKey != "" {
		keyShiftEntry.SetText(config.BarChangeKey)
	} else {
		keyShiftEntry.SetPlaceHolder("Digite 'v' ou '`'")
	}
	
	timeClicksEntry := widget.NewEntry()
	if config.TimingChange != "" {
		timeClicksEntry.SetText(config.TimingChange)
	} else {
		timeClicksEntry.SetPlaceHolder("Tempo em milisegundos. Exemplo: 1000 = 1 segundo. 200 = 0.2 segundos, quando menor, mais rapido.")
	}
	
	// Dynamic item keys conta	iner
	itemKeysContainer := container.NewVBox()
	var itemKeyEntries []*widget.Entry

	// Function to update item keys fields
	updateItemKeys := func(numItems int) {
		itemKeysContainer.RemoveAll()
		itemKeyEntries = make([]*widget.Entry, numItems)
		
		
		for i := 0; i < numItems; i++ {
			entry := widget.NewEntry()
			
			label := widget.NewLabel(fmt.Sprintf("Item %d:", i+1))
			// In case we have saved keys, auto populate them
			if i < len(config.Keys) && config.Keys[i] != "" {
				entry.SetText(config.Keys[i])
			} else {
				entry.SetPlaceHolder(fmt.Sprintf("Tecla do item %d", i+1))
			}
			itemKeyEntries[i] = entry
			itemKeysContainer.Add(container.NewHBox(label, entry))
		}
		itemKeysContainer.Refresh()
	}

	// Update item keys when number of items changes
	numItemsEntry.OnChanged = func(text string) {
		if num, err := strconv.Atoi(text); err == nil && num >= 1 && num <= 11 {
			updateItemKeys(num)
		}
	}

	// Initialize item keys on app startup if config exists
	if len(config.Keys) > 0 {
		updateItemKeys(len(config.Keys))
	}

	// Status label
	statusLabel := widget.NewLabel("Configure os campos acima e clique em 'Iniciar'")

	// HWID display (for support purposes)
	hwidLabel := widget.NewLabel("Carregando HWID...")
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

	// Start button (already declared above)
	var stopMonitoring chan bool

	startButton = widget.NewButton("Iniciar Monitoramento", func() {
		if startButton.Text == "Parar Monitoramento" {
			if stopMonitoring != nil {
				close(stopMonitoring)
			}
			statusLabel.SetText("Monitoramento parado.")
			startButton.SetText("Iniciar Monitoramento")
			return
		}

		// Validate email first
		email := emailEntry.Text
		if !IsValidEmail(email) {
			log.Printf("Invalid email. Email used: %s", email)
			statusLabel.SetText("Erro: Digite um email válido")
			return
		}

		// Check if email and subscription are valid (already verified)
		if !isEmailValid || !isSubscriptionValid {
			log.Printf("Email is either invalid or usubscribed. Invalid: %v. Unsubscribed %v", isEmailValid, isSubscriptionValid)
			statusLabel.SetText("Erro: Email deve estar registrado e assinatura ativa")
			return
		}

		// Validate inputs
		numItems, err := strconv.Atoi(numItemsEntry.Text)
		if err != nil || numItems < 1 || numItems > 11 {
			log.Printf("User entered wrong number of items")
			statusLabel.SetText("Erro: Número de items deve ser entre 1 e 11")
			return
		}

		keyShift := keyShiftEntry.Text
		if keyShift != "v" && keyShift != "`" && keyShift != "'" {
			log.Printf("User select wrong key to change bar")
			statusLabel.SetText("Erro: Tecla deve ser *v* ou *`* ou *'*")
			return
		}

		timeClicks, err := strconv.Atoi(timeClicksEntry.Text)
		if err != nil || timeClicks < 0 {
			log.Printf("User selected wrong timing in between changes")
			statusLabel.SetText("Erro: Tempo deve ser um número válido")
			return
		}

		// Collect item keys
		itemKeys := make([]string, numItems)
		for i := 0; i < numItems; i++ {
			if i < len(itemKeyEntries) && itemKeyEntries[i].Text != "" {
				itemKeys[i] = itemKeyEntries[i].Text
			} else {
				log.Printf("User attempted to started before filling all keys. Missing key %d", i+1)
				statusLabel.SetText(fmt.Sprintf("Erro: Digite a tecla para o item %d", i+1))
				return
			}

		}

		// Check subscription before starting monitoring
		statusLabel.SetText("Verificando assinatura...")
		startButton.SetText("Verificando...")

		// Check subscription with retry
		user := ValidadeUser(email, hwid)
		if !user.Active {
			statusLabel.SetText(user.Error)
			startButton.SetText("Iniciar Monitoramento")
			startButton.Disable()
			return
		}

		// Setup configuration
		g.setup = &SetupEquip{
			NumberItems: numItems,
			KeyChange:   keyShift,
			TimeClicks:  timeClicks,
			ItemKeys:    itemKeys,
			CurrentSet:  1,
		}

		startButton.Enable()
		statusLabel.SetText("Assinatura ativa! Monitoramento iniciado. Pressione Q para trocar de set.")
		startButton.SetText("Parar Monitoramento")

		// Start monitoring in a goroutine
		stopMonitoring = make(chan bool)
		log.Printf("Saving configuration into file")
		errConfig := SaveConfig(email, keyShift, timeClicksEntry.Text, itemKeys)
		if errConfig != nil {
			log.Printf("Failed to save config. Error %s", errConfig)
			return
		}
		go g.startMonitoring(statusLabel, stopMonitoring)
	})

	// Initialize button as disabled
	startButton.Disable()

	// Check if saved email is valid and has active subscription
	if config.Email != "" && IsValidEmail(config.Email) {
		go func() {
			fyne.Do(func() {
				user := RegisterEmailWithHWID(config.Email, hwid)
				if user.Active {
					userSub := ValidadeUser(config.Email, hwid)
					if err == nil && userSub.Active {
						emailStatusLabel.SetText("✅ Email registrado e assinatura ativa")
						isEmailValid = true
						isSubscriptionValid = true
						updateButtonState()
					}
				}
			})
		}()
	}

	resetHWIDButton = widget.NewButton("Resetar HWID", func() {
		user := ResetHWID(emailEntry.Text, hwid)
		if user.Active {
			emailStatusLabel.SetText("✅ HWID resetado")
		} else {
			emailStatusLabel.SetText("⚠️ " + user.Error)
		}
	})
	resetHWIDButton.Disable()

	// Form layout
	form := container.NewVBox(
		title,
		instructions,
		widget.NewForm(
			widget.NewFormItem("Email usado na compra do programa:", container.NewVBox(emailEntry, emailStatusLabel)),
			widget.NewFormItem("Quantos items deseja trocar?", numItemsEntry),
			widget.NewFormItem("Tecla para mudar barras de skills:", keyShiftEntry),
			widget.NewFormItem("Tempo entre clicks (em milisegundos):", timeClicksEntry),
		),
		widget.NewLabel("Teclas dos Items:"),
		itemKeysContainer,
		startButton,
		statusLabel,
		widget.NewSeparator(),
		resetHWIDButton,
		hwidLabel,
	)

	scrollContainer := container.NewScroll(form)
	g.window.SetContent(scrollContainer)
	g.window.ShowAndRun()
}

func (g *GuiApp) startMonitoring(statusLabel *widget.Label, stopChan chan bool) {
	evChan := hook.Start()
	defer hook.End()

	for {
		select {
		case <-stopChan:
			return
		case ev := <-evChan:
			if ev.Kind == hook.KeyDown && ev.Keycode == 16 { // Q key
				statusLabel.SetText("Trocando set...")
				ChangeItems(g.setup)
				statusLabel.SetText("Set trocado! Pressione Q novamente para trocar.")
			}
		}
	}
}
