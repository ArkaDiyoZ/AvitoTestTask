package models

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserSegment struct {
	UserID    int `json:"user_id"`
	SegmentID int `json:"segment_id"`
}
