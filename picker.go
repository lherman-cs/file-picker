package picker

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type nodeRef struct {
	FullPath string
	Info     os.FileInfo
}

// FilePicker ...
type FilePicker struct {
	*tview.TreeView
	doneHandler   func(selected *os.File) error
	cancelHandler func()
}

// NewFilePicker is a specialized modal to pick a file
func NewFilePicker(rootDir string) *FilePicker {
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	picker := FilePicker{
		TreeView: tree,
	}

	picker.add(root, rootDir)
	tree.SetSelectedFunc(picker.onSelected)
	return &picker
}

// SetDoneFunc set handler to be called after user selects a file
func (p *FilePicker) SetDoneFunc(handler func(selected *os.File) error) *FilePicker {
	p.doneHandler = handler
	return p
}

// SetCancelFunc set handler to be called after user cancels file picker
func (p *FilePicker) SetCancelFunc(handler func()) *FilePicker {
	p.cancelHandler = handler
	return p
}

func (p *FilePicker) add(target *tview.TreeNode, path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetReference(nodeRef{
				FullPath: filepath.Join(path, file.Name()),
				Info:     file,
			}).SetSelectable(true)
		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
		}
		target.AddChild(node)
	}

	return nil
}

func (p *FilePicker) onSelected(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}

	ref := reference.(nodeRef)
	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		p.add(node, ref.FullPath)
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}
