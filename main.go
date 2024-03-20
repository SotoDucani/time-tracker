package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type elapsedTickMsg struct {
	t time.Time
}

func elapsedtimeTick(timeUpdateInterval time.Duration) tea.Cmd {
	return tea.Tick(timeUpdateInterval, func(t time.Time) tea.Msg {
		return elapsedTickMsg{t: t}
	})
}

func (m model) Init() tea.Cmd {
	log.Printf("%v", m.buckets)
	return tea.Batch(m.spin.Tick, tea.ClearScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		log.Printf("Window Size Msg")
		m.help.Width = msg.Width
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch {
		default:
			log.Printf("Default Handler - KeyMsg recieved: %s", msg.String())
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			log.Printf("KeyMsg recieved: %s", msg.String())
			return m, tea.Quit
		case key.Matches(msg, m.keys.Redraw):
			log.Printf("KeyMsg recieved: %s", msg.String())
			return m, tea.ClearScreen
		case key.Matches(msg, m.keys.Up):
			log.Printf("KeyMsg recieved: %s", msg.String())
			if m.cursor > 1 {
				m.cursor--
			}
			return m, nil
		case key.Matches(msg, m.keys.Down):
			log.Printf("KeyMsg recieved: %s", msg.String())
			if m.cursor < len(m.buckets)-1 {
				m.cursor++
			}
			return m, nil
		case key.Matches(msg, m.keys.Select):
			log.Printf("KeyMsg recieved: %s", msg.String())
			return selectBucket(&m)
		case key.Matches(msg, m.keys.ResetDay):
			log.Printf("KeyMsg recieved: %s", msg.String())
			return resetDay(&m)
		case key.Matches(msg, m.keys.ResetBucket):
			log.Printf("KeyMsg recieved: %s", msg.String())
			return resetBucket(&m)
		}
	case elapsedTickMsg:
		// Doing this at slower tick pace so we're not tying updates to the spinner FPS
		// Pace is dictated by m.timeUpdateInterval
		if m.activeSelection {
			startTime := time.Now()
			// Add time to the selected bucket
			addElapsedTime(startTime, &m)
			return m, elapsedtimeTick(m.timeUpdateInterval)
		} else {
			return m, nil
		}
	default:
		var cmd tea.Cmd
		//log.Printf("Default msg handler")
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	s := fmt.Sprintf("Total Time Tracked: %v\n\n", m.buckets[0].elapsedTime.Round(time.Second).String())
	s += "Select active time bucket:\n"

	// Determines the output string for each time bucket in turn
	for i, choice := range m.buckets {
		if choice.id == 0 {
			continue
		}

		// 3 Characters
		cursor := "   "
		if m.cursor == i {
			cursor = ">> "
		}

		// 3 Characters - Points type
		itSpins := "   "
		if m.selected == choice.id {
			itSpins = fmt.Sprintf("%s", m.spin.View())
		}

		// 0 to 4 characters
		lvl := ""
		if choice.level == "second" {
			lvl = " |->"
		}

		// 3 Characters
		checked := "[ ]"
		if m.selected == choice.id {
			checked = fmt.Sprintf("[x]")
		}

		// Display line cursor, itSpins, lvl, checked, name ,elapsedTime
		s += fmt.Sprintf("%s%s%s%s %s == %v\n", cursor, itSpins, lvl, checked, choice.name, choice.elapsedTime.Round(time.Second).String())
	}

	// Writes the output string for the keybinds
	m.help.ShowAll = true
	helpView := m.help.View(m.keys)
	s += fmt.Sprintf("\n%s\n", helpView)

	return s
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(initializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error occured: %v", err)
		os.Exit(1)
	}
}
