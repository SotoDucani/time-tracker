package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/peterbourgon/diskv/v3"
)

func diskv_init() *diskv.Diskv {
	flatTransform := func(s string) []string { return []string{} }

	d := diskv.New(diskv.Options{
		BasePath:     "data",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	return d
}

func sanitizeKeyName(key string) string {
	// lowercase and dashes only
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, " ", "-")
	emptySpace := []byte("")
	regexStr := regexp.MustCompile(`[^a-z-]`)
	key = string(regexStr.ReplaceAll([]byte(key), emptySpace))
	return key
}

func dateKeyName(key string) string {
	dateStr := time.Now().Local().Format("2006-01-02")
	key = fmt.Sprintf("%s-%s", dateStr, key)
	return key
}

func storeBucketData(bucket timeBucket, d *diskv.Diskv) error {
	key := sanitizeKeyName(bucket.name)
	key = dateKeyName(key)
	// we store elapsed in seconds
	elapsedTimeSeconds := bucket.elapsedTime.Seconds()
	storageString := strconv.Itoa(int(elapsedTimeSeconds))
	err := d.Write(key, []byte(storageString))
	if err != nil {
		log.Printf("Error writing data %v, %s", storageString, err)
		return err
	}
	return nil
}

func readBucketData(bucket timeBucket, d *diskv.Diskv) (time.Duration, error) {
	key := sanitizeKeyName(bucket.name)
	key = dateKeyName(key)
	value, err := d.Read(key)
	if err != nil {
		log.Printf("Error reading data from %s, %s", key, err)
		return 0 * time.Second, err
	}
	// force value to time.Duration as that's what we deal in
	elapsedTime, err := strconv.Atoi(string(value))
	if err != nil {
		log.Printf("Error converting data %v, %s", elapsedTime, err)
		return 0 * time.Second, err
	}
	return time.Duration(elapsedTime) * time.Second, err
}
