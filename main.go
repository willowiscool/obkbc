package main

import (
	"log"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
)

var (
	combs = readCombs()
)

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)
	addForm := tview.NewForm()
	pages := tview.NewPages().
		AddPage("add", addForm, true, true).
		AddPage("combs", table, true, true)
	pages.
		SetBorder(true).
		SetTitle("OpenBox Keyboard Customizer")
	
	//==========TABLE==========
	//set headers & add combs
	addHeader(table, 0, "Keybind")
	addHeader(table, 1, "Command")
	addHeader(table, 2, "Edit")
	addHeader(table, 3, "Delete")
	addCombs(table)
	//Handle events
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape || key == 'q' {
			app.Stop()
		}
	}).SetSelectedFunc(func(row, col int) {
		switch col {
			case 3:
				modal := tview.NewModal().
					SetText("Are you sure you want to delete the combination for " + combs[row-1].key + "?").
					AddButtons([]string{"Yes", "Cancel"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if (buttonLabel == "Yes") {
							deleteComb(combs[row-1])
							table.RemoveRow(row)
						}
						app.SetRoot(pages, true).SetFocus(pages)
					})
				app.SetRoot(modal, false).SetFocus(modal)
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'a' {
			pages.SwitchToPage("add")
			return nil
		}
		if event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return event
	})
	table.Select(1, 2).SetSelectable(true, true)

	//==========FORM==========
	combInput := tview.NewInputField().SetLabel("Keybind")
	combInput.SetFormAttributes(0, tcell.ColorWhite, tcell.ColorBlack, tcell.ColorWhite, tcell.ColorGreen)
	commandInput := tview.NewInputField().SetLabel("Command")
	commandInput.SetFormAttributes(0, tcell.ColorWhite, tcell.ColorBlack, tcell.ColorWhite, tcell.ColorGreen)
	addForm.
		AddFormItem(combInput).
		AddFormItem(commandInput).
		AddButton("Add", func() {
			for _, comb := range combs {
				if comb.key == combInput.GetText() {
					modal := tview.NewModal().
						SetText("A combination for that key already exists!").
						AddButtons([]string{"ok"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							app.SetRoot(pages, true).SetFocus(pages)
						})
					app.SetRoot(modal, false).SetFocus(modal)
					return
				}
			}
			addComb(Keybinding{
				key: combInput.GetText(),
				command: commandInput.GetText(),
			})
			combs = readCombs()
			addCombs(table)
			combInput.SetText("")
			commandInput.SetText("")
			pages.SwitchToPage("combs")
		}).AddButton("Cancel", func() {
			combInput.SetText("")
			commandInput.SetText("")
			pages.SwitchToPage("combs")
		}).AddButton("Quit", func() {
			app.Stop()
		})

	err := app.SetRoot(pages, true).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func addHeader(table *tview.Table, col int, name string) {
	table.SetCell(0, col, tview.NewTableCell(name).
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetSelectable(false).
		SetExpansion(1))
}

func addCombs(table *tview.Table) {
	for i, comb := range combs {
		table.SetCell(i+1, 0, tview.NewTableCell(comb.key).
			SetSelectable(false))
		table.SetCell(i+1, 1, tview.NewTableCell(comb.command).
			SetSelectable(false))
		table.SetCell(i+1, 2, tview.NewTableCell("Edit").
			SetSelectable(true))
		table.SetCell(i+1, 3, tview.NewTableCell("Delete").
			SetSelectable(true))
	}
}
