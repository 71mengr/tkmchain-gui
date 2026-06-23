package ui

import (
	"fmt"
	"strings"
	"time"
	"math/big"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ethereum/go-ethereum/common"
	"gtkm-wallet/internal/api"
)

type TransactionView struct {
	rpc          *api.GTKMClient
	window       fyne.Window
	transList    *widget.List
	transactions []api.Transaction // in-memory list shown
	detailLabel  *widget.Label
}

// NewTransactionView ...
func NewTransactionView(rpc *api.GTKMClient, window fyne.Window) *TransactionView {
	return &TransactionView{rpc: rpc, window: window}
}

func (tv *TransactionView) CreateUI() fyne.CanvasObject {
	tv.detailLabel = widget.NewLabel("Select a transaction or enter a hash below")

	tv.transList = widget.NewList(
		func() int { return len(tv.transactions) },
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewLabel("Tx:"), widget.NewLabel("From:"), widget.NewLabel("To:"), widget.NewLabel("Value:"))
		},
		func(i int, o fyne.CanvasObject) {
			if i >= len(tv.transactions) {
				return
			}
			tx := tv.transactions[i]
			v := o.(*fyne.Container)
			if len(v.Objects) >= 4 {
				v.Objects[0].(*widget.Label).SetText("Tx: " + tx.Hash.Hex()[:16] + "...")
				v.Objects[1].(*widget.Label).SetText("From: " + tx.From.Hex())
				v.Objects[2].(*widget.Label).SetText("To: " + tx.To.Hex())
				v.Objects[3].(*widget.Label).SetText("Value: " + weiToEthString(tx.Value))
			}
		},
	)

	// Manual lookup
	hashEntry := widget.NewEntry()
	hashEntry.SetPlaceHolder("Enter transaction hash (0x...)")
	lookupBtn := widget.NewButtonWithIcon("Lookup", ui.Icon16, func() {
		h := strings.TrimSpace(hashEntry.Text)
		if h == "" {
			return
		}
		go tv.lookupByHash(h)
	})

	refreshBtn := widget.NewButton("Refresh", func() { tv.refreshList() })

	logo := ui.ImageFromResource(ui.Icon32, 32, 32)
        left := container.NewVBox(logo, refreshBtn, tv.transList)
	right := container.NewVBox(widget.NewLabelWithStyle("Transaction Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), tv.detailLabel, widget.NewSeparator(), container.NewHBox(hashEntry, lookupBtn))

	content := container.NewHSplit(left, right)
	content.Offset = 0.35
	tv.refreshList()
	return content
}

func (tv *TransactionView) refreshList() {
	// Placeholder: if you have a historical indexer or server-side API, populate tv.transactions here
	tv.transList.Refresh()
}

// AddSentTransaction should be called when a tx is sent via the wallet to update the list.
func (tv *TransactionView) AddSentTransaction(tx api.Transaction) {
	tv.transactions = append([]api.Transaction{tx}, tv.transactions...)
	tv.transList.Refresh()
}

// lookupByHash - gets tx from node and shows details
func (tv *TransactionView) lookupByHash(h string) {
	if tv.rpc == nil {
		dialog.ShowError(fmt.Errorf("not connected to node"), tv.window)
		return
	}
	hh := common.HexToHash(h)
	tx, err := tv.rpc.GetTransactionByHash(hh)
	if err != nil {
		dialog.ShowError(err, tv.window)
		return
	}
	if tx == nil {
		tv.detailLabel.SetText("Transaction not found")
		return
	}
	// render details
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Hash: %s\nFrom: %s\nTo: %s\nValue: %s ETH\nGas: %d\nGasPrice: %s\nNonce: %d\nStatus: %d\nBlock: %d\n",
		tx.Hash.Hex(),
		tx.From.Hex(),
		tx.To.Hex(),
		weiToEthString(tx.Value),
		tx.Gas,
		tx.GasPrice.String(),
		tx.Nonce,
		tx.Status,
		tx.BlockNum,
	))
	// show receipt if available
	receipt, err := tv.rpc.GetTransactionReceipt(tx.Hash)
	if err == nil && receipt != nil {
		b.WriteString("\nReceipt:\n")
		for k, v := range receipt {
			b.WriteString(fmt.Sprintf("%s: %v\n", k, v))
		}
	}
	tv.detailLabel.SetText(b.String())
}

// helper to show TKM from wei
func weiToEthString(v *big.Int) string {
	if v == nil {
		return "0"
	}
	f := new(big.Float).Quo(new(big.Float).SetInt(v), new(big.Float).SetFloat64(1e18))
	return f.Text('f', 6)
}
