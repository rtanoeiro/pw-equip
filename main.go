package main

import (
	"flag"
	"pw-equip-change/equip"
)

func main() {
	// Command line flag to choose between GUI and console mode
	useGUI := flag.Bool("gui", true, "Use GUI interface (default: true)")
	flag.Parse()

	if *useGUI {
		// Run GUI version
		app := equip.NewGuiApp()
		app.RunGUI()
	} else {
		// Run console version
		equip.Run()
	}
}
