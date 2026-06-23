package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gtkm-wallet/internal/api"
)

type NodeView struct {
	rpc    *api.GTKMClient
	window fyne.Window

	versionLabel  *widget.Label
	peersLabel    *widget.Label
	miningLabel   *widget.Label
	hashrateLabel *widget.Label
	blockLabel    *widget.Label

	refreshTicker *time.Ticker
	stopChan      chan struct{}
}

func NewNodeView(rpc *api.GTKMClient, window fyne.Window) *NodeView {
	return &NodeView{rpc: rpc, window: window}
}

func (nv *NodeView) CreateUI() fyne.CanvasObject {
	nv.versionLabel = widget.NewLabel("Loading...")
	nv.peersLabel = widget.NewLabel("Loading...")
	nv.miningLabel = widget.NewLabel("Loading...")
	nv.hashrateLabel = widget.NewLabel("Loading...")
	nv.blockLabel = widget.NewLabel("Loading...")

	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() { nv.refreshOnce() })
	stopBtn := widget.NewButton("Stop Auto", func() { nv.stopPeriodic() })
	startBtn := widget.NewButton("Start Auto", func() { nv.startPeriodic(15 * time.Second) })

	grid := container.NewGridWithColumns(2,
		widget.NewCard("Version", "", nv.versionLabel),
		widget.NewCard("Peers", "", nv.peersLabel),
		widget.NewCard("Mining", "", nv.miningLabel),
		widget.NewCard("Hashrate", "", nv.hashrateLabel),
		widget.NewCard("Current Block", "", nv.blockLabel),
	)

	top := container.NewHBox(refreshBtn, startBtn, stopBtn)
	content := container.NewBorder(top, nil, nil, nil, grid)
	nv.startPeriodic(15 * time.Second)
	return content
}

func (nv *NodeView) refreshOnce() {
	if nv.rpc == nil {
		return
	}
	go func() {
		info, err := nv.rpc.GetNodeInfo()
		if err != nil {
			nv.versionLabel.SetText("error")
			return
		}
		blockNum, _ := nv.rpc.GetBlockNumber()
		nv.versionLabel.SetText(fmt.Sprintf("GTKM %s (net %d)", info.Version, info.NetworkID))
		nv.peersLabel.SetText(fmt.Sprintf("%d", info.PeerCount))
		if info.IsMining {
			nv.miningLabel.SetText("Active")
		} else {
			nv.miningLabel.SetText("Idle")
		}
		nv.hashrateLabel.SetText(fmt.Sprintf("%.2f", info.Hashrate))
		nv.blockLabel.SetText(fmt.Sprintf("#%d", blockNum))
	}()
}

func (nv *NodeView) startPeriodic(interval time.Duration) {
	if nv.refreshTicker != nil {
		return
	}
	nv.refreshTicker = time.NewTicker(interval)
	nv.stopChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-nv.refreshTicker.C:
				nv.refreshOnce()
			case <-nv.stopChan:
				return
			}
		}
	}()
}

func (nv *NodeView) stopPeriodic() {
	if nv.refreshTicker == nil {
		return
	}
	nv.refreshTicker.Stop()
	close(nv.stopChan)
	nv.refreshTicker = nil
	nv.stopChan = nil
}
