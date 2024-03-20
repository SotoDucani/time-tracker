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

func addElapsedTime(bucket *timeBucket, startTime time.Time, m *model) {
	bucket.elapsedTime += time.Since(bucket.startTime)
	storeBucketData(*bucket, m.datastore)
	bucket.startTime = startTime
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
