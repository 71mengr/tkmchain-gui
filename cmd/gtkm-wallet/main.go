package main

import (
	"errors"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"gtkm-wallet/internal/api"
	"gtkm-wallet/internal/config"
	"gtkm-wallet/internal/ui"
)

func main() {
	// Create application
	myApp := app.NewWithID("com.gtkm.wallet")
	myWindow := myApp.NewWindow("GTKM Wallet - Rotating King Edition")

	// Set icon (optional - skip if file doesn't exist)
	// if icon, err := fyne.LoadResourceFromPath("assets/icons/gtkm.png"); err == nil {
	// 	myWindow.SetIcon(icon)
	// }

	// Load configuration
	cfg := config.Load()
	log.Printf("GTKM Wallet starting...")
	log.Printf("Connecting to gtkm node at %s", cfg.NodeURL)

	// Try to connect to gtkm node
	rpcClient, err := api.NewGTKMClient(cfg.NodeURL)
	if err != nil {
		// Show error dialog
		dialog.ShowError(
			errors.New("Failed to connect to gtkm node.\n\nMake sure gtkm is running with RPC enabled.\n\nRPC URL: "+cfg.NodeURL),
			myWindow,
		)
		log.Printf("Failed to connect to gtkm node: %v", err)
		// Continue with offline mode (rpcClient will be nil)
	}

	// Create main UI
	mainUI := ui.NewMainUI(rpcClient, cfg, myWindow)

	// Setup window
	mainUI.SetupWindow(myWindow)

	// Set window size
	myWindow.Resize(fyne.NewSize(1400, 900))
	myWindow.CenterOnScreen()

	// Show and run
	myWindow.ShowAndRun()
}
