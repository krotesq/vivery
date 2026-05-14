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
	StreamKey string    `json:"streamKey"`
	CreatedAt time.Time `json:"createdAt"`
}

type sourceWithDetailsDTO struct {
	Source sourceDTO `json:"source"`
	// add all available types here
	Rtmp *rtmpDTO `json:"rtmp"`
}

// ----------------------
// CREATE DTOs
// ----------------------

type createSourceDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type createRtmpDTO struct {
	URL       string `json:"url"`
	StreamKey string `json:"stream_key"`
}

type createSourceWithRtmpDTO struct {
	Source createSourceDTO `json:"source"`
	Rtmp   createRtmpDTO   `json:"rtmp"`
}
