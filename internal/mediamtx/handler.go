package mediamtx

import (
	"fmt"
	"net/http"

	"github.com/krotesq/vivery/internal/response"
	"github.com/krotesq/vivery/internal/util"
)

type handler struct {
	service *service
}

func newHandler(service *service) *handler {
	return &handler{service: service}
}

func (handler *handler) auth(w http.ResponseWriter, r *http.Request) {
	var dto authDTO
	if err := util.ParseBody(r.Body, &dto); err != nil {
		response.Send(w, http.StatusBadRequest, "Could not authenticate", nil)
		return
	}

	fmt.Println("User: ", dto.User)
	fmt.Println("Password: ", dto.Password)

	response.Send(w, http.StatusOK, "test", nil)
}
