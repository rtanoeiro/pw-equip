package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	configPath, errPath := GetConfigPath()
	if errPath != nil {
		log.Fatalf("Failed to get config path: %v", errPath)
	}

	// Create a log file for the application
	logFile, err := os.OpenFile(
		fmt.Sprintf("%s", strings.Replace(configPath, "config.json", "app.log", 1)),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("\n\n\n\nApplication started at %s", time.Now().Format(time.RFC3339))

	title := widget.NewLabel("Bem vindo ao seu auxilio de troca de set")
	title.TextStyle.Bold = true

	hwid, err := GetHWID()
	if err != nil {
		hwid = "Erro ao obter HWID"
	}

	instructions := widget.NewRichTextFromMarkdown(`
**Instruções:**
- Deixe 3 barras livres para serem rotacionadas
- Em sua barra principal, deixe suas skills/boticarios como deseja usa-los
- Se deseja iniciar com equipamentos de ataque, na segunda barra deixe os Equipamentos de ataque
- Na ultima barra, deixe os Equipamentos de defesa
- Para trocar de set clique no botão "Clique para definir tecla" e pressione a tecla desejada!
	`)

	config, errConfig := LoadConfig()
	if errConfig != nil {
		log.Printf("failed to load config file")
		return
	}
	log.Printf("Successfully loaded configuration! Config: %v", config)

	emailEntry := widget.NewEntry()

	// Load saved email on config file
	if config.Email != "" {
		emailEntry.SetText(config.Email)
	} else {
		emailEntry.SetPlaceHolder("Digite seu email usado na compra do programa")
	}

	// Email status label defaulted to empty, it's soon replaced with the email status
	emailStatusLabel := widget.NewLabel("")

	// Creation of the base buttons and states for button state management
	var startButton *widget.Button
	var resetHWIDButton *widget.Button
	var isEmailValid bool = false
	var isSubscriptionValid bool = false
	var isMonitoring bool = false

	// Functions to update button states. We use pointers to button pointers so we can access buttons created after this function is called
	updateButtonState := UpdateButtonState(&startButton, &resetHWIDButton, &isEmailValid, &isSubscriptionValid)

	emailEntry.OnChanged = OnEmailChanged(emailStatusLabel, &isEmailValid, &isSubscriptionValid, updateButtonState)

	emailEntry.OnSubmitted = OnEmailSubmitted(hwid, emailStatusLabel, &isEmailValid, &isSubscriptionValid, updateButtonState)

	// Base config for app usage
	numItemsEntry := widget.NewEntry()
	// In case we have configured keys on config file, auto load them
	if config.Keys != nil {
		numItemsEntry.SetText(fmt.Sprintf("%d", len(config.Keys)))
	} else {
		numItemsEntry.SetPlaceHolder("Digite um número de 1 a 11")
	}

	keyBarShiftEntry := widget.NewEntry()
	if config.BarChangeKey != "" {
		keyBarShiftEntry.SetText(config.BarChangeKey)
	} else {
		keyBarShiftEntry.SetPlaceHolder("Digite 'v' ou '`'")
	}

	// Button to capture key press for changing sets
	// We use a button instead of an entry to avoid conflicts with text input
	var changeSetKeyButton *widget.Button
	changeSetKeyButton = widget.NewButton("Clique para definir tecla", OnChangeSetKeyButtonClicked(&changeSetKeyButton, &config, &isMonitoring))

	// If we have a saved key, show it on the button
	if config.ChangeSetKeyChar != "" {
		changeSetKeyButton.SetText(fmt.Sprintf("Tecla: %s (código: %d)", config.ChangeSetKeyChar, config.ChangeSetKeyCode))
	}

	InBetweenTimeClicksEntry := widget.NewEntry()
	if config.InBetweenTimeClicks > 0 {
		InBetweenTimeClicksEntry.SetText(fmt.Sprintf("%d", config.InBetweenTimeClicks))
	} else {
		InBetweenTimeClicksEntry.SetPlaceHolder("Tempo em milisegundos. Exemplo: 1000 = 1 segundo. 200 = 0.2 segundos, quando menor, mais rapido.")
	}

	// Dynamic item keys fields
	itemKeysContainer := container.NewVBox()
	var itemKeyEntries []*widget.Entry

	// Function to update item keys fields. This is where all items
	updateItemKeys := UpdateItemKeys(&config, len(config.Keys), itemKeysContainer, &itemKeyEntries)

	// Update item keys when number of items changes
	numItemsEntry.OnChanged = UpdateItemKeysOnChanged(numItemsEntry, updateItemKeys)

	// Initialize item keys on app startup if config exists
	if len(config.Keys) > 0 {
		updateItemKeys(len(config.Keys))
	}

	// Status label
	statusLabel := widget.NewLabel("Configure os campos acima e clique em 'Iniciar'")

	// HWID display (for support purposes)
	hwidLabel := widget.NewLabel("Carregando HWID...")
	go GetHWIDLabel(hwidLabel)

	// Start button (already declared above)
	var stopMonitoring chan bool

	startButton = widget.NewButton("Iniciar Monitoramento", g.StartButton(
		emailEntry,
		hwid,
		keyBarShiftEntry,
		InBetweenTimeClicksEntry,
		&isEmailValid,
		&isSubscriptionValid,
		&config,
		&itemKeyEntries,
		numItemsEntry,
		statusLabel,
		&startButton,
		&stopMonitoring,
		&isMonitoring,
	),
	)

	// Initialize button as disabled
	startButton.Disable()

	// Check if saved email is valid and has active subscription
	if config.Email != "" && IsValidEmail(config.Email) {
		go func() {
			fyne.Do(func() {
				user := RegisterEmailWithHWID(config.Email, hwid)
				if user.Active {
					userSub := ValidadeUser(config.Email, hwid)
					if userSub.Active {
						emailStatusLabel.SetText("✅ Email registrado e assinatura ativa")
						isEmailValid = true
						isSubscriptionValid = true
						updateButtonState()
					}
				}
			})
		}()
	}

	resetHWIDButton = widget.NewButton("Resetar HWID", ResetHWIDButton(emailEntry, hwid, emailStatusLabel))
	resetHWIDButton.Disable()

	// Form layout
	form := container.NewVBox(
		title,
		instructions,
		widget.NewForm(
			widget.NewFormItem("Email usado na compra do programa:", container.NewVBox(emailEntry, emailStatusLabel)),
			widget.NewFormItem("Quantos items deseja trocar?", numItemsEntry),
			widget.NewFormItem("Tecla para mudar a barra de skills:", keyBarShiftEntry),
			widget.NewFormItem("Tecla para trocar o set:", changeSetKeyButton),
			widget.NewFormItem("Tempo de click entre items (em milisegundos):", InBetweenTimeClicksEntry),
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

func (g *GuiApp) startMonitoring(statusLabel *widget.Label, stopChan chan bool, keyCode uint16) {
	evChan := hook.Start()
	defer hook.End()

	for {
		select {
		case <-stopChan:
			return
		case ev := <-evChan:
			// Check if the pressed key matches the configured key
			if ev.Kind == hook.KeyDown && ev.Keycode == keyCode {
				fyne.Do(func() {
					statusLabel.SetText("Trocando set...")
					ChangeItems(g.setup)
					statusLabel.SetText("Set trocado! Pressione a tecla novamente para trocar.")
				})
			}
		}
	}
}
