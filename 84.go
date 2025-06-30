package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the TUI
var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Padding(0, 1)
	modalStyle  = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Background(lipgloss.Color("235")).
			Padding(1).
			Margin(1).
			Width(60).
			Height(20).
			Align(lipgloss.Center, lipgloss.Center)
	menuStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1).
			MarginBottom(1).
			Width(80)
)

// Model represents the editor state
type model struct {
	textarea       textarea.Model
	textinput      textinput.Model
	filePath       string
	err            error
	previewMode    bool
	previewContent string
	helpMode       bool
	helpContent    string
	searchMode     bool
	menuVisible    bool
}

// InitialModel initializes the editor
func InitialModel(filePath string) model {
	// Ensure the file has a .md extension for new files
	if !strings.HasSuffix(filePath, ".md") && !strings.HasSuffix(filePath, ".bookmark.md") {
		filePath = filePath + ".md"
	}

	ta := textarea.New()
	ta.Placeholder = "Start typing your Markdown..."
	ta.Focus()
	ta.CharLimit = 0 // No limit
	ta.SetWidth(80)
	ta.SetHeight(20)

	// Read existing file content
	content, err := os.ReadFile(filePath)
	if err == nil {
		ta.SetValue(string(content))
	}

	ti := textinput.New()
	ti.Placeholder = "Enter search term..."
	ti.CharLimit = 80
	ti.Width = 50

	helpContent := `84 Keybindings:
F1: Show this help
F2: Save and exit
F3: Search text
F4: Toggle function key menu
F5: Toggle Markdown preview
F10/Esc: Quit without saving
Ctrl+H: Insert header
Ctrl+L: Insert list item
Ctrl+B: Insert bold markers
Ctrl+I: Insert italic markers
Ctrl+K: Insert link template
Ctrl+M: Insert inline code`

	return model{
		textarea:    ta,
		textinput:   ti,
		filePath:    filePath,
		previewMode: false,
		helpMode:    false,
		searchMode:  false,
		menuVisible: true,
		helpContent: helpContent,
	}
}

// Init returns the initial command
func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// getCursorPosition calculates the character index for the current line
func (m model) getCursorPosition() int {
	content := m.textarea.Value()
	lines := strings.Split(content, "\n")
	currentLine := m.textarea.Line()
	if currentLine >= len(lines) {
		currentLine = len(lines) - 1
	}
	if currentLine < 0 {
		currentLine = 0
	}
	cursorPos := 0
	for i := 0; i < currentLine; i++ {
		cursorPos += len(lines[i]) + 1 // +1 for newline
	}
	return cursorPos
}

// setCursorPosition moves the cursor to the specified character index
func (m model) setCursorPosition(pos int) {
	content := m.textarea.Value()
	lines := strings.Split(content, "\n")
	cursorPos := 0
	lineIndex := 0
	for i, line := range lines {
		if cursorPos+len(line)+1 > pos {
			lineIndex = i
			break
		}
		cursorPos += len(line) + 1
	}
	col := pos - cursorPos
	if col < 0 {
		col = 0
	}
	if col > len(lines[lineIndex]) {
		col = len(lines[lineIndex])
	}
	m.textarea.SetValue(content) // Ensure content is updated
	// Move cursor by simulating key presses (workaround for lack of SetCursor)
	for i := 0; i < lineIndex; i++ {
		m.textarea, _ = m.textarea.Update(tea.KeyMsg{Type: tea.KeyDown})
	}
	for i := 0; i < col; i++ {
		m.textarea, _ = m.textarea.Update(tea.KeyMsg{Type: tea.KeyRight})
	}
}

// Update handles user input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.previewMode || m.helpMode {
			switch msg.String() {
			case "esc", "f5":
				m.previewMode = false
				return m, nil
			case "f1":
				m.helpMode = false
				return m, nil
			}
			return m, nil
		}

		if m.searchMode {
			switch msg.String() {
			case "esc":
				m.searchMode = false
				m.textinput.Reset()
				m.textarea.Focus()
				return m, textarea.Blink
			case "enter":
				query := m.textinput.Value()
				if query != "" {
					content := m.textarea.Value()
					if idx := strings.Index(strings.ToLower(content), strings.ToLower(query)); idx >= 0 {
						m.setCursorPosition(idx)
					}
				}
				m.searchMode = false
				m.textinput.Reset()
				m.textarea.Focus()
				return m, textarea.Blink
			}
			var cmd tea.Cmd
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "esc", "f10":
			return m, tea.Quit
		case "f1":
			m.helpMode = true
			return m, nil
		case "f2":
			// Ensure directory exists
			if err := os.MkdirAll(filepath.Dir(m.filePath), 0755); err != nil {
				m.err = err
				return m, nil
			}
			err := os.WriteFile(m.filePath, []byte(m.textarea.Value()), 0644)
			if err != nil {
				m.err = err
			}
			return m, tea.Quit
		case "f3":
			m.searchMode = true
			m.textarea.Blur()
			m.textinput.Focus()
			return m, textinput.Blink
		case "f4":
			m.menuVisible = !m.menuVisible
			if m.menuVisible {
				m.textarea.SetHeight(20)
			} else {
				m.textarea.SetHeight(22) // Extra space when menu is hidden
			}
			return m, nil
		case "f5":
			rendered, err := glamour.Render(m.textarea.Value(), "dark")
			if err != nil {
				rendered = m.textarea.Value()
			}
			m.previewContent = rendered
			m.previewMode = true
			return m, nil
		case "ctrl+h":
			// Insert header, increment level if already a header
			content := m.textarea.Value()
			currentLine := m.textarea.Line()
			lines := strings.Split(content, "\n")
			if currentLine < len(lines) {
				lineContent := lines[currentLine]
				if strings.HasPrefix(lineContent, "#") {
					// Increment header level (up to 6)
					count := strings.Count(lineContent, "#")
					if count < 6 {
						m.textarea.InsertString("#")
					}
				} else {
					m.textarea.InsertString("# ")
				}
			} else {
				m.textarea.InsertString("# ")
			}
			return m, nil
		case "ctrl+l":
			// Insert unordered list item
			m.textarea.InsertString("- ")
			return m, nil
		case "ctrl+b":
			// Insert bold markers
			m.textarea.InsertString("****")
			m.setCursorPosition(m.getCursorPosition() - 2)
			return m, nil
		case "ctrl+i":
			// Insert italic markers
			m.textarea.InsertString("**")
			m.setCursorPosition(m.getCursorPosition() - 1)
			return m, nil
		case "ctrl+k":
			// Insert link template
			m.textarea.InsertString("[text](url)")
			m.setCursorPosition(m.getCursorPosition() - 10)
			return m, nil
		case "ctrl+m":
			// Insert inline code
			m.textarea.InsertString("``")
			m.setCursorPosition(m.getCursorPosition() - 1)
			return m, nil
		}

	case tea.MouseMsg:
		if m.previewMode || m.helpMode || m.searchMode {
			return m, nil
		}
		if msg.Type == tea.MouseLeft {
			// Left-click to set cursor position
			row := msg.Y - 3 // Adjust for menu, header, and padding
			if m.menuVisible {
				row--
			}
			if row >= 0 && row < m.textarea.Height() {
				col := msg.X
				lines := strings.Split(m.textarea.Value(), "\n")
				lineIndex := row
				if lineIndex >= len(lines) {
					lineIndex = len(lines) - 1
				}
				if lineIndex < 0 {
					lineIndex = 0
				}
				line := lines[lineIndex]
				charIndex := 0
				for i, r := range line {
					if i >= col {
						break
					}
					charIndex += len(string(r))
				}
				cursorPos := 0
				for i := 0; i < lineIndex; i++ {
					cursorPos += len(lines[i]) + 1 // +1 for newline
				}
				cursorPos += charIndex
				if cursorPos > len(m.textarea.Value()) {
					cursorPos = len(m.textarea.Value())
				}
				m.setCursorPosition(cursorPos)
			}
		}

	case error:
		m.err = msg
		return m, nil
	}

	if !m.previewMode && !m.searchMode && !m.helpMode {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s\n\n%s", headerStyle.Render("Error"), m.err.Error())
	}

	if m.helpMode {
		return lipgloss.Place(
			80, 24,
			lipgloss.Center, lipgloss.Center,
			modalStyle.Render(m.helpContent),
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("235")),
		)
	}

	if m.previewMode {
		return lipgloss.Place(
			80, 24,
			lipgloss.Center, lipgloss.Center,
			modalStyle.Render(m.previewContent),
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("235")),
		)
	}

	if m.searchMode {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			headerStyle.Render("Search in "+m.filePath),
			m.textinput.View(),
			"Press Enter to search, Esc to cancel",
		)
	}

	menu := ""
	if m.menuVisible {
		menu = menuStyle.Render("F1:Help  F2:Save  F3:Search  F4:Hide Menu  F5:Preview  F10:Quit")
	}

	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		menu,
		headerStyle.Render("Editing "+m.filePath),
		m.textarea.View(),
	)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: 84 <file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	p := tea.NewProgram(InitialModel(filePath), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
