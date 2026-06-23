package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gtkm-wallet/internal/api"
)

type NodeView struct {
	rpc    *api.GTKMClient
	window fyne.Window
	
	infoLabel *widget.Label
	versionLabel *widget.Label
	peersLabel *widget.Label
	miningLabel *widget.Label
	hashrateLabel *widget.Label
	blockLabel *widget.Label
}

func NewNodeView(rpc *api.GTKMClient, window fyne.Window) *NodeView {
	return &NodeView{
		rpc:    rpc,
		window: window,
	}
}

func (nv *NodeView) CreateUI() fyne.CanvasObject {
	nv.infoLabel = widget.NewLabel("Node Information")
	
	versionCard := widget.NewCard("Version", "", widget.NewLabel("Loading..."))
	peersCard := widget.NewCard("Peers", "", widget.NewLabel("Loading..."))
	miningCard := widget.NewCard("Mining", "", widget.NewLabel("Loading..."))
	hashrateCard := widget.NewCard("Hashrate", "", widget.NewLabel("Loading..."))
	blockCard := widget.NewCard("Current Block", "", widget.NewLabel("Loading..."))
	
	// Get references to the labels inside cards
	nv.versionLabel = widget.NewLabel("Loading...")
	nv.peersLabel = widget.NewLabel("Loading...")
	nv.miningLabel = widget.NewLabel("Loading...")
	nv.hashrateLabel = widget.NewLabel("Loading...")
	nv.blockLabel = widget.NewLabel("Loading...")
	
	// Recreate cards with proper labels
	versionCard = widget.NewCard("Version", "", nv.versionLabel)
	peersCard = widget.NewCard("Peers", "", nv.peersLabel)
	miningCard = widget.NewCard("Mining", "", nv.miningLabel)
	hashrateCard = widget.NewCard("Hashrate", "", nv.hashrateLabel)
	blockCard = widget.NewCard("Current Block", "", nv.blockLabel)
	
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		nv.refresh()
	})

	grid := container.NewGridWithColumns(3,
		versionCard,
		peersCard,
		miningCard,
		hashrateCard,
		blockCard,
	)

	content := container.NewBorder(
		container.NewVBox(nv.infoLabel, refreshBtn),
		nil,
		nil,
		nil,
		grid,
	)

	nv.refresh()
	return content
}

func (nv *NodeView) refresh() {
	if nv.rpc == nil {
		return
	}
	
	go func() {
		info, err := nv.rpc.GetNodeInfo()
		if err == nil {
			blockNum, _ := nv.rpc.GetBlockNumber()
			
			nv.versionLabel.SetText(fmt.Sprintf("GTKM v%s", info.Version))
			nv.peersLabel.SetText(fmt.Sprintf("%d peers", info.PeerCount))
			if info.IsMining {
				nv.miningLabel.SetText("✅ Active")
			} else {
				nv.miningLabel.SetText("❌ Idle")
			}
			nv.hashrateLabel.SetText(fmt.Sprintf("%.2f H/s", info.Hashrate))
			nv.blockLabel.SetText(fmt.Sprintf("#%d", blockNum))
		}
	}()
}
