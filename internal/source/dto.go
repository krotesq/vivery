package source

import "time"

type sourceDTO struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	SourceTypeID string    `json:"sourceTypeId"`
	AccountID    string    `json:"accountId"`
	CreatedAt    time.Time `json:"createdAt"`
}

type rtmpDTO struct {
	SourceID  string    `json:"sourceId"`
	URL       string    `json:"url"`
	StreamKey string    `json:"stream_key"`
	CreatedAt time.Time `json:"createdAt"`
}

type sourceWithDetailsDTO struct {
	Source sourceDTO `json:"source"`
	// add all available types here
	Rtmp *rtmpDTO `json:"rtmp"`
}
