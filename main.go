package main

import (
	"fmt"
	"log"

	hook "github.com/robotn/gohook"
)

var numItemsQuestion = "Quantos items voce deseja trocar?"
var keyShiftQuestion = "Quantas barras o seu V ou ` vao trocar?"

func main() {

	fmt.Println("Bem vindo ao seu auxilio de troca de set.")
	fmt.Println("Para que esse programa funciona do jeito esperado, deixe 3 barras livres para serem rotacionadas.")
	fmt.Println("Em sua barra principal, deixe suas skills/boticarios como deseja usa-los.")
	fmt.Println("Digamos que voce deseja iniciar com equipamentos de ataque, na segunda barra deixe os Equipamentos de ataque")
	fmt.Println("Na ultima barra, deixe os Equipamentos de defesa")
	fmt.Println("Caso deseje iniciar com os de defesa, deixe as barras de forma contraria")
	fmt.Println("Para trocar de set aperte o botao da roda do mouse!")
	fmt.Println()
	fmt.Println()

	numItem := AskQuestionInt(numItemsQuestion)
	errorNumItem := ValidadeNumberEquips(numItem)
	if !errorNumItem {
		log.Fatal("O numero maximo de items permitidos eh 11")
	}

	keyShift := AskQuestion(keyShiftQuestion)
	errorShift := ValidateKeyShift(keyShift)
	if !errorShift {
		log.Fatal("Tecla para mudar de items deve ser 'v' ou '`'")
	}

	basicSetup := SetupEquip{
		NumberItems: numItem,
		KeyChange:   keyShift,
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
		defer hook.End()
		for ev := range evChan {
			if ev.Kind == hook.MouseWheel {
				fmt.Println("Clique na roda do mouse, trocando o set...")
				ChangeItems(&basicSetup)
			}
		}
	}
}
