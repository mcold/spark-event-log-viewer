package main

import "github.com/rivo/tview"

type pageHelpType struct {
	helpText string
	*tview.TextView
}

var pageHelp pageHelpType

func (pageHelp *pageHelpType) build() {

	pageHelp.TextView = tview.NewTextView()
	pageHelp.TextView.SetBorder(true)
	pageHelp.TextView.SetBorderPadding(1, 1, 1, 1)
	pageHelp.TextView.SetTitleAlign(tview.AlignCenter)
	pageHelp.TextView.SetRegions(true)
	pageHelp.TextView.SetWordWrap(true)
	pageHelp.TextView.SetWrap(true)
	pageHelp.TextView.SetScrollable(true)

	pageHelp.helpText = `KEYBOARD SHORTCUTS

[green]Global:[-:-:-:-]
[yellow]esc[-:-:-:-]: Quit application.

[green]Main UI:[-:-:-:-]
[yellow]F1[-:-:-:-]: Help
[yellow]F2[-:-:-:-]: Event List
[yellow]F3[-:-:-:-]: Spark Tree Plan
[yellow]F4[-:-:-:-]: Sources
`

	pageHelp.TextView.SetText(pageHelp.helpText)

	application.pages.AddPage("help", pageHelp.TextView, true, false)

}
