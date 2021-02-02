package probe

import (
	"net/http"
)

type healthRouter struct {
	*http.ServeMux
}

func NewHealthRouter() http.Handler {
	mux := http.NewServeMux()
	r := &healthRouter{mux}

	r.HandleFunc("/healthz", r.liveness)

	return r
}

func (h *healthRouter) liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
