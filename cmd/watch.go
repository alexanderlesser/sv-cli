/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/alexanderlesser/sv-cli/datastore"
	"github.com/alexanderlesser/sv-cli/internal/components"
	"github.com/alexanderlesser/sv-cli/types"
	"github.com/alexanderlesser/sv-cli/utils/helpers"
	"github.com/fsnotify/fsnotify"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

type LastChanges struct {
	Name     string
	fullPath string
	content  []byte
}

var fileContents []LastChanges

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch [command] Watches selected directory for changes",
	Long: `
	Watch [command] Watches selected directory for changes.
	When it detects changes it deploys the files.
	
	For more information about a specific command run record [command] -h`,
	Run: func(cmd *cobra.Command, args []string) {

		records, err := datastore.Load()
		if err != nil {
			fmt.Println("Cannot load records: ", err)
			os.Exit(1)
		}

		app := tview.NewApplication()
		pages := tview.NewPages()

		textView := tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetChangedFunc(func() {
				app.Draw()
			})

		textView.SetBorder(true)

		table := components.DisplayRecordsTable(records, func(k tcell.Key) {
			if k == tcell.KeyEscape {
				app.Stop()
			}
		}, func(row, column int) {
			record := records[row-1]

			pages.AddPage("view", textView, true, false)
			pages.SwitchToPage("view")

			go startGulp(app, textView, record)

		})

		pages.AddPage("table", table, true, true)

		if err := app.SetRoot(pages, true).Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startGulp(app *tview.Application, textView *tview.TextView, record types.Record) {
	// Get the absolute path of the provided directory
	absPath, err := filepath.Abs(record.Path)
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
	var initialStartup = true
	go func() {
		defer pw.Close() // Close the write end when done reading

		buf := make([]byte, 4096)
		for {
			n, err := pr.Read(buf)
			if n > 0 {
				if initialStartup {
					fmt.Fprintf(textView, "%s ", buf[:n])
				}

				if strings.Contains(string(buf[:n]), "Serving files from: assets/") {
					initialStartup = false
					go watcher(app, textView, record)
				}

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

func watcher(app *tview.Application, textView *tview.TextView, record types.Record) {
	fmt.Fprintf(textView, "%s ", "--------------------------------------\n")
	fmt.Fprintf(textView, "%s ", "Starting watcher...")

	path, err := helpers.GetCSSPath(record.Path)

	if err != nil {
		app.Stop()
		fmt.Println(err)
		os.Exit(1)
	}

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) {
					name := filepath.Base(event.Name)
					var file types.File
					file.Name = name
					file.Path = event.Name

					// Read file content
					content, err := os.ReadFile(file.Path)
					if err != nil {
						app.Stop()
						fmt.Println(err)
						os.Exit(1)
					}

					isCss := helpers.IsCssFile(file)
					if isCss {
						e := checkContentMatch(file, content)

						if !e {
							done := make(chan struct{})

							go func() {
								updatedRecord, err := helpers.DeployFile(record, file)
								if err != nil {
									app.Stop()
									log.Fatal(err)
									os.Exit(1)
								}

								now := time.Now()

								if updatedRecord.Success {

									message := "\nSave occcured:" + now.Format("15:04:05") + "\nmodified file: " + name + "\n"
									app.QueueUpdateDraw(func() {
										fmt.Fprintf(textView, "%s ", message)
									})
								} else {
									message := "\nFailed to deploy file: " + name + "\n"
									app.QueueUpdateDraw(func() {
										fmt.Fprintf(textView, "%s ", message)
									})
								}

								// Signal that DeployFile is done
								done <- struct{}{}
							}()
						}
					}

					// deployFile(record, cssFile)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func checkContentMatch(f types.File, c []byte) bool {
	var equal bool
	for i, content := range fileContents {
		if f.Name == content.Name {
			equal = bytes.Equal(content.content, c)
			if !equal {
				fileContents[i].content = c // Update with the new content 'c'
			}
			break
		}
	}

	// If the file wasn't found in fileContents, add it
	if !equal {
		fileContents = append(fileContents, LastChanges{
			Name:     f.Name,
			fullPath: f.Path,
			content:  c,
		})
	}

	return equal
}
