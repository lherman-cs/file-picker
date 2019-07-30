package main

import (
	fp "github.com/lherman-cs/file-picker"
	"github.com/rivo/tview"
)

func main() {
	picker := fp.NewFilePicker(".")

	if err := tview.NewApplication().SetRoot(picker, true).Run(); err != nil {
		panic(err)
	}
}
