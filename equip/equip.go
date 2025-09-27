package equip

import (
	"fmt"
	"log"

	hook "github.com/robotn/gohook"
)

func Run() {

	numItemsQuestion := "Quantos items voce deseja trocar?"
	keyShiftQuestion := "Qual a letra do teclado que muda as barras de skills?"
	timeBetweenClicks := "Quanto tempo voce deseja esperar entre os clicks? (em segundos)"

	fmt.Println("Bem vindo ao seu auxilio de troca de set.")
	fmt.Println("Para que esse programa funciona do jeito esperado, deixe 3 barras livres para serem rotacionadas.")
	fmt.Println("Em sua barra principal, deixe suas skills/boticarios como deseja usa-los.")
	fmt.Println("Digamos que voce deseja iniciar com equipamentos de ataque, na segunda barra deixe os Equipamentos de ataque")
	fmt.Println("Na ultima barra, deixe os Equipamentos de defesa")
	fmt.Println("Caso deseje iniciar com os de defesa, deixe as barras de forma contraria")
	fmt.Println("Para trocar de set aperte a letra Q!")
	fmt.Println()
	fmt.Println()

	numItem := AskQuestionInt(numItemsQuestion)
	errorNumItem := ValidadeNumberEquips(numItem)
	if !errorNumItem {
		log.Fatal("O numero maximo de items permitidos eh 11")
	}

	keyShiftLetter := AskQuestion(keyShiftQuestion)
	errorShiftLetter := ValidateKeyShift(keyShiftLetter)
	if !errorShiftLetter {
		log.Fatal("Resposta invalida, tente usar 'v' ou '`'")
	}

	keyShiftNumber := AskQuestion(keyShiftQuestion)
	errorShiftNumber := ValidateKeyShift(keyShiftNumber)
	if !errorShiftNumber {
		log.Fatal("Tecla para mudar de items deve ser 'v' ou '`'")
	}

	timeClicks := AskQuestionInt(timeBetweenClicks)

	basicSetup := SetupEquip{
		NumberItems: numItem,
		KeyChange:   keyShiftLetter,
		TimeClicks:  timeClicks,
		CurrentSet:  1,
	}

	itemKeys := make([]string, numItem)
	for item := 0; item < numItem; item++ {
		currentItemKey := AskQuestion(fmt.Sprintf("Digite a tecla do item numero %d", item+1))
		itemKeys[item] = currentItemKey
	}
	basicSetup.ItemKeys = itemKeys

	for {
		evChan := hook.Start()
		for ev := range evChan {
			if ev.Kind == hook.KeyDown && ev.Keycode == 16 {
				fmt.Println("Clique na roda do mouse, trocando o set...")
				ChangeItems(&basicSetup)
			}
		}
	}
}
