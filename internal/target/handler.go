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
		response.Send(w, http.StatusUnauthorized, "Unauthorized: missing account", nil)
		return
	}

	targetModels, err := handler.service.findAll(r.Context(), accountID)
	if err != nil {
		response.Send(w, http.StatusInternalServerError, "Failed to retrieve targets", nil)
		return
	}

	var targetDTOs []targetDTO
	for _, targetModel := range targetModels {
		targetDTOs = append(targetDTOs, toTargetDTO(&targetModel))
	}

	response.Send(w, http.StatusOK, "Targets retrieved", targetDTOs)
}

func (handler *handler) findWithRtmpByID(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		response.Send(w, http.StatusUnauthorized, "Unauthorized: missing account", nil)
		return
	}

	id := chi.URLParam(r, "id")

	targetModel, rtmpModel, err := handler.service.findWithRtmpByID(r.Context(), id, accountID)
	if err != nil {
		response.Send(w, http.StatusNotFound, "RTMP target not found", nil)
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
		response.Send(w, http.StatusUnauthorized, "Unauthorized: missing account", nil)
		return
	}

	var _targetWithDetailsDTO targetWithDetailsDTO
	if err := util.ParseBody(r.Body, &_targetWithDetailsDTO); err != nil {
		response.Send(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	targetModel, rtmpModel, err := handler.service.createWithRtmp(r.Context(), _targetWithDetailsDTO.Target.Name, _targetWithDetailsDTO.Target.Description, _targetWithDetailsDTO.Rtmp.URL, _targetWithDetailsDTO.Rtmp.StreamKey, accountID)
	if err != nil {
		response.Send(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	rtmpDto := toRtmpDTO(rtmpModel)
	result := targetWithDetailsDTO{
		Target: toTargetDTO(targetModel),
		Rtmp:   &rtmpDto,
	}
	response.Send(w, http.StatusOK, "RTMP target created", result)
}

func (handler *handler) deleteByID(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		response.Send(w, http.StatusUnauthorized, "Unauthorized: missing account", nil)
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
