package account

import (
	"net/http"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/krotesq/vivery/internal/response"
	"github.com/krotesq/vivery/internal/util"
)

type handler struct {
	s *service
}

func newHandler(s *service) *handler {
	return &handler{s: s}
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	var dto createDTO
	if err := util.ParseBody(r.Body, &dto); err != nil {
		response.Send(w, http.StatusBadRequest, "Could not create account", nil)
		return
	}

	acc, err := h.s.create(r.Context(), dto.Username, dto.Password)
	if err != nil {
		response.Send(w, http.StatusBadRequest, "Could not create account", nil)
		return
	}

	response.Send(w, http.StatusOK, "Account created", toAccountDTO(acc))
}

func (h *handler) findByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	acc, err := h.s.findByID(r.Context(), id)
	if err != nil {
		response.Send(w, http.StatusNotFound, "Account not found", nil)
		return
	}

	response.Send(w, http.StatusOK, "Account found", toAccountDTO(acc))
}

func (h *handler) deleteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	acc, err := h.s.deleteByID(r.Context(), id)
	if err != nil {
		response.Send(w, http.StatusNotFound, "Account not found", nil)
		return
	}

	response.Send(w, http.StatusOK, "Account deleted", toAccountDTO(acc))
}

func (h *handler) deactivateByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	acc, err := h.s.deactivateByID(r.Context(), id)
	if err != nil {
		response.Send(w, http.StatusNotFound, "Account not found", nil)
		return
	}

	response.Send(w, http.StatusOK, "Account deactivated", toAccountDTO(acc))
}

func (h *handler) activateByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	acc, err := h.s.activateByID(r.Context(), id)
	if err != nil {
		response.Send(w, http.StatusNotFound, "Account not found", nil)
		return
	}

	response.Send(w, http.StatusOK, "Account activated", toAccountDTO(acc))
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var dto loginDTO
	if err := util.ParseBody(r.Body, &dto); err != nil {
		response.Send(w, http.StatusBadRequest, "Could not login", nil)
		return
	}

	ip := util.GetClientIP(r)

	acc, accessToken, refreshToken, err := h.s.login(r.Context(), dto.Username, dto.Password, r.UserAgent(), ip)
	if err != nil {
		response.Send(w, http.StatusBadRequest, "Could not login", nil)
		return
	}

	res := response.NewBuilder(w)
	res.SetStatus(201)
	res.SetSimpleCookie("access_token", accessToken)
	res.SetSimpleCookie("refresh_token", refreshToken)
	res.SetBody("Account logged in", toAccountDTO(acc))
	res.Send()
}

func (h *handler) refresh(w http.ResponseWriter, r *http.Request) {
	// get cookie
	cookie, err := r.Cookie("refresh_token")
	if errors.Is(err, http.ErrNoCookie) {
		response.Send(w, http.StatusUnauthorized, "Could not refresh", nil)
		return
	}
	if err != nil {
		response.Send(w, http.StatusInternalServerError, "Could not refresh", nil)
		return
	}

	// create new refresh token
	ip := util.GetClientIP(r)
	refreshToken, jwt, err := h.s.refresh(r.Context(), cookie.Value, r.UserAgent(), ip)
	if err != nil {
		response.Send(w, 500, "Could not refresh", nil)
		return
	}

	// TODO: zu res muss expire hinzugefügt werden (refreshToken.expires_at == cookie.expire)
	res := response.NewBuilder(w)
	res.SetStatus(201)
	res.SetSimpleCookie("access_token", jwt)
	res.SetSimpleCookie("refresh_token", refreshToken)
	res.SetBody("Token refreshed", nil)
	res.Send()
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	// revoke refresh token in db if cookie exists
	cookie, err := r.Cookie("refresh_token")
	if err == nil && cookie.Value != "" {
		h.s.revokeRefreshToken(r.Context(), cookie.Value)
	}

	res := response.NewBuilder(w)
	res.DeleteCookie("access_token")
	res.DeleteCookie("refresh_token")
	res.SetStatus(http.StatusOK)
	res.SetBody("Account logged out.", nil)
	res.Send()
}

func (h *handler) me(w http.ResponseWriter, r *http.Request) {
	acc, err := h.s.me(r.Context())
	if err != nil {
		response.Send(w, http.StatusInternalServerError, "Unable to load account.", nil)
		return
	}
	response.Send(w, http.StatusOK, "Account found.", toAccountDTO(acc))
}
