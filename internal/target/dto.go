package target

import (
	"time"
)

type targetDTO struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TargetTypeID string    `json:"targetTypeId"`
	AccountID    string    `json:"accountId"`
	CreatedAt    time.Time `json:"createdAt"`
}

type rtmpDTO struct {
	TargetID  string    `json:"targetId"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}

type targetWithDetailsDTO struct {
	Target targetDTO `json:"target"`
	// add all available types here
	Rtmp *rtmpDTO `json:"rtmp"`
}

// ----------------------
// CREATE DTOs
// ----------------------

type createTargetDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type createRtmpDTO struct {
	URL       string `json:"url"`
	StreamKey string `json:"stream_key"`
}

type createTargetWithRtmpDTO struct {
	Target createTargetDTO `json:"target"`
	Rtmp   createRtmpDTO   `json:"rtmp"`
}
