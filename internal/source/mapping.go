package source

// takes an source and returns the source dto
func toSourceDTO(source *source) sourceDTO {
	return sourceDTO{
		ID:           source.ID,
		Name:         source.Name,
		Description:  source.Description,
		SourceTypeID: source.SourceTypeID,
		AccountID:    source.AccountID,
		CreatedAt:    source.CreatedAt,
	}
}

func toRtmpDTO(rtmp *rtmp) rtmpDTO {
	return rtmpDTO{
		SourceID:  rtmp.SourceID,
		URL:       rtmp.URL,
		StreamKey: "********",
		CreatedAt: rtmp.CreatedAt,
	}
}

func toRtmpDTOWithStreamKey(rtmp *rtmp) rtmpDTO {
	return rtmpDTO{
		SourceID:  rtmp.SourceID,
		URL:       rtmp.URL,
		StreamKey: rtmp.StreamKey,
		CreatedAt: rtmp.CreatedAt,
	}
}
