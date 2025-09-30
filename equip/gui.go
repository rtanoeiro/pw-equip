package equip

import (
	"fmt"
	"strconv"
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
	myWindow.Resize(fyne.NewSize(1200, 900))

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

	// Instructions
	instructions := widget.NewRichTextFromMarkdown(`
**Instruções:**
- Deixe 3 barras livres para serem rotacionadas
- Em sua barra principal, deixe suas skills/boticarios como deseja usa-los
- Se deseja iniciar com equipamentos de ataque, na segunda barra deixe os Equipamentos de ataque
- Na ultima barra, deixe os Equipamentos de defesa
- Para trocar de set aperte a tecla Q!
	`)

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Digite seu email usado na compra do programa")

	// Load saved email
	if savedEmail, err := LoadEmail(); err == nil && savedEmail != "" {
		emailEntry.SetText(savedEmail)
	}

	// Email status label
	emailStatusLabel := widget.NewLabel("")

	// Goroutine 1: Handle email validation and local saving
	emailEntry.OnChanged = func(email string) {
		if email == "" {
			emailStatusLabel.SetText("")
			return
		}

		if !IsValidEmail(email) {
			emailStatusLabel.SetText("❌ Email inválido")
			return
		}

		emailStatusLabel.SetText("✅ Email válido - Pressione Enter para registrar")

		if err := SaveEmail(email); err != nil {
			// Handle error silently or log it
			fmt.Printf("Error saving email: %v\n", err)
		}
	}

	// Goroutine 2: Register email with HWID when submitted (Enter pressed)
	emailEntry.OnSubmitted = func(email string) {
		if !IsValidEmail(email) {
			emailStatusLabel.SetText("❌ Email inválido")
			return
		}

		emailStatusLabel.SetText("🔄 Registrando email...")

		hwid, err := GetHWID()
		if err != nil {
			emailStatusLabel.SetText("⚠️ Erro ao obter HWID: " + err.Error())
			return
		}

		errorRegister := RegisterEmailWithHWID(email, hwid)
		if errorRegister != nil {
			emailStatusLabel.SetText("⚠️ Erro no registro: " + errorRegister.Error())
		} else {
			emailStatusLabel.SetText("✅ Email registrado com sucesso")
		}
	}

	// Form fields
	numItemsEntry := widget.NewEntry()
	numItemsEntry.SetPlaceHolder("Digite um número de 1 a 11")

	keyShiftEntry := widget.NewEntry()
	keyShiftEntry.SetPlaceHolder("Digite 'v' ou '`'")

	timeClicksEntry := widget.NewEntry()
	timeClicksEntry.SetPlaceHolder("Tempo em segundos")

	// Dynamic item keys container
	itemKeysContainer := container.NewVBox()
	var itemKeyEntries []*widget.Entry

	// Function to update item keys fields
	updateItemKeys := func(numItems int) {
		itemKeysContainer.RemoveAll()
		itemKeyEntries = make([]*widget.Entry, numItems)

		for i := 0; i < numItems; i++ {
			entry := widget.NewEntry()
			entry.SetPlaceHolder(fmt.Sprintf("Tecla do item %d", i+1))
			itemKeyEntries[i] = entry
			label := widget.NewLabel(fmt.Sprintf("Item %d:", i+1))
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

	// Status label
	statusLabel := widget.NewLabel("Configure os campos acima e clique em 'Iniciar'")

	// HWID display (for support purposes)
	hwidLabel := widget.NewLabel("Carregando HWID...")
	go func() {
		hwid, err := GetHWID()
		if err == nil {
			hwidLabel.SetText(fmt.Sprintf("HWID: %s", hwid))
		} else {
			hwidLabel.SetText("Erro ao obter HWID")
		}
	}()

	// Start button
	var startButton *widget.Button
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
			statusLabel.SetText("Erro: Digite um email válido")
			return
		}

		// Validate inputs
		numItems, err := strconv.Atoi(numItemsEntry.Text)
		if err != nil || numItems < 1 || numItems > 11 {
			statusLabel.SetText("Erro: Número de items deve ser entre 1 e 11")
			return
		}

		keyShift := keyShiftEntry.Text
		if keyShift != "v" && keyShift != "`" {
			statusLabel.SetText("Erro: Tecla deve ser 'v' ou '`'")
			return
		}

		timeClicks, err := strconv.Atoi(timeClicksEntry.Text)
		if err != nil || timeClicks < 0 {
			statusLabel.SetText("Erro: Tempo deve ser um número válido")
			return
		}

		// Collect item keys
		itemKeys := make([]string, numItems)
		for i := 0; i < numItems; i++ {
			if i < len(itemKeyEntries) && itemKeyEntries[i].Text != "" {
				itemKeys[i] = itemKeyEntries[i].Text
			} else {
				statusLabel.SetText(fmt.Sprintf("Erro: Digite a tecla para o item %d", i+1))
				return
			}
		}

		// Check subscription before starting monitoring
		statusLabel.SetText("Verificando assinatura...")
		startButton.SetText("Verificando...")

		go func() {
			// Check subscription with retry
			isActive, err := ValidateSubscriptionWithRetry(3)
			if err != nil {
				statusLabel.SetText(fmt.Sprintf("Erro ao verificar assinatura: %v", err))
				startButton.SetText("Iniciar Monitoramento")
				return
			}

			if !isActive {
				statusLabel.SetText("Assinatura inativa. Entre em contato com o suporte.")
				startButton.SetText("Iniciar Monitoramento")
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

			statusLabel.SetText("Assinatura ativa! Monitoramento iniciado. Pressione Q para trocar de set.")
			startButton.SetText("Parar Monitoramento")

			// Start monitoring in a goroutine
			stopMonitoring = make(chan bool)
			go g.startMonitoring(statusLabel, startButton, stopMonitoring)
		}()
	})

	// Form layout
	form := container.NewVBox(
		title,
		instructions,
		widget.NewForm(
			widget.NewFormItem("Email usado na compra do programa:", container.NewVBox(emailEntry, emailStatusLabel)),
			widget.NewFormItem("Quantos items deseja trocar?", numItemsEntry),
			widget.NewFormItem("Tecla para mudar barras de skills:", keyShiftEntry),
			widget.NewFormItem("Tempo entre clicks (segundos):", timeClicksEntry),
		),
		widget.NewLabel("Teclas dos Items:"),
		itemKeysContainer,
		startButton,
		statusLabel,
		widget.NewSeparator(),
		hwidLabel,
	)

	scrollContainer := container.NewScroll(form)
	g.window.SetContent(scrollContainer)
	g.window.ShowAndRun()
}

func (g *GuiApp) startMonitoring(statusLabel *widget.Label, startButton *widget.Button, stopChan chan bool) {
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
				time.Sleep(time.Duration(g.setup.TimeClicks) * time.Second)
				statusLabel.SetText("Set trocado! Pressione Q novamente para trocar.")
			}
		}
	}
}
