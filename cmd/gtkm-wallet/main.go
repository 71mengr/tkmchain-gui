package main

import (
	"errors"
	"fmt"
	"log"
	"time"

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

	// Load configuration
	cfg := config.Load()
	log.Printf("GTKM Wallet starting...")
	log.Printf("Connecting to gtkm node at %s", cfg.NodeURL)

	// Try to connect to gtkm node with retries and timeout from config
	dialTimeout := time.Duration(cfg.DialTimeoutSeconds) * time.Second
	retryInterval := time.Duration(cfg.RetryIntervalSeconds) * time.Second

	rpcClient, err := api.NewClientWithRetry(cfg.NodeURL, dialTimeout, cfg.RetryCount, retryInterval)
	if err != nil {
		// Show error dialog with details; allow user to continue in offline mode
		msg := fmt.Sprintf("Failed to connect to gtkm node at %s\n\nReason: %v\n\nMake sure gtkm is running with RPC enabled.", cfg.NodeURL, err)
		dialog.ShowError(errors.New(msg), myWindow)
		log.Printf("Failed to connect to gtkm node: %v", err)
		// rpcClient == nil => the UI should defend against nil client (offline mode)
	} else {
		log.Printf("Connected to gtkm node. chainID=%s", rpcClient.ChainID.String())
		// Keep client for UI usage and make sure it will be closed when app exits.
		defer rpcClient.Close()
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
