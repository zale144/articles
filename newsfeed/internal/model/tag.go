package model

import "time"

type Card struct {
	Title     string    `bson:"title"`
	Timestamp time.Time `bson:"timestamp"`
	Tags      []string  `bson:"tags"`
}
