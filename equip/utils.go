package equip

import (
	"time"

	"github.com/go-vgo/robotgo"
)

func ClickButton(button string) {
	robotgo.KeyPress(button)
	// Note: Errors are silently ignored to prevent console output in GUI mode
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
	}

	if equipSetup.CurrentSet == 2 {
		ClickButton(equipSetup.KeyChange)
		for _, itemToPress := range equipSetup.ItemKeys {
			ClickButton(itemToPress)
			time.Sleep(time.Duration(equipSetup.TimeClicks) * time.Millisecond)
		}
		ClickButton(equipSetup.KeyChange)
		ClickButton(equipSetup.KeyChange)
		equipSetup.CurrentSet = 1
	}
}
