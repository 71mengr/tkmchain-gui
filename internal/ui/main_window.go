package ui

import (
	"errors"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gtkm-wallet/internal/api"
	"gtkm-wallet/internal/config"
)

type MainUI struct {
	rpc     *api.GTKMClient
	config  *config.Config
	window  fyne.Window
	
	// Tabs
	kingTab    *container.TabItem
	walletTab  *container.TabItem
	transTab   *container.TabItem
	monitorTab *container.TabItem
	nodeTab    *container.TabItem
	
	// Status
	statusLabel *widget.Label
	statusIcon  *canvas.Circle
	
	// Refresh timer
	refreshTicker *time.Ticker
}

func NewMainUI(rpc *api.GTKMClient, cfg *config.Config, window fyne.Window) *MainUI {
	ui := &MainUI{
		rpc:    rpc,
		config: cfg,
		window: window,
	}
	
	ui.setupTabs()
	return ui
}

func (ui *MainUI) SetupWindow(window fyne.Window) {
	// Status bar
	ui.statusLabel = widget.NewLabel("Connecting to gtkm node: " + ui.config.NodeURL)
	ui.statusIcon = canvas.NewCircle(color.RGBA{255, 255, 0, 255}) // Yellow for connecting
	
	// Status container
	statusContainer := container.NewHBox(
		ui.statusIcon,
		ui.statusLabel,
		layout.NewSpacer(),
		widget.NewLabel("GTKM v"+api.Version),
	)
	
	// Toolbar
	toolbar := ui.createToolbar()
	
	// Main tabs
	tabs := container.NewAppTabs(
		ui.kingTab,
		ui.walletTab,
		ui.transTab,
		ui.monitorTab,
		ui.nodeTab,
	)
	tabs.SetTabLocation(container.TabLocationTop)
	
	// Main container
	content := container.NewBorder(
		toolbar,
		statusContainer,
		nil,
		nil,
		tabs,
	)
	
	window.SetContent(content)
	
	// Initial connection check and refresh
	go ui.checkConnection()
	go ui.startAutoRefresh()
}

func (ui *MainUI) createToolbar() fyne.CanvasObject {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			ui.refreshAll()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HomeIcon(), func() {
			ui.showAboutDialog()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			ui.showSettingsDialog()
		}),
	)
}

func (ui *MainUI) setupTabs() {
	// Create the dashboard
	dashboard := NewKingDashboard(ui.rpc, ui.window)
	ui.kingTab = container.NewTabItem("�� Kings", dashboard.CreateUI())
	
	// Wallet tab
	wallet := NewWalletView(ui.rpc, ui.window)
	ui.walletTab = container.NewTabItem("�� Wallet", wallet.CreateUI())
	
	// Transactions tab
	transactions := NewTransactionView(ui.rpc, ui.window)
	ui.transTab = container.NewTabItem("�� Transactions", transactions.CreateUI())
	
	// Monitoring tab
	monitoring := NewMonitoringView(ui.rpc, ui.window)
	ui.monitorTab = container.NewTabItem("�� Monitoring", monitoring.CreateUI())
	
	// Node info tab
	node := NewNodeView(ui.rpc, ui.window)
	ui.nodeTab = container.NewTabItem("��️ Node", node.CreateUI())
}

func (ui *MainUI) checkConnection() {
	if ui.rpc == nil {
		ui.updateStatus("❌ Not connected to gtkm node", color.RGBA{255, 0, 0, 255})
		return
	}
	
	err := ui.rpc.CheckConnection()
	if err != nil {
		ui.updateStatus("❌ Connection failed: "+err.Error(), color.RGBA{255, 0, 0, 255})
		dialog.ShowError(
			errors.New("Failed to connect to gtkm node.\n\nMake sure gtkm is running with RPC enabled.\n\nRPC URL: "+ui.config.NodeURL),
			ui.window,
		)
	} else {
		ui.updateStatus("✅ Connected to gtkm node on port "+ui.config.NodeURL, color.RGBA{0, 255, 0, 255})
		ui.refreshAll()
	}
}

func (ui *MainUI) updateStatus(text string, clr color.Color) {
	ui.statusLabel.SetText(text)
	ui.statusIcon.FillColor = clr
	ui.statusIcon.Refresh()
}

func (ui *MainUI) startAutoRefresh() {
	interval := time.Duration(ui.config.AutoRefresh) * time.Second
	ui.refreshTicker = time.NewTicker(interval)
	
	for range ui.refreshTicker.C {
		if ui.rpc != nil {
			// Check if connected by calling a method
			err := ui.rpc.CheckConnection()
			if err == nil {
				ui.refreshAll()
			}
		}
	}
}

func (ui *MainUI) refreshAll() {
	// Refresh each tab
	// This will be handled by each view's refresh method
}

func (ui *MainUI) showAboutDialog() {
	dialog.ShowInformation("About GTKM Wallet",
		"GTKM Wallet - Rotating King Edition\n\n"+
		"Connected to: "+ui.config.NodeURL+"\n"+
		"Network ID: 8979\n\n"+
		"Features:\n"+
		"• Rotating King Consensus\n"+
		"• Main King (10% rewards)\n"+
		"• Rotating King (40% rewards)\n"+
		"• Miner (50% rewards)\n\n"+
		"Version: 1.0.0\n"+
		"Built with Fyne and Go",
		ui.window,
	)
}

func (ui *MainUI) showSettingsDialog() {
	// Settings dialog
	nodeEntry := widget.NewEntry()
	nodeEntry.SetText(ui.config.NodeURL)
	nodeEntry.SetPlaceHolder("Node URL (e.g., http://localhost:8545)")
	
	refreshEntry := widget.NewEntry()
	refreshEntry.SetText(fmt.Sprintf("%d", ui.config.AutoRefresh))
	refreshEntry.SetPlaceHolder("Refresh interval (seconds)")
	
	content := container.NewVBox(
		widget.NewLabel("Node Settings"),
		nodeEntry,
		widget.NewLabel("Refresh Settings"),
		refreshEntry,
	)
	
	dialog.ShowCustomConfirm("Settings", "Save", "Cancel", content,
		func(save bool) {
			if save {
				// Save settings
				ui.config.NodeURL = nodeEntry.Text
				// TODO: Save config to file
				ui.updateStatus("Settings saved", color.RGBA{0, 255, 0, 255})
				// Reconnect
				go ui.checkConnection()
			}
		},
		ui.window,
	)
}
