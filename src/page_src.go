package main

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageSrcType struct {
	*tview.Flex
	*tview.TextArea
	*tview.Button
}

var pageSrc pageSrcType

func (pageSrc *pageSrcType) build() {

	pageSrc.TextArea = tview.NewTextArea()
	pageSrc.TextArea.SetTitle("sources").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue)

	pageSrc.Button = tview.NewButton("Copy")

	pageSrc.Button.SetBorderColor(tcell.ColorBlue)
	pageSrc.Button.SetSelectedFunc(copySrc)

	pageSrc.Flex = tview.NewFlex().
		AddItem(pageSrc.TextArea, 0, 10, true).
		AddItem(pageSrc.Button, 0, 1, false)

	pageSrc.Flex.SetDirection(tview.FlexRow).
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	application.pages.AddPage("src", pageSrc.Flex, true, false)
}

func copySrc() {
	err := clipboard.WriteAll(pageSrc.TextArea.GetText())
	check(err)
}
