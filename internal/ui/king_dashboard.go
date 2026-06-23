package ui

import (
	"fmt"
	"image/color"
	"math/big"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ethereum/go-ethereum/common"
	"gtkm-wallet/internal/api"
)

type KingDashboard struct {
	rpc    *api.GTKMClient
	window fyne.Window
	
	// UI components
	currentKingLabel  *widget.Label
	mainKingLabel     *widget.Label
	nextKingLabel     *widget.Label
	kingList          *widget.List
	rotationInfoLabel *widget.Label
	balanceLabels     map[common.Address]*widget.Label
	
	// Data
	kingInfo *api.KingInfo
}

func NewKingDashboard(rpc *api.GTKMClient, window fyne.Window) *KingDashboard {
	kd := &KingDashboard{
		rpc:           rpc,
		window:        window,
		balanceLabels: make(map[common.Address]*widget.Label),
	}
	
	// Initialize labels
	kd.currentKingLabel = widget.NewLabel("Loading...")
	kd.mainKingLabel = widget.NewLabel("Loading...")
	kd.nextKingLabel = widget.NewLabel("Loading...")
	kd.rotationInfoLabel = widget.NewLabel("Loading...")
	
	return kd
}

func (kd *KingDashboard) CreateUI() fyne.CanvasObject {
	// Header
	header := canvas.NewText("�� Rotating King Dashboard", color.RGBA{255, 215, 0, 255})
	header.TextSize = 24
	header.TextStyle = fyne.TextStyle{Bold: true}
	
	// Current King Card
	currentKingCard := widget.NewCard("Current King", 
		"The currently active rotating king",
		container.NewVBox(
			widget.NewLabelWithStyle("Address:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			kd.currentKingLabel,
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Status:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel("�� Active"),
		),
	)
	
	// Main King Card
	mainKingCard := widget.NewCard("Main King (10%)", 
		"Permanent king with 10% reward share",
		container.NewVBox(
			widget.NewLabelWithStyle("Address:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			kd.mainKingLabel,
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Balance:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel("Loading..."),
		),
	)
	
	// Next King Card
	nextKingCard := widget.NewCard("Next King", 
		"Upcoming king in rotation",
		container.NewVBox(
			widget.NewLabelWithStyle("Address:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			kd.nextKingLabel,
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Blocks until rotation:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			kd.rotationInfoLabel,
		),
	)
	
	// All Kings List
	kd.kingList = widget.NewList(
		func() int {
			if kd.kingInfo == nil {
				return 0
			}
			return len(kd.kingInfo.AllKings)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("��"),
				widget.NewLabel("Address: "),
				widget.NewLabel(""),
				widget.NewLabel("Balance: "),
				widget.NewLabel(""),
			)
		},
		func(i int, o fyne.CanvasObject) {
			if kd.kingInfo == nil || i >= len(kd.kingInfo.AllKings) {
				return
			}
			hbox := o.(*fyne.Container)
			addr := kd.kingInfo.AllKings[i]
			labels := hbox.Objects
			if len(labels) >= 5 {
				if label, ok := labels[2].(*widget.Label); ok {
					label.SetText(addr.Hex()[:10] + "...")
				}
				if label, ok := labels[4].(*widget.Label); ok {
					label.SetText("Loading...")
				}
			}
		},
	)
	
	allKingsCard := widget.NewCard("All Kings", 
		"Total kings in rotation",
		container.NewVBox(
			kd.kingList,
		),
	)
	
	// Refresh button
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		kd.refresh()
	})
	
	// Top row: Current, Main, Next kings
	topRow := container.NewGridWithColumns(3,
		currentKingCard,
		mainKingCard,
		nextKingCard,
	)
	
	// Main content
	content := container.NewBorder(
		container.NewVBox(header, refreshBtn),
		nil,
		nil,
		nil,
		container.NewVBox(
			topRow,
			allKingsCard,
		),
	)
	
	// Initial refresh
	kd.refresh()
	
	return content
}

func (kd *KingDashboard) refresh() {
	if kd.rpc == nil {
		return
	}
	
	go func() {
		info, err := kd.rpc.GetFullKingInfo()
		if err != nil {
			dialog.ShowError(err, kd.window)
			return
		}
		
		kd.kingInfo = info
		
		// Update UI directly - Fyne widgets are thread-safe for SetText
		if kd.currentKingLabel != nil {
			kd.currentKingLabel.SetText(info.CurrentKing.Hex())
		}
		
		if kd.mainKingLabel != nil {
			kd.mainKingLabel.SetText(info.MainKing.Hex())
		}
		
		if kd.nextKingLabel != nil {
			kd.nextKingLabel.SetText(info.NextKing.Hex())
		}
		
		if kd.rotationInfoLabel != nil {
			blocksUntil := info.RotationInfo.BlocksUntilNext
			kd.rotationInfoLabel.SetText(fmt.Sprintf("%d blocks", blocksUntil))
		}
		
		if kd.kingList != nil {
			kd.kingList.Refresh()
		}
		
		kd.fetchBalances()
	}()
}

func (kd *KingDashboard) fetchBalances() {
	if kd.kingInfo == nil {
		return
	}
	
	for _, addr := range kd.kingInfo.AllKings {
		go func(address common.Address) {
			balance, err := kd.rpc.GetBalance(address)
			if err != nil {
				return
			}
			
			ethBalance := new(big.Float).Quo(
				new(big.Float).SetInt(balance),
				new(big.Float).SetFloat64(1e18),
			)
			
			if label, ok := kd.balanceLabels[address]; ok {
				label.SetText(fmt.Sprintf("%.4f ETH", ethBalance))
			}
		}(addr)
	}
}
