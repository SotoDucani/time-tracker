package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type timeBucket struct {
	id           int
	name         string
	level        string
	parentBucket int
	startTime    time.Time
	elapsedTime  time.Duration
}

func addIncrementTime(increment time.Duration, m *model) {
	// We operate on the CURSOR targeted bucket here
	// Update the cursor bucket
	if m.buckets[m.cursor].elapsedTime >= increment.Abs() || increment > 0*time.Second {
		m.buckets[m.cursor].elapsedTime += increment
	} else {
		m.buckets[m.cursor].elapsedTime = 0 * time.Second
	}
	storeBucketData(m.buckets[m.cursor], m.datastore)

	// Update the parent bucket
	if m.buckets[m.buckets[m.cursor].parentBucket].elapsedTime >= increment.Abs() || increment > 0*time.Second {
		m.buckets[m.buckets[m.cursor].parentBucket].elapsedTime += increment
	} else {
		m.buckets[m.buckets[m.cursor].parentBucket].elapsedTime = 0 * time.Second
	}
	storeBucketData(m.buckets[m.buckets[m.cursor].parentBucket], m.datastore)

	// Update the total, if we happen to be operating on a second level bucket
	if m.buckets[m.cursor].level == "second" {
		if m.buckets[0].elapsedTime >= increment.Abs() || increment > 0*time.Second {
			m.buckets[0].elapsedTime += increment
		} else {
			m.buckets[0].elapsedTime += 0 * time.Second
		}
		storeBucketData(m.buckets[0], m.datastore)
	}
}

func addElapsedTime(startTime time.Time, m *model) {
	// Update the selected bucket
	m.buckets[m.selected].elapsedTime += time.Since(m.buckets[m.selected].startTime)
	storeBucketData(m.buckets[m.selected], m.datastore)
	m.buckets[m.selected].startTime = startTime

	// Update the parent bucket
	m.buckets[m.buckets[m.selected].parentBucket].elapsedTime += time.Since(m.buckets[m.buckets[m.selected].parentBucket].startTime)
	storeBucketData(m.buckets[m.buckets[m.selected].parentBucket], m.datastore)
	m.buckets[m.buckets[m.selected].parentBucket].startTime = startTime

	// Update the total, if we happen to be operating on a second level bucket
	if m.buckets[m.selected].level == "second" {
		m.buckets[0].elapsedTime += time.Since(m.buckets[0].startTime)
		storeBucketData(m.buckets[0], m.datastore)
		m.buckets[0].startTime = startTime
	}
}

func selectBucket(m *model) (tea.Model, tea.Cmd) {
	startTime := time.Now()
	if m.activeSelection {
		// If we were tracking time already, update that bucket and it's parents
		addElapsedTime(startTime, m)
	}
	if m.cursor == m.selected {
		// Clear the selection marker and indicate no active selection
		m.selected = -1
		m.activeSelection = false
		return m, nil
	} else {
		// Update selected to where the cursor is and indicate active selection
		m.selected = m.cursor
		m.activeSelection = true
		m.buckets[m.selected].startTime = startTime
		m.buckets[m.buckets[m.selected].parentBucket].startTime = startTime
		if m.buckets[m.selected].level == "second" {
			m.buckets[0].startTime = startTime
		}
		// Start the tick loop for the selection
		return m, elapsedtimeTick(m.timeUpdateInterval)
	}
}

func resetDay(m *model) (tea.Model, tea.Cmd) {
	// Halt active ticking, we don't care about the data
	if m.activeSelection {
		m.selected = -1
		m.activeSelection = false
	}
	// Loop through all buckets and set time to 0
	for i := range m.buckets {
		m.buckets[i].elapsedTime = 0 * time.Second
		storeBucketData(m.buckets[i], m.datastore)
	}
	return m, nil
}

// TODO: Calling a reset action on a parent still needs to keep the sum of it's child buckets
// parent 6 hours; child1 2 hours, child2 2 hours
// resetBucket on the parent needs to set 4 hours and not 0 hours
func resetBucket(m *model) (tea.Model, tea.Cmd) {
	// Move the selection away if we're reseting the active bucket
	if m.cursor == m.selected {
		m.selected = -1
		m.activeSelection = false
	}
	// Pull the elapsed time out of the parent bucket
	m.buckets[m.buckets[m.cursor].parentBucket].elapsedTime -= m.buckets[m.cursor].elapsedTime
	storeBucketData(m.buckets[m.buckets[m.cursor].parentBucket], m.datastore)
	// Pull the elapsed time out of the total, if we happen to be operating on a second level bucket
	if m.buckets[m.cursor].level == "second" {
		m.buckets[0].elapsedTime -= m.buckets[m.cursor].elapsedTime
		storeBucketData(m.buckets[0], m.datastore)
	}
	// Set the actioned bucket time to 0
	m.buckets[m.cursor].elapsedTime = 0 * time.Second
	storeBucketData(m.buckets[m.cursor], m.datastore)
	return m, nil
}

func minuteUp(m *model) (tea.Model, tea.Cmd) {
	addIncrementTime(15*time.Minute, m)
	return m, nil
}

func minuteDown(m *model) (tea.Model, tea.Cmd) {
	addIncrementTime(-15*time.Minute, m)
	return m, nil
}

func hourUp(m *model) (tea.Model, tea.Cmd) {
	addIncrementTime(1*time.Hour, m)
	return m, nil
}

func hourDown(m *model) (tea.Model, tea.Cmd) {
	addIncrementTime(-1*time.Hour, m)
	return m, nil
}
