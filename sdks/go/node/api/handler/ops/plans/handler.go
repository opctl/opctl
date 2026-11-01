package plans

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

type Handler interface {
	Handle(http.ResponseWriter, *http.Request)
}

func NewHandler(node node.Node) Handler { return handler{node: node} }

type handler struct{ node node.Node }

func (h handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var request model.PlanOpReq
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	spec, err := h.node.PlanOp(r.Context(), request)
	if err != nil {
		status := http.StatusInternalServerError
		var declaration *node.OpDeclarationError
		switch {
		case errors.Is(err, model.ErrDataProviderAuthentication{}):
			status = http.StatusUnauthorized
		case errors.Is(err, model.ErrDataProviderAuthorization{}):
			status = http.StatusForbidden
		case errors.Is(err, model.ErrDataRefResolution{}), errors.Is(err, model.ErrDataUnableToResolve{}):
			status = http.StatusNotFound
		case errors.Is(err, model.ErrDataMissingVersion{}), errors.Is(err, model.ErrDataGitInvalidRef{}), errors.As(err, &declaration):
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spec); err != nil {
		return
	}
}
