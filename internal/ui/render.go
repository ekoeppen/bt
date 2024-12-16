package ui

import (
	"fmt"
	"math"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"

	"github.com/ekoeppen/bt/internal/state"
	t "github.com/ekoeppen/bt/internal/tree"
	"github.com/ekoeppen/bt/pkg/stack"
)

const (
	previewBytesLimit int64 = 10_000

	minHeight = 10
	minWidth  = 10

	arrow               = " <-"
	indentParent        = "│  "
	indentCurrent       = "├─ "
	indentCurrentLast   = "└─ "
	indentEmpty         = "   "
	emptydirContentName = "..."

	tooSmall                 = "too small =("
	binaryContentPlaceholder = "<binary content>"
)

type Renderer struct {
	Style       Stylesheet
	EdgePadding int
	Preview     bool
	offsetMem   int
	previewBuff [previewBytesLimit]byte
}

func (r *Renderer) Render(s *state.State, winHeight, winWidth int) string {
	if winWidth < minWidth || winHeight < minHeight {
		return tooSmall
	}

	renderedHeading, headLen := r.renderHeading(s, winWidth)

	// section is half a screen, devided vertically
	// left for tree, right for file preview
	var sectionWidth int
	if r.Preview {
		sectionWidth = int(math.Floor(0.5 * float64(winWidth)))
	} else {
		sectionWidth = winWidth
	}

	renderedTree := r.renderTree(s.Tree, winHeight-headLen, sectionWidth)

	var rightPane string

	if s.HelpToggle {
		renderedContent := r.renderSelectedFileContent(s.Tree, winHeight-headLen, sectionWidth)
		rightPane = lipgloss.JoinVertical(lipgloss.Left, renderedContent)
	} else {
		renderedContent := r.renderSelectedFileContent(s.Tree, winHeight-headLen, sectionWidth)
		rightPane = renderedContent
	}

	renderedTreeWithContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderedTree,
		rightPane,
	)

	return renderedHeading + "\n" + renderedTreeWithContent
}

func (r *Renderer) renderHeading(s *state.State, width int) (string, int) {
	markedPath := ""
	if s.Tree.Marked != nil {
		markedPath = s.Tree.Marked.Path
	}

	operationBar := fmt.Sprintf(": %s", s.OpBuf.Repr())
	if markedPath != "" {
		operationBar += fmt.Sprintf(" [%s]", markedPath)
	}

	if s.OpBuf.IsInput() {
		operationBar += fmt.Sprintf(" │ %s │", r.Style.OperationBarInput.Render(string(s.InputBuf)))
	}

	header := []string{
		r.Style.OperationBar.Render(operationBar),
		r.Style.ErrBar.Render(s.ErrBuf),
	}
	return strings.Join(header, "\n"), len(header)
}

func (r *Renderer) renderTree(tree *t.Tree, height, width int) string {
	renderedTreeLines, selectedRow := r.renderTreeFull(tree, width)
	croppedTreeLines := r.cropTree(renderedTreeLines, selectedRow, height)

	treeStyle := lipgloss.
		NewStyle().
		MaxWidth(width).
		MarginRight(width)

	return treeStyle.Render(strings.Join(croppedTreeLines, "\n"))
}

func (r *Renderer) renderSelectedFileContent(tree *t.Tree, height, width int) string {
	if !r.Preview {
		return ""
	}
	n, err := tree.ReadSelectedChildContent(r.previewBuff[:], previewBytesLimit)
	if err != nil {
		return ""
	}
	content := r.previewBuff[:n]

	contentStyle := r.Style.ContentPreview.MaxWidth(width - 1) // -1 for border...

	var contentLines []string
	if !utf8.Valid(content) {
		contentLines = []string{binaryContentPlaceholder}
	} else {
		contentLines = strings.Split(string(content), "\n")
		contentLines = contentLines[:max(min(height, len(contentLines)), 0)]
	}
	return contentStyle.Render(strings.Join(contentLines, "\n"))
}

// Crops tree lines, such that current line is visible and view is consistent.
func (r *Renderer) cropTree(lines []string, currentLine int, height int) []string {
	linesLen := len(lines)

	// determining offset and limit based on selected row
	offset := r.offsetMem
	limit := linesLen

	// cursor is out for 'top' boundary
	if currentLine+1 > height+offset-r.EdgePadding {
		offset = max(min(currentLine+1-height+r.EdgePadding, linesLen-height), 0)
	}
	// cursor is out for 'bottom' boundary
	if currentLine < r.EdgePadding+offset {
		offset = max(currentLine-r.EdgePadding, 0)
	}
	r.offsetMem = offset
	limit = min(height+offset, linesLen)
	return lines[offset:limit]
}

// Returns lines as slice and index of selected line.
func (r *Renderer) renderTreeFull(tree *t.Tree, width int) ([]string, int) {
	linen := -1
	currentLine := 0

	type stackEl struct {
		*t.Node
		string
		bool
	}
	lines := []string{}
	s := stack.NewStack(stackEl{tree.Root, "", false})

	for s.Len() > 0 {
		el := s.Pop()
		linen += 1

		node := el.Node
		isLast := el.bool
		parentIndent := el.string

		var indent string
		if node == tree.Root {
			indent = ""
		} else if isLast {
			indent = parentIndent + indentCurrentLast
			parentIndent = parentIndent + indentEmpty
		} else {
			indent = parentIndent + indentCurrent
			parentIndent = parentIndent + indentParent
		}

		if node == nil {
			continue
		}

		name := node.Info.Name()
		nameRuneCountNoStyle := utf8.RuneCountInString(name)
		indentRuneCount := utf8.RuneCountInString(indent)

		if nameRuneCountNoStyle+indentRuneCount > width {
			name = string([]rune(name)[:max(0, width-indentRuneCount)]) + "..."
		}

		indent = r.Style.TreeIndent.Render(indent)

		if tree.GetSelectedChild() == node {
			name = r.Style.TreeSelectionLine.Render(name)
		} else if node.Info.IsDir() {
			name = r.Style.TreeDirectoryName.Render(name)
		} else if node.Info.Mode()&os.ModeSymlink == os.ModeSymlink {
			name = r.Style.TreeLinkName.Render(name)
		} else {
			name = r.Style.TreeRegularFileName.Render(name)
		}

		if tree.Marked == node {
			name = r.Style.TreeMarkedNode.Render(name)
		}

		repr := indent + name

		lines = append(lines, repr)

		if node.Children != nil {
			// current directory is empty
			if len(node.Children) == 0 && tree.CurrentDir == node {
				emptyIndent := r.Style.TreeIndent.Render(parentIndent + indentCurrentLast)
				lines = append(lines, emptyIndent+emptydirContentName+r.Style.TreeSelectionArrow.Render(arrow))
				currentLine = linen + 1
			}
			for i := len(node.Children) - 1; i >= 0; i-- {
				ch := node.Children[i]
				s.Push(stackEl{ch, parentIndent, i == len(node.Children)-1})
			}
		}
	}
	return lines, currentLine
}

var sizes = [...]string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

func formatSize(s float64, base float64) string {
	unitsLimit := len(sizes)
	i := 0
	for s >= base && i < unitsLimit {
		s = s / base
		i++
	}
	f := "%.0f %s"
	if i > 1 {
		f = "%.2f %s"
	}
	return fmt.Sprintf(f, s, sizes[i])
}
