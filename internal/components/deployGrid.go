package components

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/alexanderlesser/sv-cli/types"
	"github.com/alexanderlesser/sv-cli/utils/helpers"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func generateList(app *tview.Application, record types.Record, files []types.File, js bool, selectedFunc func(i int, s1, s2 string, r rune)) tview.List {

	var list = tview.NewList()

	list.SetSelectedFunc(selectedFunc).SetBorder(true)

	list.ShowSecondaryText(false)
	if js {
		list.SetTitle("JS files")
	} else {
		list.SetTitle("CSS files")
	}

	for i, f := range files {
		name := f.Name
		list.AddItem(name, "Some explanatory text", rune(49+i), nil)
	}
	list.AddItem("Quit", "Quits application", 'q', func() {
		app.Stop()
		os.Exit(0)
	})

	return *list
}

func generateEntryTable(data types.Record, deployText *tview.TextView) *tview.Table {
	entryTable := tview.NewTable()

	entryTable.SetSelectable(true, false)
	entryTable.SetEvaluateAllRows(true).SetBorder(true).SetTitle("Deploys")

	for i, j := 0, len(data.Entries)-1; i < j; i, j = i+1, j-1 {
		data.Entries[i], data.Entries[j] = data.Entries[j], data.Entries[i]
	}

	entryTable.SetCell(0, 0, &tview.TableCell{Text: "Name", Color: tcell.ColorYellow, Align: tview.AlignCenter, MaxWidth: 0, Expansion: 3})
	entryTable.SetCell(0, 1, &tview.TableCell{Text: "Date", Color: tcell.ColorYellow, Align: tview.AlignCenter, MaxWidth: 0, Expansion: 10})
	entryTable.SetCell(0, 2, &tview.TableCell{Text: "Time", Color: tcell.ColorYellow, Align: tview.AlignCenter, MaxWidth: 0, Expansion: 10})

	for row, entry := range data.Entries {
		var entryTableColor tcell.Color

		if entry.ErrorWarning {
			entryTableColor = tcell.ColorRed
		} else {
			entryTableColor = tcell.ColorGreen
		}

		entryTable.SetCell(row+1, 0, &tview.TableCell{
			Text:  entry.Name,
			Color: entryTableColor,
			Align: tview.AlignCenter,
		})
		entryTable.SetCell(row+1, 1, &tview.TableCell{
			Text:  entry.Date,
			Color: entryTableColor,
			Align: tview.AlignCenter,
		})
		entryTable.SetCell(row+1, 2, &tview.TableCell{
			Text:  entry.Time,
			Color: entryTableColor,
			Align: tview.AlignCenter,
		})
	}

	return entryTable
}

// Displays a spinner with a text
func spinTitle(app *tview.Application, tv *tview.TextView, startText string, loadText string, action func()) {
	done := make(chan bool)

	//action
	go func() {
		action()
		done <- true
		close(done)
	}()

	// spinner
	go func() {
		spinners := []string{"⠋", "⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		var i int
		for {
			select {
			case _ = <-done:
				app.QueueUpdateDraw(func() {
					tv.SetText(startText)
				})
				return
			case <-time.After(100 * time.Millisecond):
				spin := i % len(spinners)
				app.QueueUpdateDraw(func() {
					tv.SetText("\n" + spinners[spin] + loadText)
				})
				i++
			}
		}
	}()
}

func setDeployText(app *tview.Application, s string, file types.File, record types.Record, resultChan chan types.DeploySuccess, deployText *tview.TextView) {

	deployText.Clear()
	deployText.SetTextAlign(tview.AlignCenter)

	spinTitle(app, deployText, " ", s, func() {

		updatedRecord, err := helpers.DeployFile(record, file)
		if err != nil {
			app.Stop()
			log.Fatal(err)
			os.Exit(1)
		}

		resultChan <- updatedRecord
	})
}

func generateTextView(app *tview.Application) *tview.TextView {

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorder(true).SetTitle("Output")

	return textView
}

func GenerateGrid(app *tview.Application, grid *tview.Grid, record types.Record, files []types.File, js bool, startGulp bool) {
	tabIndex := 0

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	deployText := tview.NewTextView()

	entryTable := generateEntryTable(record, deployText)

	list := generateList(app, record, files, js, func(i int, s1, s2 string, r rune) {
		file := files[i]

		text := " Deploying file " + file.Name

		updateMsg := make(chan bool)
		resultChan := make(chan types.DeploySuccess)
		go func() {
			res := <-resultChan
			var textColor tcell.Color
			var statusText string
			if res.Success {
				statusText = "\nDeploy Successful"
				textColor = tcell.ColorGreen
			} else {
				statusText = "\nDeploy failed"
				textColor = tcell.ColorRed
			}

			var cellColor tcell.Color

			if res.Entry.ErrorWarning {
				cellColor = tcell.ColorRed
			} else {
				cellColor = tcell.ColorGreen
			}

			entryTable.InsertRow(1).SetCell(1, 0, &tview.TableCell{
				Text:  res.Entry.Name,
				Color: cellColor,
				Align: tview.AlignCenter,
			}).SetCell(1, 1, &tview.TableCell{
				Text:  res.Entry.Date,
				Color: cellColor,
				Align: tview.AlignCenter,
			}).SetCell(1, 2, &tview.TableCell{
				Text:  res.Entry.Time,
				Color: cellColor,
				Align: tview.AlignCenter,
			})

			app.QueueUpdateDraw(func() {
				deployText.SetTextColor(textColor).SetText(statusText)
			})

			// Notify the channel to update the message after 10 seconds
			updateMsg <- true
		}()

		setDeployText(app, text, file, record, resultChan, deployText)

		go func() {
			<-updateMsg

			time.Sleep(3 * time.Second)

			app.QueueUpdateDraw(func() {
				deployText.SetTextColor(tcell.ColorWhite).Clear()
			})
		}()
	})

	grid.SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorder(true)

	txtView := generateTextView(app)

	rightMenu := txtView
	leftMenu := list
	header := newPrimitive("\n" + record.Name)
	main := entryTable

	if startGulp {
		go runGulp(app, record.Path, txtView)
	}

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Key input capture for grid
		if event.Key() == tcell.KeyTab {
			if tabIndex == 0 {
				tabIndex = 1
				app.SetFocus(main)
			} else if tabIndex == 1 {

				if startGulp {
					tabIndex = 2
					app.SetFocus(rightMenu)
					app.QueueUpdateDraw(func() {
						rightMenu.SetBorderColor(tcell.ColorGreen)
					})
				} else {
					tabIndex = 0
					app.SetFocus(&leftMenu)
				}
			} else {
				tabIndex = 0
				app.SetFocus(&leftMenu)
			}
			return nil
		}
		return event
	})

	grid.SetRows(3, 0).
		SetColumns(30, 0, 30).
		AddItem(newPrimitive("\nSelect file to deploy"), 0, 0, 1, 1, 1, 1, false).
		AddItem(header, 0, 0, 0, 0, 0, 0, false).
		AddItem(header, 0, 1, 1, 3, 1, 100, false).
		AddItem(deployText, 0, 1, 1, 5, 1, 0, false).
		AddItem(deployText, 0, 4, 1, 2, 1, 100, false)

	grid.AddItem(&leftMenu, 1, 0, 1, 1, 0, 0, true)

	if startGulp {
		grid.AddItem(rightMenu, 1, 4, 1, 2, 0, 130, false)
		grid.AddItem(main, 1, 1, 1, 3, 0, 0, false)
	} else {
		grid.AddItem(main, 1, 1, 1, 5, 0, 0, false)
	}

	app.SetFocus(&leftMenu)
}

func runGulp(app *tview.Application, path string, textView *tview.TextView) {
	// Get the absolute path of the provided directory
	absPath, err := filepath.Abs(path)
	if err != nil {
		app.Stop()
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Change working directory
	err = os.Chdir(absPath)
	if err != nil {
		app.Stop()
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Create a pipe to capture the command output
	pr, pw := io.Pipe()

	// Run "gulp watch" command
	cmd := exec.Command("gulp", "watch")
	cmd.Stdout = pw // Use the write end of the pipe as stdout
	cmd.Stderr = pw // Use the write end of the pipe as stderr

	// Start the command but don't wait for it to finish
	err = cmd.Start()
	if err != nil {
		app.Stop()
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	go func() {
		defer pw.Close() // Close the write end when done reading

		buf := make([]byte, 4096)
		for {
			n, err := pr.Read(buf)
			if n > 0 {
				fmt.Fprintf(textView, "%s ", buf[:n])
			}
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading output:", err)
				}
				// break
			}
		}
	}()

	// Wait for a termination signal
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, os.Interrupt, syscall.SIGTERM)
	<-terminateSignal

	// Terminate the command
	err = cmd.Process.Kill()
	if err != nil {
		fmt.Println("Error:", err)
	}

	message := "\nGulp watch command terminated"
	app.QueueUpdateDraw(func() {
		fmt.Fprintf(textView, "%s ", message)
	})
}
