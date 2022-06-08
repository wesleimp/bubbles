package textarea

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	textarea := newTextArea()
	view := textarea.View()

	if !strings.Contains(view, ">") {
		t.Log(view)
		t.Error("Textarea did not render the prompt")
	}

	if !strings.Contains(view, "World!") {
		t.Log(view)
		t.Error("Textarea did not render the placeholder")
	}
}

func TestInput(t *testing.T) {
	textarea := newTextArea()

	input := "foo"

	for _, k := range []rune(input) {
		textarea, _ = textarea.Update(keyPress(k))
	}

	view := textarea.View()

	if !strings.Contains(view, input) {
		t.Log(view)
		t.Error("Textarea did not render the input")
	}

	if textarea.col != len(input) {
		t.Log(view)
		t.Error("Textarea did not move the cursor to the correct position")
	}
}

func TestWrap(t *testing.T) {
	textarea := newTextArea()
	textarea.Width = 5
	textarea.LineLimit = 5
	textarea.Height = 5
	textarea.CharLimit = 60

	textarea, _ = textarea.Update(initialBlinkMsg{})

	input := "foo bar baz"

	for _, k := range []rune(input) {
		textarea, _ = textarea.Update(keyPress(k))
	}

	view := textarea.View()

	for _, word := range strings.Split(input, " ") {
		if !strings.Contains(view, word) {
			t.Log(view)
			t.Error("Textarea did not render the input")
		}
	}

	// Due to the word wrapping, each word will be on a new line and the
	// textarea will look like this:
	//
	// > foo
	// > bar
	// > baz█
	if textarea.row != 2 && textarea.col != 3 {
		t.Log(view)
		t.Error("Textarea did not move the cursor to the correct position")
	}
}

func TestLineNumbers(t *testing.T) {
	textarea := newTextArea()
	textarea.ShowLineNumbers = true

	lines := 5

	textarea.LineLimit = lines
	textarea.Height = lines

	textarea, _ = textarea.Update(initialBlinkMsg{})

	view := textarea.View()

	for i := 0; i < lines; i++ {
		if !strings.Contains(view, fmt.Sprint(i+1)) {
			t.Log(view)
			t.Error("Textarea did not render the line numbers")
		}
	}
}

func TestCharLimit(t *testing.T) {
	textarea := newTextArea()

	// First input (foo bar) should be accepted as it will fall within the
	// CharLimit. Second input (baz) should not appear in the input.
	input := []string{"foo bar", "baz"}
	textarea.CharLimit = len(input[0])

	for _, k := range []rune(strings.Join(input, " ")) {
		textarea, _ = textarea.Update(keyPress(k))
	}

	view := textarea.View()
	if strings.Contains(view, input[1]) {
		t.Log(view)
		t.Error("Textarea should not include input past the character limit")
	}
}

func newTextArea() Model {
	textarea := New()

	textarea.Prompt = "> "
	textarea.Placeholder = "Hello, World!"

	textarea.Focus()

	textarea, _ = textarea.Update(initialBlinkMsg{})

	return textarea
}

func keyPress(key rune) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}, Alt: false}
}
