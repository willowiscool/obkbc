package main

import (
	"log"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
)

func main() {
	combs, _, err := readCombs()
	if err != nil {
		log.Fatal("Problem reading keyboard combinations: " + err.Error())
	}
	// Set up main variables, and set up pages
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)
	addForm := tview.NewForm()
	editForm := tview.NewForm()
	pages := tview.NewPages().
		AddPage("edit", editForm, true, true).
		AddPage("add", addForm, true, true).
		AddPage("combs", table, true, true)
	pages.
		SetBorder(true).
		SetTitle("OpenBox Keyboard Customizer")

	//==========FORM==========
	// Set up form and inputs
	// I use separate inputs here because I need to be able to read from each one in the future
	combInput := tview.NewInputField().SetLabel("Keybind")
	combInput.SetFormAttributes(0, tcell.ColorWhite, tcell.ColorBlack, tcell.ColorWhite, tcell.ColorGreen)
	commandInput := tview.NewInputField().SetLabel("Command")
	commandInput.SetFormAttributes(0, tcell.ColorWhite, tcell.ColorBlack, tcell.ColorWhite, tcell.ColorGreen)
	addForm.
		AddFormItem(combInput).
		AddFormItem(commandInput).
		AddButton("Add", func() {
			// Check if a combination for that key already exists
			for _, comb := range combs {
				if comb.key == combInput.GetText() {
					// Open up a modal. Doesn't draw over the existing screen, but takes it up
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
			addComb(Kcommand{
				key: combInput.GetText(),
				command: commandInput.GetText(),
			})
			// Update list of keyboard combinations, and then reset form and go back
			combs = append(combs, Kcommand{
				key: combInput.GetText(),
				command: commandInput.GetText(),
			})
			addCombs(table, combs)
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

	// Separate form for editing combinations
	editCommandInput := tview.NewInputField().SetLabel("Command for ___")
	editCommandInput.SetFormAttributes(0, tcell.ColorWhite, tcell.ColorBlack, tcell.ColorWhite, tcell.ColorGreen)
	editForm.
		AddFormItem(editCommandInput).
		// Set add function to nil because it gets changed in the table func
		AddButton("Edit", nil).
		AddButton("Cancel", func() {
			editCommandInput.SetText("")
			pages.SwitchToPage("combs")
		}).AddButton("Quit", func() {
			app.Stop()
		})

	//==========TABLE==========
	//set headers & add combs
	addHeader(table, 0, "Keybind")
	addHeader(table, 1, "Command")
	addHeader(table, 2, "Edit")
	addHeader(table, 3, "Delete")
	addCombs(table, combs)
	//Handle events
	table.SetDoneFunc(func(key tcell.Key) {
		// key == 'q' shouldn't matter here, but meh
		if key == tcell.KeyEscape || key == 'q' {
			app.Stop()
		}
	}).SetSelectedFunc(func(row, col int) {
		switch col {
			case 2:
				// Column 2 is "Edit". This makes the form command-specific, with the present value in the input
				editCommandInput.
					SetLabel("Command for " + combs[row-1].key).
					SetText(combs[row-1].command)
				// Set the function to edit that *specific* combination
				editForm.GetButton(0).SetSelectedFunc(func() {
					editComb(combs[row-1], Kcommand{
						key: combs[row-1].key,
						command: editCommandInput.GetText(),
					})
					// Once again reset the combinations.
					combs[row-1] = Kcommand{
						key: combs[row-1].key,
						command: editCommandInput.GetText(),
					}
					addCombs(table, combs)
					pages.SwitchToPage("combs")
					editForm.GetButton(0).SetSelectedFunc(nil)
				})
				pages.SwitchToPage("edit")
			case 3:
				// Column 3 is "Delete". This opens a modal to make sure the combination wants to be deleted.
				modal := tview.NewModal().
					SetText("Are you sure you want to delete the combination for " + combs[row-1].key + "?").
					AddButtons([]string{"Yes", "Cancel"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if (buttonLabel == "Yes") {
							deleteComb(combs[row-1].key)
							combs = append(combs[:row-1], combs[row:]...)
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

	err = app.SetRoot(pages, true).Run()
	if err != nil {
		log.Fatal(err)
	}
}

// addHeader adds a header to a table, with a yellow centered text that cannot be selected
func addHeader(table *tview.Table, col int, name string) {
	table.SetCell(0, col, tview.NewTableCell(name).
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetSelectable(false).
		SetExpansion(1))
}

// addCombs adds the keyboard combination using the combs given
func addCombs(table *tview.Table, combs []Kcommand) {
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
