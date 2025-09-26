package main

import (
	"fmt"
	"time"

	hook "github.com/robotn/gohook"
)

var numItems = "Quantos items voce deseja trocar?"
var numberShift = "Quantas barras o seu V ou ` vao trocar?"

func main() {

	fmt.Printf("Bem vindo ao seu auxilio de troca de set.")
	fmt.Printf("Para que esse programa funciona do jeito esperado, deixe 3 barras livres para serem rotacionadas.")
	fmt.Printf("Em sua barra principal, deixe suas skills/boticarios como deseja usa-los.")
	fmt.Printf("Digamos que voce deseja iniciar com equipamentos de ataque, na segunda barra deixe os Equipamentos de ataque")
	fmt.Printf("Na ultima barra, deixe os Equipamentos de defesa")
	fmt.Printf("Caso deseje iniciar com os de defesa, deixe as barras de forma contraria")
	fmt.Println("Para trocar de set aperte o botao da roda do mouse!")
	fmt.Printf("Iniciando")
	for i := 0; i < 3; i++ {
		fmt.Printf(".")
		time.Sleep(1000)
	}
	numItem := AskQuestionInt(numItems)
	keyShift := AskQuestion(numberShift)

	basicSetup := SetupEquip{
		NumberItems: numItem,
		KeyChange:   keyShift,
		CurrentSet:  1,
	}

	itemKeys := make([]string, numItem)
	for item := 1; item <= numItem; item++ {
		currentItemKey := AskQuestion(fmt.Sprintln("Digite a tecla do item numero ", item))
		itemKeys = append(itemKeys, currentItemKey)
	}
	basicSetup.ItemKeys = itemKeys

	for {
		evChan := hook.Start()
		defer hook.End()
		for ev := range evChan {
			if ev.Kind == hook.WheelDown {
				ChangeItems(basicSetup)
			}
		}
	}
}
