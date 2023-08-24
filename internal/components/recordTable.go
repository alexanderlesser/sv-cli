package components

import (
	"strconv"

	"github.com/alexanderlesser/sv-cli/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DisplayRecordsTable(data []types.Record, setDoneFunc func(tcell.Key), setSelectedFunc func(row int, column int)) *tview.Table {

	table := tview.NewTable().SetSelectable(true, false)

	// Create table headers
	table.SetCell(0, 0, &tview.TableCell{Text: "ID", Color: tcell.ColorYellow, Align: tview.AlignLeft, MaxWidth: 0, Expansion: 3})
	table.SetCell(0, 1, &tview.TableCell{Text: "Name", Color: tcell.ColorYellow, Align: tview.AlignLeft, MaxWidth: 0, Expansion: 10})
	table.SetCell(0, 2, &tview.TableCell{Text: "Domain", Color: tcell.ColorYellow, Align: tview.AlignLeft, MaxWidth: 0, Expansion: 10})
	table.SetCell(0, 3, &tview.TableCell{Text: "Username", Color: tcell.ColorYellow, Align: tview.AlignLeft, MaxWidth: 0, Expansion: 10})
	table.SetCell(0, 4, &tview.TableCell{Text: "Path", Color: tcell.ColorYellow, Align: tview.AlignLeft, MaxWidth: 0, Expansion: 10})

	// Iterate through the data entries and populate the table rows
	for row, entry := range data {

		table.SetCell(row+1, 0, &tview.TableCell{
			Text:      strconv.Itoa(int(entry.ID)),
			Color:     tcell.ColorWhite,
			Align:     tview.AlignLeft,
			Expansion: 3,
		})
		table.SetCell(row+1, 1, &tview.TableCell{
			Text:      entry.Name,
			Color:     tcell.ColorWhite,
			Align:     tview.AlignLeft,
			Expansion: 10,
		})
		table.SetCell(row+1, 2, &tview.TableCell{
			Text:      entry.Domain,
			Color:     tcell.ColorWhite,
			Align:     tview.AlignLeft,
			Expansion: 10,
		})
		table.SetCell(row+1, 3, &tview.TableCell{
			Text:      entry.Username,
			Color:     tcell.ColorWhite,
			Align:     tview.AlignLeft,
			Expansion: 10,
		})
		table.SetCell(row+1, 4, &tview.TableCell{
			Text:      entry.Path,
			Color:     tcell.ColorWhite,
			Align:     tview.AlignLeft,
			Expansion: 10,
		})
	}

	if setDoneFunc != nil {
		table.SetDoneFunc(setDoneFunc)
	}

	if setSelectedFunc != nil {
		table.SetSelectedFunc(setSelectedFunc)
	}

	return table
}
