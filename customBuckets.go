package main

// Hackjob instead of a config file for now
// Misses out on a lot of features and is fragile
// ---------------
// Order matters both for buckets in the array as well as IDs
// First bucket must always be Total
// Yes you must start at 0
// Do not name ANY buckets the same
// ---------------
// Total          id 0
// |-> Parent     id 1
//     |-> Child  id 2
// |-> Parent     id 3

func setupBuckets() []timeBucket {
	return []timeBucket{
		{
			id:    0,
			name:  "Total",
			level: "none",
		},
		{
			id:           1,
			name:         "Parent Bucket One",
			level:        "first",
			parentBucket: 0,
		},
		{
			id:           2,
			name:         "Child Bucket One",
			level:        "second",
			parentBucket: 1,
		},
		{
			id:           3,
			name:         "Child Bucket Two",
			level:        "second",
			parentBucket: 1,
		},
		{
			id:           4,
			name:         "Parent Bucket Two",
			level:        "first",
			parentBucket: 0,
		},
		{
			id:           5,
			name:         "Child Bucket Three",
			level:        "second",
			parentBucket: 4,
		},
		{
			id:           6,
			name:         "Parent Bucket Three",
			level:        "first",
			parentBucket: 0,
		},
	}
}
