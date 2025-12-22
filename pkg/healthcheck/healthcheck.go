package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/alexliesenfeld/health"
)

type Healthcheck struct {
	cfg Config

	mux *http.ServeMux
	srv *http.Server

	mu       sync.RWMutex
	probes   map[string]*Probe
	handlers map[string]http.Handler
}

func New(cfg Config) (*Healthcheck, error) {
	if cfg.Addr == "" {
		return nil, errors.New("healthcheck: addr is empty")
	}

	h := &Healthcheck{
		cfg:      cfg,
		mux:      http.NewServeMux(),
		probes:   make(map[string]*Probe),
		handlers: make(map[string]http.Handler),
	}

	h.srv = &http.Server{
		Addr:    cfg.Addr,
		Handler: h.mux,
	}

	return h, nil
}

func (h *Healthcheck) Register(p *Probe) error {
	if p == nil {
		return errors.New("healthcheck: probe is nil")
	}
	if p.Name() == "" {
		return errors.New("healthcheck: probe name is empty")
	}
	route := p.Route()
	if route == "" || route[0] != '/' {
		return errors.New("healthcheck: probe route must start with '/'")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.probes[p.Name()]; exists {
		return errors.New("healthcheck: probe with same name already registered: " + p.Name())
	}
	h.probes[p.Name()] = p

	if _, exists := h.handlers[route]; exists {
		return errors.New("healthcheck: route already registered: " + route)
	}

	checker := health.NewChecker(
		health.WithCheck(health.Check{
			Name: p.Name(),
			Check: func(ctx context.Context) error {
				if p.IsEnabled() {
					return nil
				}
				return errors.New("probe disabled")
			},
		}),
	)

	handler := health.NewHandler(checker)
	h.handlers[route] = handler
	h.mux.Handle(route, handler)

	return nil
}

func (h *Healthcheck) Run() error {
	err := h.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (h *Healthcheck) Shutdown(ctx context.Context) error {
	return h.srv.Shutdown(ctx)
}
