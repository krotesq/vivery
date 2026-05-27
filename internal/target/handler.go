package target

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/krotesq/vivery/internal/auth"
	"github.com/krotesq/vivery/internal/response"
	"github.com/krotesq/vivery/internal/util"
)

type handler struct {
	service *service
}

func newHandler(service *service) *handler {
	return &handler{service: service}
}

func (handler *handler) findAll(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		response.Send(w, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	targetModels, err := handler.service.findAll(r.Context(), accountID)
	if err != nil {
		response.Send(w, http.StatusInternalServerError, "Failed to retrieve targets", nil)
		return
	}

	targetDTOs := make([]targetDTO, 0, len(targetModels))
	for _, targetModel := range targetModels {
		targetDTOs = append(targetDTOs, toTargetDTO(&targetModel))
	}

	response.Send(w, http.StatusOK, "Targets retrieved", targetDTOs)
}

func (handler *handler) findWithRtmpByID(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		response.Send(w, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	id := chi.URLParam(r, "id")

	targetModel, rtmpModel, err := handler.service.findWithRtmpByID(r.Context(), id, accountID)
	if err != nil {
		response.Send(w, http.StatusNotFound, "Target not found", nil)
		return
	}

	rtmpDto := toRtmpDTO(rtmpModel)
	result := targetWithDetailsDTO{
		Target: toTargetDTO(targetModel),
		Rtmp:   &rtmpDto,
	}
	response.Send(w, http.StatusOK, "RTMP target found", result)
}

func (handler *handler) createWithRtmp(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		response.Send(w, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	var targetWithRtmpDTO createTargetWithRtmpDTO
	if err := util.ParseBody(r.Body, &targetWithRtmpDTO); err != nil {
		response.Send(w, http.StatusBadRequest, "Could not parse request body", nil)
		return
	}

	targetModel, rtmpModel, err := handler.service.createWithRtmp(r.Context(), targetWithRtmpDTO.Target.Name, targetWithRtmpDTO.Target.Description, targetWithRtmpDTO.Rtmp.URL, targetWithRtmpDTO.Rtmp.StreamKey, accountID)
	if err != nil {
		response.Send(w, http.StatusBadRequest, "Could not create target", nil)
		return
	}

	rtmpDto := toRtmpDTOWithStreamKey(rtmpModel)
	result := targetWithDetailsDTO{
		Target: toTargetDTO(targetModel),
		Rtmp:   &rtmpDto,
	}
	response.Send(w, http.StatusOK, "RTMP target created", result)
}

func (handler *handler) deleteByID(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		response.Send(w, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	id := chi.URLParam(r, "id")

	err := handler.service.deleteByID(r.Context(), id, accountID)
	if err != nil {
		response.Send(w, http.StatusNotFound, "Target not found", nil)
		return
	}

	response.Send(w, http.StatusOK, "Target deleted", nil)
}
