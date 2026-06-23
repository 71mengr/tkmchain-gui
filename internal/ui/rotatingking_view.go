package ui

import (
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gtkm-wallet/internal/api"
)

// RotatingKingView provides RK control panel: add/list/stats/status
type RotatingKingView struct {
	rpc       *api.GTKMClient
	window    fyne.Window
	output    *widget.MultiLineEntry
	addParams *widget.Entry
}

func NewRotatingKingView(rpc *api.GTKMClient, window fyne.Window) *RotatingKingView {
	return &RotatingKingView{rpc: rpc, window: window}
}

func (rv *RotatingKingView) CreateUI() fyne.CanvasObject {
	rv.output = widget.NewMultiLineEntry()
	rv.output.SetReadOnly(true)
	rv.addParams = widget.NewEntry()
	rv.addParams.SetPlaceHolder(`JSON params for rk_add, e.g. {"addr":"0x...","meta":"..."} `)

	addBtn := widget.NewButton("rk_add", func() { rv.doAdd() })
	listBtn := widget.NewButton("rk_list", func() { rv.doList() })
	statsBtn := widget.NewButton("rk_stats", func() { rv.doStats() })
	statusBtn := widget.NewButton("rk_status", func() { rv.doStatus() })

	btns := container.NewHBox(addBtn, listBtn, statsBtn, statusBtn)
	logo := ui.ImageFromResource(ui.Icon48, 48, 48)
        title := container.NewHBox(logo, widget.NewLabelWithStyle("RotatingKing Control", fyne.TextAlignLeading, fyne.TextStyle{Bold:true}))
        content := container.NewBorder(title, nil, nil, nil, container.NewVBox(rv.addParams, btns, rv.output))

	return content
}

func (rv *RotatingKingView) doAdd() {
	if rv.rpc == nil {
		rv.output.SetText("not connected")
		return
	}
	paramStr := rv.addParams.Text
	if paramStr == "" {
		rv.output.SetText("please enter params JSON")
		return
	}
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(paramStr), &params); err != nil {
		rv.output.SetText(fmt.Sprintf("invalid JSON: %v", err))
		return
	}
	res, err := rv.rpc.RKAdd(params)
	if err != nil {
		rv.output.SetText(fmt.Sprintf("rk_add error: %v", err))
		return
	}
	rv.output.SetText(prettyJSON(res))
}

func (rv *RotatingKingView) doList() {
	if rv.rpc == nil {
		rv.output.SetText("not connected")
		return
	}
	res, err := rv.rpc.RKList()
	if err != nil {
		rv.output.SetText(fmt.Sprintf("rk_list error: %v", err))
		return
	}
	rv.output.SetText(prettyJSON(res))
}

func (rv *RotatingKingView) doStats() {
	if rv.rpc == nil {
		rv.output.SetText("not connected")
		return
	}
	res, err := rv.rpc.RKStats()
	if err != nil {
		rv.output.SetText(fmt.Sprintf("rk_stats error: %v", err))
		return
	}
	rv.output.SetText(prettyJSON(res))
}

func (rv *RotatingKingView) doStatus() {
	if rv.rpc == nil {
		rv.output.SetText("not connected")
		return
	}
	res, err := rv.rpc.RKStatus()
	if err != nil {
		rv.output.SetText(fmt.Sprintf("rk_status error: %v", err))
		return
	}
	rv.output.SetText(prettyJSON(res))
}

func prettyJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
