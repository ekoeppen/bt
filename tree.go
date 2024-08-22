package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type Tree struct {
	Root        *Node
	CurrentDir  *Node
	Marked      *Node
	sortingFunc NodeSortingFunc
}

func (t *Tree) GetSelectedChild() *Node {
	if len(t.CurrentDir.Children) > 0 {
		return t.CurrentDir.Children[t.CurrentDir.selectedChildIdx]
	}
	return nil
}
func (t *Tree) RenameMarked(name string) error {
	if t.Marked == nil {
		return nil
	}
	err := os.Rename(t.Marked.Path, filepath.Join(t.Marked.Parent.Path, name))
	if err != nil {
		return err
	}
	err = t.Marked.Parent.readChildren(t.sortingFunc)
	if err != nil {
		return err
	}
	t.Marked = nil
	return nil
}
func (t *Tree) CreateFileInCurrent(name string) error {
	_, err := os.Create(filepath.Join(t.CurrentDir.Path, name))
	if err != nil {
		return err
	}
	return t.CurrentDir.readChildren(t.sortingFunc)
}
func (t *Tree) CreateDirectoryInCurrent(name string) error {
	err := os.Mkdir(filepath.Join(t.CurrentDir.Path, name), os.ModePerm)
	if err != nil {
		return err
	}
	return t.CurrentDir.readChildren(t.sortingFunc)
}
func (t *Tree) ReadSelectedChildContent(buf []byte, limit int64) (int, error) {
	selectedNode := t.GetSelectedChild()
	if selectedNode == nil || !selectedNode.Info.Mode().IsRegular() {
		return 0, fmt.Errorf("file not selected or is irregular")
	}
	f, err := os.Open(selectedNode.Path)
	defer f.Close()
	if err != nil {
		return 0, err
	}
	limitedReader := io.LimitReader(f, limit)
	n, err := limitedReader.Read(buf)
	if err != nil {
		return 0, err
	}
	return n, nil
}
func (t *Tree) SelectNextChild() {
	if t.CurrentDir.selectedChildIdx < len(t.CurrentDir.Children)-1 {
		t.CurrentDir.selectedChildIdx += 1
	}
}
func (t *Tree) SelectPreviousChild() {
	if t.CurrentDir.selectedChildIdx > 0 {
		t.CurrentDir.selectedChildIdx -= 1
	}
}
func (t *Tree) SetSelectedChildAsCurrent() error {
	selectedChild := t.GetSelectedChild()
	if selectedChild == nil {
		return nil
	}
	if !selectedChild.Info.IsDir() {
		return nil
	}
	if selectedChild.Children == nil {
		err := selectedChild.readChildren(t.sortingFunc)
		if err != nil {
			return err
		}
	}
	t.CurrentDir = selectedChild
	return nil
}
func (t *Tree) SetParentAsCurrent() {
	if t.CurrentDir.Parent != nil {
		currentName := t.CurrentDir.Info.Name()
		// setting parent current to match the directory, that we're leaving
		newParentIdx := slices.IndexFunc(t.CurrentDir.Parent.Children, func(n *Node) bool { return n.Info.Name() == currentName })
		t.CurrentDir.Parent.selectedChildIdx = newParentIdx

		t.CurrentDir = t.CurrentDir.Parent
	}
}
func (t *Tree) MarkSelectedChild() bool {
	if selected := t.GetSelectedChild(); selected != nil {
		t.Marked = selected
		return true
	}
	return false
}
func (t *Tree) DropMark() {
	t.Marked = nil
}
func (t *Tree) DeleteMarked() error {
	if t.Marked == nil {
		return nil
	}
	cmd := exec.Command("rm", "-r", t.Marked.Path)
	err := cmd.Run()
	if err != nil {
		return err // todo: this is not the same error...?
	}
	err = t.Marked.Parent.readChildren(t.sortingFunc)
	if err != nil {
		return err
	}
	t.Marked = nil
	return nil
}
func (t *Tree) CopyMarkedToCurrentDir() error {
	if t.Marked == nil {
		return nil
	}
	targetDir := t.CurrentDir.Path
	targetFileName, err := generateNewFileName(t.Marked.Info.Name(), targetDir)
	if err != nil {
		return err
	}
	targetPath := filepath.Join(targetDir, targetFileName)

	cmd := exec.Command("cp", "-r", t.Marked.Path, targetPath)
	err = cmd.Run()
	if err != nil {
		return err // todo: this is not the same error...?
	}
	err = t.CurrentDir.readChildren(t.sortingFunc)
	if err != nil {
		return err
	}
	t.Marked = nil
	return nil
}
func (t *Tree) MoveMarkedToCurrentDir() error {
	if t.Marked == nil {
		return nil
	}
	targetDir := t.CurrentDir.Path
	targetFileName, err := generateNewFileName(t.Marked.Info.Name(), targetDir)
	if err != nil {
		return err
	}
	targetPath := filepath.Join(targetDir, targetFileName)

	cmd := exec.Command("mv", t.Marked.Path, targetPath)
	err = cmd.Run()
	if err != nil {
		return err // todo: this is not the same error...?
	}
	err = t.CurrentDir.readChildren(t.sortingFunc)
	if err != nil {
		return err
	}
	err = t.Marked.Parent.readChildren(t.sortingFunc)
	if err != nil {
		return err
	}
	t.Marked = nil
	return nil
}
func (t *Tree) CollapseOrExpandSelected() error {
	selectedChild := t.GetSelectedChild()
	if selectedChild == nil {
		return nil
	}
	if selectedChild.Children != nil {
		selectedChild.orphanChildren()
	} else {
		err := selectedChild.readChildren(t.sortingFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitTree(dir string, sortingFunc NodeSortingFunc) (Tree, error) {
	var tree Tree
	rootInfo, err := os.Lstat(dir)
	if err != nil {
		return tree, err
	}
	if !rootInfo.IsDir() {
		return tree, fmt.Errorf("%s is not a directory", dir)
	}
	if sortingFunc == nil {
		sortingFunc = defaultNodeSorting
	}

	root := &Node{
		Path:     dir,
		Info:     rootInfo,
		Parent:   nil,
		Children: []*Node{},
	}

	err = root.readChildren(sortingFunc)
	if err != nil {
		return tree, err
	}
	if len(root.Children) == 0 {
		return tree, fmt.Errorf("Can't initialize on empty directory '%s'", dir)
	}

	tree = Tree{Root: root, CurrentDir: root, sortingFunc: sortingFunc}
	return tree, nil
}

type NodeSortingFunc func(a, b *Node) int

type Node struct {
	Path             string
	Info             fs.FileInfo
	Children         []*Node // nil - not read or it's a file
	Parent           *Node
	selectedChildIdx int
}

func (n *Node) SelectLast() {
	// can I just check for len here?
	if n.Children != nil && len(n.Children) > 0 {
		n.selectedChildIdx = len(n.Children) - 1
	}
}
func (n *Node) SelectFirst() {
	n.selectedChildIdx = 0
}
func (n *Node) readChildren(sortFunc NodeSortingFunc) error {
	if !n.Info.IsDir() {
		return nil
	}
	children, err := os.ReadDir(n.Path)
	if err != nil {
		return err
	}
	chNodes := []*Node{}

	for _, ch := range children {
		ch := ch
		chInfo, err := ch.Info()
		if err != nil {
			return err
		}
		// Looking if child already exist
		var childToAdd *Node
		if n.Children != nil {
			for _, ech := range n.Children {
				if ech.Info.Name() == chInfo.Name() {
					childToAdd = ech
					break
				}
			}
		}
		if childToAdd == nil {
			childToAdd = &Node{
				Path:     filepath.Join(n.Path, chInfo.Name()),
				Info:     chInfo,
				Children: nil,
				Parent:   n,
			}
		}
		chNodes = append(chNodes, childToAdd)
	}
	slices.SortFunc(chNodes, sortFunc)
	n.Children = chNodes

	// updateing selected child index if it's out of bounds after update
	n.selectedChildIdx = max(min(n.selectedChildIdx, len(n.Children)-1), 0)
	return nil
}
func (n *Node) orphanChildren() {
	n.Children = nil
}

// Checks if fname already exists in targetDir.
// Adds "copy_" prefix (multiple times), until new file name becomes unique in derecotry.
func generateNewFileName(fname, targetDir string) (string, error) {
	currentDirContent, err := os.ReadDir(targetDir)
	if err != nil {
		return "", err
	}
	for slices.ContainsFunc(currentDirContent, func(e fs.DirEntry) bool { return e.Name() == fname }) {
		fname = "copy_" + fname
	}
	return fname, nil
}

func defaultNodeSorting(a, b *Node) int {
	// dirs first
	if a.Info.IsDir() != b.Info.IsDir() {
		if a.Info.IsDir() {
			return -1
		} else {
			return 1
		}
	}
	return strings.Compare(strings.ToLower(a.Info.Name()), strings.ToLower(b.Info.Name()))
}
