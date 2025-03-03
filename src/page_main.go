package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageMainType struct {
	*tview.Flex
	*tview.List
	flDetails   *tview.Flex
	Description *tview.TextView
	Details     *tview.TextView
	Events      []Event
}

var pageMain pageMainType

func (pageMain *pageMainType) build() {

	pageMain.Events = get_events()

	pageMain.List = tview.NewList()
	pageMain.List.SetTitle("events").
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageMain.Details = tview.NewTextView()
	pageMain.Details.SetTitle("details").
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageMain.Description = tview.NewTextView()
	pageMain.Description.SetTitle("description").
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageMain.flDetails = tview.NewFlex().SetDirection(tview.FlexRow)
	pageMain.flDetails.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageMain.flDetails.
		AddItem(pageMain.Description, 0, 1, false).
		AddItem(pageMain.Details, 0, 5, false)

	pageMain.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageMain.List, 0, 1, true).
		AddItem(pageMain.flDetails, 0, 1, false)

	for _, ev := range pageMain.Events {
		pageMain.List.AddItem(ev.EventName, "", rune(0), func() {
			pageMain.Details.SetText(pageMain.Events[pageMain.List.GetCurrentItem()].Details)
			pageMain.Description.SetText(pageMain.Events[pageMain.List.GetCurrentItem()].Description)
		})
	}
	application.pages.AddPage("main", pageMain.Flex, true, true)
}
