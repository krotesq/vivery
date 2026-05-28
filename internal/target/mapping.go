package target

// takes a target and returns the target DTO
func toTargetDTO(target *target) targetDTO {
	return targetDTO{
		ID:           target.ID,
		Name:         target.Name,
		Description:  target.Description,
		TargetTypeID: target.TargetTypeID,
		AccountID:    target.AccountID,
		CreatedAt:    target.CreatedAt,
	}
}

func toRtmpDTO(rtmp *rtmp) rtmpDTO {
	return rtmpDTO{
		TargetID:  rtmp.TargetID,
		URL:       rtmp.URL,
		StreamKey: "********",
		CreatedAt: rtmp.CreatedAt,
	}
}

func toRtmpDTOWithStreamKey(rtmp *rtmp) rtmpDTO {
	return rtmpDTO{
		TargetID:  rtmp.TargetID,
		URL:       rtmp.URL,
		StreamKey: rtmp.StreamKey,
		CreatedAt: rtmp.CreatedAt,
	}
}
