package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gtkm-wallet/internal/api"
)

type MonitoringView struct {
	rpc    *api.GTKMClient
	window fyne.Window
	monitorList *widget.List
	categories []api.MonitoringCategory
}

func NewMonitoringView(rpc *api.GTKMClient, window fyne.Window) *MonitoringView {
	return &MonitoringView{
		rpc:    rpc,
		window: window,
	}
}

func (mv *MonitoringView) CreateUI() fyne.CanvasObject {
	mv.monitorList = widget.NewList(
		func() int { return len(mv.categories) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Category: "),
				widget.NewLabel("Priority: "),
				widget.NewLabel("Description: "),
				widget.NewLabel("Tasks: "),
			)
		},
		func(i int, o fyne.CanvasObject) {
			if i >= len(mv.categories) {
				return
			}
			cat := mv.categories[i]
			vbox := o.(*fyne.Container)
			if len(vbox.Objects) >= 4 {
				if label, ok := vbox.Objects[0].(*widget.Label); ok {
					label.SetText("Category: " + cat.Name)
				}
				if label, ok := vbox.Objects[1].(*widget.Label); ok {
					label.SetText("Priority: " + cat.Priority)
				}
				if label, ok := vbox.Objects[2].(*widget.Label); ok {
					label.SetText("Description: " + cat.Description)
				}
				if label, ok := vbox.Objects[3].(*widget.Label); ok {
					label.SetText("Tasks: " + fmt.Sprintf("%v", cat.Tasks))
				}
			}
		},
	)

	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		mv.refresh()
	})

	// Status grid
	statusGrid := container.NewGridWithColumns(3,
		widget.NewCard("Network Status", "", widget.NewLabel("Connected")),
		widget.NewCard("King Status", "", widget.NewLabel("Active")),
		widget.NewCard("Mining Status", "", widget.NewLabel("Running")),
	)

	content := container.NewBorder(
		container.NewVBox(statusGrid, refreshBtn),
		nil,
		nil,
		nil,
		container.NewVBox(
			widget.NewLabelWithStyle("Monitoring Responsibilities", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			mv.monitorList,
		),
	)

	mv.refresh()
	return content
}

func (mv *MonitoringView) refresh() {
	if mv.rpc == nil {
		return
	}
	
	go func() {
		categories, err := mv.rpc.GetMonitoringResponsibilities()
		if err == nil {
			mv.categories = categories
			mv.monitorList.Refresh()
		}
	}()
}
