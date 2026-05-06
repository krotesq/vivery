package target

import "time"

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
	StreamKey string    `json:"stream_key"`
	CreatedAt time.Time `json:"createdAt"`
}

type targetWithDetailsDTO struct {
	Target targetDTO `json:"target"`
	// add all available types here
	Rtmp *rtmpDTO `json:"rtmp"`
}
