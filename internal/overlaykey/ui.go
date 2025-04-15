package overlaykey

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/colors"

	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	FocusWidth    key.Binding
	FocusInterval key.Binding
	FocusShift    key.Binding
	FocusStart    key.Binding
	RemoveStart   key.Binding
	Increase      key.Binding
	Decrease      key.Binding
	Escape        key.Binding
}

var keys = keymap{
	FocusWidth:    Key("Focus Width", ":"),
	FocusInterval: Key("Focus Interval", "/"),
	FocusShift:    Key("Focus Shift", "^"),
	FocusStart:    Key("Focus Start", "S"),
	RemoveStart:   Key("Remove Start", "s"),
	Increase:      Key("Increase", "+"),
	Decrease:      Key("Decrease", "-"),
	Escape:        Key("Escape", "esc", "enter"),
}

func Key(help string, keyboardKey ...string) key.Binding {
	return key.NewBinding(key.WithKeys(keyboardKey...), key.WithHelp(keyboardKey[0], help))
}

type Model struct {
	overlayKey        OverlayPeriodicity
	focus             focus
	firstDigitApplied bool
}

func InitModel() Model {
	return Model{ROOT, FOCUS_NOTHING, false}
}

func (m *Model) SetOverlayKey(op OverlayPeriodicity) {
	m.overlayKey = op
}

func (m Model) GetKey() OverlayPeriodicity {
	return m.overlayKey
}

func (m *Model) Focus(shouldFocus bool) {
	if shouldFocus {
		m.focus = FOCUS_SHIFT
	} else {
		m.focus = FOCUS_NOTHING
	}
}

type focus int

const (
	FOCUS_NOTHING focus = iota
	FOCUS_SHIFT
	FOCUS_WIDTH
	FOCUS_INTERVAL
	FOCUS_START
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() >= "0" && msg.String() <= "9":
			numberString := msg.String()
			newDigit, _ := strconv.Atoi(numberString)
			m.ApplyDigit(newDigit)
			m.firstDigitApplied = true
			if m.focus == FOCUS_START && m.overlayKey.StartCycle == 0 {
				m.focus = FOCUS_SHIFT
			}
			if m.focus == FOCUS_WIDTH && m.overlayKey.Width == 0 {
				m.focus = FOCUS_SHIFT
			}
		case key.Matches(msg, keys.Escape):
			m.focus = FOCUS_NOTHING
			m.firstDigitApplied = false
			return m, Updated(m.overlayKey, false)
		case key.Matches(msg, keys.FocusWidth):
			m.focus = FOCUS_WIDTH
			if m.overlayKey.Width == 0 {
				m.overlayKey.Width = 1
			}
			m.firstDigitApplied = false
		case key.Matches(msg, keys.FocusInterval):
			m.focus = FOCUS_INTERVAL
			m.firstDigitApplied = false
		case key.Matches(msg, keys.FocusShift):
			m.focus = FOCUS_SHIFT
			m.firstDigitApplied = false
		case key.Matches(msg, keys.FocusStart):
			m.focus = FOCUS_START
			if m.overlayKey.StartCycle == 0 {
				m.overlayKey.StartCycle = 1
			}
			m.firstDigitApplied = false
		case key.Matches(msg, keys.RemoveStart):
			m.focus = FOCUS_SHIFT
			m.overlayKey.StartCycle = 0
			m.firstDigitApplied = false
		case key.Matches(msg, keys.Increase):
			switch m.focus {
			case FOCUS_SHIFT:
				m.overlayKey.IncrementShift()
			case FOCUS_INTERVAL:
				m.overlayKey.IncrementInterval()
			case FOCUS_WIDTH:
				m.overlayKey.IncrementWidth()
			case FOCUS_START:
				m.overlayKey.IncrementStartCycle()
			}
		case key.Matches(msg, keys.Decrease):
			switch m.focus {
			case FOCUS_SHIFT:
				m.overlayKey.DecrementShift()
			case FOCUS_INTERVAL:
				m.overlayKey.DecrementInterval()
			case FOCUS_WIDTH:
				m.overlayKey.DecrementWidth()
			case FOCUS_START:
				m.overlayKey.DecrementStartCycle()
			}
		}
	}
	return m, Updated(m.overlayKey, true)
}

func (m *Model) ApplyDigit(newDigit int) {
	switch m.focus {
	case FOCUS_SHIFT:
		m.overlayKey.Shift = m.UnshiftDigit(m.overlayKey.Shift, newDigit)
		if m.overlayKey.Shift == 0 {
			m.overlayKey.Shift = 1
		}
	case FOCUS_INTERVAL:
		m.overlayKey.Interval = m.UnshiftDigit(m.overlayKey.Interval, newDigit)
		if m.overlayKey.Interval == 0 {
			m.overlayKey.Interval = 1
		}
	case FOCUS_WIDTH:
		m.overlayKey.Width = m.UnshiftDigit(m.overlayKey.Width, newDigit)
	case FOCUS_START:
		m.overlayKey.StartCycle = m.UnshiftDigit(m.overlayKey.StartCycle, newDigit)
	}
}

func (m Model) UnshiftDigit(digits uint8, newDigit int) uint8 {
	if m.firstDigitApplied {
		return uint8((int(digits)%10)*10 + newDigit)
	} else {
		return uint8(newDigit)
	}
}

type UpdatedOverlayKey struct {
	OverlayKey OverlayPeriodicity
	HasFocus   bool
}

func Updated(overlayKey OverlayPeriodicity, maintainsFocus bool) tea.Cmd {
	return func() tea.Msg {
		return UpdatedOverlayKey{
			OverlayKey: overlayKey,
			HasFocus:   maintainsFocus,
		}
	}
}

func View(ok OverlayPeriodicity) string {
	var shift, interval, width, start string
	var buf strings.Builder

	shift = NormalColor(ok.Shift)
	interval = NormalColor(ok.Interval)
	width = NormalColor(ok.Width)
	start = NormalColor(ok.StartCycle)

	buf.WriteString(shift)
	if ok.Width > 0 {
		buf.WriteString(":")
		buf.WriteString(width)
	}
	buf.WriteString("/")
	buf.WriteString(interval)
	if ok.StartCycle > 0 {
		buf.WriteString("S")
		buf.WriteString(start)
	}
	return buf.String()
}

func (m Model) ViewOverlay() string {
	var shift, interval, width, start string
	var buf strings.Builder

	shift = NumberColor(m.overlayKey.Shift)
	interval = NumberColor(m.overlayKey.Interval)
	width = NumberColor(m.overlayKey.Width)
	start = NumberColor(m.overlayKey.StartCycle)

	switch m.focus {
	case FOCUS_SHIFT:
		shift = SelectedColor(m.overlayKey.Shift)
	case FOCUS_WIDTH:
		width = SelectedColor(m.overlayKey.Width)
	case FOCUS_INTERVAL:
		interval = SelectedColor(m.overlayKey.Interval)
	case FOCUS_START:
		start = SelectedColor(m.overlayKey.StartCycle)
	}

	buf.WriteString(shift)
	if m.overlayKey.Width > 0 {
		buf.WriteString(":")
		buf.WriteString(width)
	}
	buf.WriteString("/")
	buf.WriteString(interval)
	if m.overlayKey.StartCycle > 0 {
		buf.WriteString("S")
		buf.WriteString(start)
	}
	return buf.String()
}

func NumberColor(number uint8) string {
	return colors.NumberStyle.Render(strconv.Itoa(int(number)))
}

func SelectedColor(number uint8) string {
	return colors.SelectedStyle.Render(strconv.Itoa(int(number)))
}

func NormalColor(number uint8) string {
	return strconv.Itoa(int(number))
}
