package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type applicationType struct {
	pages         *tview.Pages
	ListShortcuts []rune
}

var app *tview.Application

func (application *applicationType) init() {
	file, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE, 0666)
	check(err)
	log.SetOutput(file)

	app = tview.NewApplication()

	application.pages = tview.NewPages()
	pageMain.build()
	pagePlan.build()

	pageConfirm.build()
	application.registerGlobalShortcuts()
	app.SetFocus(pageMain.Flex)

	if err := app.SetRoot(application.pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {

		panic(err)
	}
}

func (application *applicationType) registerGlobalShortcuts() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			application.ConfirmQuit()
		case tcell.KeyF2:
			application.pages.SwitchToPage("main")
		case tcell.KeyF3:
			application.pages.SwitchToPage("plan")
			app.SetFocus(pagePlan.TreeView)
		default:
			return event

		}
		return nil
	})
}

func (application *applicationType) ConfirmQuit() {
	pageConfirm.show("Are you sure you want to exit?", application.Quit)
}

func (application *applicationType) Quit() {
	app.Stop()
}

func check(err interface{}) {
	if err != nil {
		_, fileName, lineNo, _ := runtime.Caller(1) // Получаем информацию о вызывающем файле
		log.Printf("%s: %d\n", filepath.Base(fileName), lineNo)
		log.Println(err)
		panic(err)
	}
}
