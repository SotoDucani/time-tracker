package main

import "time"

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
