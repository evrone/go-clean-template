package router

import (
	"fmt"
	"net/http"
)

type ProbeRouter struct {
	*http.ServeMux
}

func NewProbeRouter() *ProbeRouter {
	mux := http.NewServeMux()
	p := &ProbeRouter{mux}

	p.HandleFunc("/healthz", p.Liveness)

	return p
}

func (p *ProbeRouter) Liveness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("It's alive!")
	w.WriteHeader(http.StatusOK)
}
