package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
)

func ClickButton(button string) {
	errorPress := robotgo.KeyPress(button)
	if errorPress != nil {
		log.Println("Erro ao pressionar botao", button)
	}
}

func AskQuestion(question string) string {
	fmt.Println(question)
	var variableToRead string
	_, errorQuestion := fmt.Scanln(&variableToRead)
	if errorQuestion != nil {
		fmt.Println("Erro ao ler resposta. Tente novamente.")
		return AskQuestion(question)
	}
	fmt.Println("Resposta: ", variableToRead)
	return variableToRead
}

func AskQuestionInt(question string) int {
	variableToRead := AskQuestion(question)
	intVariable, errConvert := StringToInt(variableToRead)
	if errConvert != nil {
		return 0
	}
	return intVariable

}

func StringToInt(value string) (int, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error converting string to int: %v", err)
		return 0, err
	}
	return intValue, nil
}

func ValidadeNumberEquips(numEquips int) bool {
	if numEquips <= 11 && numEquips >= 1 {
		return true
	}
	return false
}

func ValidateKeyShift(keyShift string) bool {
	if keyShift == "v" || keyShift == "`" {
		return true
	}
	return false
}

func ChangeItems(equipSetup *SetupEquip) {
	if equipSetup.CurrentSet == 1 {
		ClickButton(equipSetup.KeyChange)
		ClickButton(equipSetup.KeyChange)
		for _, itemToPress := range equipSetup.ItemKeys {
			ClickButton(itemToPress)
			time.Sleep(50)
		}
		ClickButton(equipSetup.KeyChange)
		equipSetup.CurrentSet = 2
		fmt.Println("Set trocado com sucesso para o segundo")
	}

	if equipSetup.CurrentSet == 2 {
		ClickButton(equipSetup.KeyChange)
		for _, itemToPress := range equipSetup.ItemKeys {
			ClickButton(itemToPress)
		}
		ClickButton(equipSetup.KeyChange)
		ClickButton(equipSetup.KeyChange)
		fmt.Println("Set trocado com sucesso para o primeiro")
		equipSetup.CurrentSet = 2
	}
}
