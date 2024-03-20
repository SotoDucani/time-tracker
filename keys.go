package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Select      key.Binding
	Quit        key.Binding
	Redraw      key.Binding
	ResetDay    key.Binding
	ResetBucket key.Binding
	MinuteUp    key.Binding
	MinuteDown  key.Binding
	HourUp      key.Binding
	HourDown    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit, k.Redraw, k.ResetDay, k.ResetBucket}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},                         //First Column
		{k.MinuteUp, k.MinuteDown, k.HourUp, k.HourDown}, // Second Column
		{k.Quit, k.Redraw, k.ResetDay, k.ResetBucket},    // Third Column
	}
}

func setKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "Move Up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "Move Down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("↵/␠", "Select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "Quit"),
		),
		Redraw: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "Redraw"),
		),
		ResetDay: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "Reset Day"),
		),
		ResetBucket: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("ctrl+b", "Reset Bucket"),
			key.WithDisabled(), // TODO
		),
		MinuteUp: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "+15 Minutes"),
		),
		MinuteDown: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "-15 Minutes"),
		),
		HourUp: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "+1 Hour"),
		),
		HourDown: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "-1 Hour"),
		),
	}
}
