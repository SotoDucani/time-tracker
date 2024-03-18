package main

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/peterbourgon/diskv/v3"
)

type model struct {
	keys               keyMap
	help               help.Model
	buckets            []timeBucket
	cursor             int
	selected           int
	activeSelection    bool
	spin               spinner.Model
	timeUpdateInterval time.Duration
	datastore          *diskv.Diskv
}

func setSpinner() spinner.Model {
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("31"))
	return spin
}

func (m model) reloadBuckets() error {
	for i := range m.buckets {
		// Ignoring error, most likely it's not needed
		m.buckets[i].elapsedTime, _ = readBucketData(m.buckets[i], m.datastore)
	}
	return nil
}

func initializeModel() model {
	m := model{
		buckets:            setupBuckets(),  // setup custom buckets
		cursor:             1,               // start at the top of the list
		selected:           -1,              // put selected out of the array so nothing is chedked
		activeSelection:    false,           // make sure we don't index-out-of-range
		spin:               setSpinner(),    // setup the spinner object
		keys:               setKeyMap(),     // setup the keymap
		help:               help.New(),      // init help
		timeUpdateInterval: 1 * time.Second, // Set interval for updating elapsedTime
		datastore:          diskv_init(),
	}
	m.reloadBuckets()
	return m
}
