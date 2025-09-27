package equip

import (
	"fmt"
	"log"
	"time"

	"github.com/go-vgo/robotgo"
)

func ClickButton(button string) {
	errorPress := robotgo.KeyPress(button)
	if errorPress != nil {
		log.Println("Erro ao pressionar botao", button)
	}
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
		fmt.Println("Set trocado com sucesso para o segundo")
	}

	if equipSetup.CurrentSet == 2 {
		ClickButton(equipSetup.KeyChange)
		for _, itemToPress := range equipSetup.ItemKeys {
			ClickButton(itemToPress)
			time.Sleep(time.Duration(equipSetup.TimeClicks) * time.Millisecond)
		}
		ClickButton(equipSetup.KeyChange)
		ClickButton(equipSetup.KeyChange)
		fmt.Println("Set trocado com sucesso para o primeiro")
		equipSetup.CurrentSet = 1
	}
}
