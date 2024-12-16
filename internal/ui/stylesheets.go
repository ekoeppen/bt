package ui

import "github.com/charmbracelet/lipgloss"

type Stylesheet struct {
	SelectedPath lipgloss.Style

	FinfoPermissions lipgloss.Style
	FinfoLastUpdated lipgloss.Style
	FinfoSize        lipgloss.Style
	FinfoSep         lipgloss.Style

	OperationBar      lipgloss.Style
	OperationBarInput lipgloss.Style

	ErrBar lipgloss.Style

	HelpMsg     lipgloss.Style
	HelpContent lipgloss.Style

	TreeRegularFileName lipgloss.Style
	TreeDirectoryName   lipgloss.Style
	TreeLinkName        lipgloss.Style
	TreeMarkedNode      lipgloss.Style
	TreeSelectionArrow  lipgloss.Style
	TreeSelectionLine   lipgloss.Style
	TreeIndent          lipgloss.Style

	ContentPreview lipgloss.Style
}

var DefaultStylesheet = Stylesheet{
	SelectedPath: lipgloss.NewStyle().Foreground(lipgloss.Color("#74AC6D")),

	FinfoPermissions: lipgloss.NewStyle().Foreground(lipgloss.Color("#ACA46D")),
	FinfoLastUpdated: lipgloss.NewStyle().Foreground(lipgloss.Color("#E6E6E6")),
	FinfoSize:        lipgloss.NewStyle().Foreground(lipgloss.Color("#E6E6E6")),
	FinfoSep:         lipgloss.NewStyle().Foreground(lipgloss.Color("#2b2b2b")),

	OperationBar:      lipgloss.NewStyle().Foreground(lipgloss.Color("#E6E6E6")),
	OperationBarInput: lipgloss.NewStyle().Background(lipgloss.Color("#3C3C3C")),

	ErrBar:  lipgloss.NewStyle().Foreground(lipgloss.Color("#AC6D74")),
	HelpMsg: lipgloss.NewStyle().Foreground(lipgloss.Color("#ACA46D")),
	HelpContent: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8c7ca6")).
		BorderForeground(lipgloss.Color("#8c7ca6")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderLeft(true),

	TreeRegularFileName: lipgloss.NewStyle().Foreground(lipgloss.Color("#E6E6E6")),
	TreeDirectoryName:   lipgloss.NewStyle().Foreground(lipgloss.Color("#6D74AC")),
	TreeLinkName:        lipgloss.NewStyle().Foreground(lipgloss.Color("#6DACA4")),
	TreeMarkedNode: lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.InnerHalfBlockBorder()).
		Background(lipgloss.Color("#363636")),
	TreeSelectionArrow: lipgloss.NewStyle().Foreground(lipgloss.Color("#ACA46D")),
	TreeSelectionLine:  lipgloss.NewStyle().Foreground(lipgloss.Color("#ACA46D")),
	TreeIndent:         lipgloss.NewStyle().Foreground(lipgloss.Color("#363636")),

	ContentPreview: lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#a8a8a8")).
		BorderForeground(lipgloss.Color("#363636")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true),
}

var SolarizedTerm = Stylesheet{
	SelectedPath: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),

	FinfoPermissions: lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
	FinfoLastUpdated: lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
	FinfoSize:        lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
	FinfoSep:         lipgloss.NewStyle().Foreground(lipgloss.Color("0")),

	OperationBar:      lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
	OperationBarInput: lipgloss.NewStyle().Background(lipgloss.Color("0")),

	ErrBar:  lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
	HelpMsg: lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
	HelpContent: lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		BorderForeground(lipgloss.Color("0")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderLeft(true),

	TreeRegularFileName: lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
	TreeDirectoryName:   lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
	TreeLinkName:        lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
	TreeMarkedNode:      lipgloss.NewStyle().Background(lipgloss.Color("7")),
	TreeSelectionArrow:  lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
	TreeSelectionLine:   lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Background(lipgloss.Color("7")),
	TreeIndent:          lipgloss.NewStyle().Foreground(lipgloss.Color("0")),

	ContentPreview: lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("0")).
		BorderForeground(lipgloss.Color("0")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true),
}
