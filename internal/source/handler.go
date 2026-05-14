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
		response.Send(w, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	sourceModels, err := handler.service.findAll(r.Context(), accountID)
	if err != nil {
		response.Send(w, http.StatusInternalServerError, "Failed to retrieve sources", nil)
		return
	}

	sourceDTOs := make([]sourceDTO, 0, len(sourceModels))
	for _, sourceModel := range sourceModels {
		sourceDTOs = append(sourceDTOs, toSourceDTO(&sourceModel))
	}

	response.Send(w, http.StatusOK, "Sources retrieved", sourceDTOs)
}

func (handler *handler) findWithDetailsByID(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (handler *handler) findWithRtmpByID(w http.ResponseWriter, r *http.Request) {
	accountID, ok := auth.AccountIDFromContext(r.Context())
	if !ok {
		response.Send(w, http.StatusInternalServerError, "Internal server error", nil)
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
		response.Send(w, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	var sourceWithRtmpDTO createSourceWithRtmpDTO
	if err := util.ParseBody(r.Body, &sourceWithRtmpDTO); err != nil {
		response.Send(w, http.StatusBadRequest, "Failed to create source", nil)
		return
	}

	sourceModel, rtmpModel, err := handler.service.createWithRtmp(r.Context(), sourceWithRtmpDTO.Source.Name, sourceWithRtmpDTO.Source.Description, sourceWithRtmpDTO.Rtmp.URL, sourceWithRtmpDTO.Rtmp.StreamKey, accountID)
	if err != nil {
		response.Send(w, http.StatusBadRequest, "Failed to create source", nil)
		return
	}

	rtmpDto := toRtmpDTOWithStreamKey(rtmpModel)
	result := sourceWithDetailsDTO{
		Source: toSourceDTO(sourceModel),
		Rtmp:   &rtmpDto,
	}
	response.Send(w, http.StatusOK, "RTMP source created", result)
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
		response.Send(w, http.StatusNotFound, "Source not found", nil)
		return
	}

	response.Send(w, http.StatusOK, "Source deleted", nil)
}
