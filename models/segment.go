package models

type Segment struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

type SegmentsRequest struct {
	Segment []string `json:"segments"`
}
