package ui

import (
//	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gtkm-wallet/internal/api"
)

type TransactionView struct {
	rpc    *api.GTKMClient
	window fyne.Window
	transList *widget.List
	transactions []api.Transaction
}

func NewTransactionView(rpc *api.GTKMClient, window fyne.Window) *TransactionView {
	return &TransactionView{
		rpc:    rpc,
		window: window,
	}
}

func (tv *TransactionView) CreateUI() fyne.CanvasObject {
	tv.transList = widget.NewList(
		func() int { return len(tv.transactions) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Tx: "),
				widget.NewLabel("From: "),
				widget.NewLabel("To: "),
				widget.NewLabel("Value: "),
			)
		},
		func(i int, o fyne.CanvasObject) {
			if i >= len(tv.transactions) {
				return
			}
			tx := tv.transactions[i]
			vbox := o.(*fyne.Container)
			if len(vbox.Objects) >= 4 {
				if label, ok := vbox.Objects[0].(*widget.Label); ok {
					label.SetText("Tx: " + tx.Hash.Hex()[:16] + "...")
				}
				if label, ok := vbox.Objects[1].(*widget.Label); ok {
					label.SetText("From: " + tx.From.Hex()[:10] + "...")
				}
				if label, ok := vbox.Objects[2].(*widget.Label); ok {
					label.SetText("To: " + tx.To.Hex()[:10] + "...")
				}
				if label, ok := vbox.Objects[3].(*widget.Label); ok {
					label.SetText("Value: " + tx.Value.String())
				}
			}
		},
	)

	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		tv.refresh()
	})

	detailsCard := widget.NewCard("Transaction Details", "", 
		widget.NewLabel("Select a transaction to view details"),
	)

	content := container.NewBorder(
		refreshBtn,
		nil,
		nil,
		nil,
		container.NewHSplit(
			tv.transList,
			detailsCard,
		),
	)

	tv.refresh()
	return content
}

func (tv *TransactionView) refresh() {
	// TODO: Load actual transactions
	tv.transList.Refresh()
}
