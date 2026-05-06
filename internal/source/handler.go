package source

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

	sourceModels, err := handler.service.findAll(r.Context(), accountID)
	if err != nil {
		response.Send(w, http.StatusInternalServerError, "Failed to retrieve sources", nil)
		return
	}

	var sourceDTOs []sourceDTO
	for _, sourceModel := range sourceModels {
		sourceDTOs = append(sourceDTOs, toSourceDTO(&sourceModel))
	}

	response.Send(w, http.StatusOK, "Sources retrieved", sourceDTOs)
}

func (handler *handler) findWithRtmpByID(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		response.Send(w, http.StatusUnauthorized, "Unauthorized: missing account", nil)
		return
	}

	id := chi.URLParam(r, "id")

	sourceModel, rtmpModel, err := handler.service.findWithRtmpByID(r.Context(), id, accountID)
	if err != nil {
		response.Send(w, http.StatusNotFound, "RTMP source not found", nil)
		return
	}

	rtmpDto := toRtmpDTO(rtmpModel)
	result := sourceWithDetailsDTO{
		Source: toSourceDTO(sourceModel),
		Rtmp:   &rtmpDto,
	}
	response.Send(w, http.StatusOK, "RTMP source found", result)
}

func (handler *handler) createWithRtmp(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		response.Send(w, http.StatusUnauthorized, "Unauthorized: missing account", nil)
		return
	}

	var _sourceWithDetailsDTO sourceWithDetailsDTO
	if err := util.ParseBody(r.Body, &_sourceWithDetailsDTO); err != nil {
		response.Send(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	sourceModel, rtmpModel, err := handler.service.createWithRtmp(r.Context(), _sourceWithDetailsDTO.Source.Name, _sourceWithDetailsDTO.Source.Description, _sourceWithDetailsDTO.Rtmp.URL, _sourceWithDetailsDTO.Rtmp.StreamKey, accountID)
	if err != nil {
		response.Send(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	rtmpDto := toRtmpDTO(rtmpModel)
	result := sourceWithDetailsDTO{
		Source: toSourceDTO(sourceModel),
		Rtmp:   &rtmpDto,
	}
	response.Send(w, http.StatusOK, "RTMP source created", result)
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
		response.Send(w, http.StatusNotFound, "Source not found", nil)
		return
	}

	response.Send(w, http.StatusOK, "Source deleted", nil)
}
