package ui

import (
//	"fmt"
	"math/big"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ethereum/go-ethereum/common"
	"gtkm-wallet/internal/api"
)

type WalletView struct {
	rpc    *api.GTKMClient
	window fyne.Window
	
	balanceLabel *widget.Label
	addressEntry *widget.Entry
	toEntry      *widget.Entry
	amountEntry  *widget.Entry
	sendBtn      *widget.Button
	accountList  *widget.List
	accounts     []common.Address
}

func NewWalletView(rpc *api.GTKMClient, window fyne.Window) *WalletView {
	return &WalletView{
		rpc:    rpc,
		window: window,
	}
}

func (wv *WalletView) CreateUI() fyne.CanvasObject {
	wv.balanceLabel = widget.NewLabel("Balance: 0 ETH")
	
	// Account list
	wv.accountList = widget.NewList(
		func() int { return len(wv.accounts) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("��"),
				widget.NewLabel("Account: "),
				widget.NewLabel(""),
				widget.NewLabel("Balance: "),
				widget.NewLabel(""),
			)
		},
		func(i int, o fyne.CanvasObject) {
			if i >= len(wv.accounts) {
				return
			}
			hbox := o.(*fyne.Container)
			if len(hbox.Objects) >= 5 {
				if label, ok := hbox.Objects[2].(*widget.Label); ok {
					label.SetText(wv.accounts[i].Hex()[:10] + "...")
				}
			}
		},
	)

	// Balance check
	wv.addressEntry = widget.NewEntry()
	wv.addressEntry.SetPlaceHolder("Enter address to check balance")
	
	checkBtn := widget.NewButton("Check Balance", func() {
		wv.checkBalance()
	})

	// Send transaction
	wv.toEntry = widget.NewEntry()
	wv.toEntry.SetPlaceHolder("Recipient address")
	wv.amountEntry = widget.NewEntry()
	wv.amountEntry.SetPlaceHolder("Amount in ETH")
	wv.sendBtn = widget.NewButtonWithIcon("Send", theme.ConfirmIcon(), func() {
		wv.sendTransaction()
	})

	// Refresh button
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		wv.refreshAccounts()
	})

	sendForm := container.NewVBox(
		widget.NewLabelWithStyle("Send Transaction", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		wv.toEntry,
		wv.amountEntry,
		wv.sendBtn,
	)

	content := container.NewBorder(
		container.NewVBox(
			refreshBtn,
			widget.NewLabelWithStyle("Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			wv.accountList,
			widget.NewSeparator(),
			container.NewHBox(wv.addressEntry, checkBtn, wv.balanceLabel),
		),
		nil,
		nil,
		nil,
		sendForm,
	)

	wv.refreshAccounts()
	return content
}

func (wv *WalletView) refreshAccounts() {
	if wv.rpc == nil {
		return
	}
	
	go func() {
		// Try to get accounts
		// For now, use some default accounts or fetch from node
		wv.accountList.Refresh()
	}()
}

func (wv *WalletView) checkBalance() {
	if wv.rpc == nil || wv.addressEntry.Text == "" {
		return
	}
	
	addr := common.HexToAddress(wv.addressEntry.Text)
	go func() {
		balance, err := wv.rpc.GetBalance(addr)
		if err != nil {
			dialog.ShowError(err, wv.window)
			return
		}
		
		ethBalance := new(big.Float).Quo(
			new(big.Float).SetInt(balance),
			new(big.Float).SetFloat64(1e18),
		)
		
		wv.balanceLabel.SetText("Balance: " + ethBalance.Text('f', 4) + " ETH")
	}()
}

func (wv *WalletView) sendTransaction() {
	if wv.rpc == nil || wv.toEntry.Text == "" || wv.amountEntry.Text == "" {
		dialog.ShowInformation("Error", "Please fill in all fields", wv.window)
		return
	}
	
	// TODO: Implement send transaction
	dialog.ShowInformation("Info", "Transaction sending will be implemented", wv.window)
}
