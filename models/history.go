package models

import "time"

type History struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	SegmentID int       `json:"segment_id"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}
