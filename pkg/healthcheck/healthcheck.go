package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/alexliesenfeld/health"
)

type Healthcheck struct {
	cfg *Config

	mux *http.ServeMux
	srv *http.Server

	mu       sync.RWMutex
	probes   map[string]*Probe
	handlers map[string]http.Handler
}

func New(cfg *Config) (*Healthcheck, error) {
	if cfg == nil {
		return nil, errors.New("healthcheck -> config is nil")
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("healthcheck -> failed to validate config -> %w", err)
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
		return errors.New("healthcheck -> probe is nil")
	}
	if p.Name() == "" {
		return errors.New("healthcheck -> probe name is empty")
	}
	route := p.Route()
	if route == "" || route[0] != '/' {
		return fmt.Errorf("healthcheck -> probe route must start with '/' -> %q", route)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.probes[p.Name()]; exists {
		return fmt.Errorf("healthcheck -> probe with same name already registered -> %s", p.Name())
	}
	h.probes[p.Name()] = p

	if _, exists := h.handlers[route]; exists {
		return fmt.Errorf("healthcheck -> route already registered -> %s", route)
	}

	checker := health.NewChecker(
		health.WithCheck(health.Check{
			Name: p.Name(),
			Check: func(ctx context.Context) error {
				if p.IsEnabled() {
					return nil
				}
				return errors.New("healthcheck -> probe disabled")
			},
		}),
	)

	handler := health.NewHandler(checker)
	h.handlers[route] = handler
	h.mux.Handle(route, handler)

	return nil
}

func (h *Healthcheck) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		err := h.srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			errCh <- nil
			return
		}
		errCh <- err
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), h.cfg.ShutdownTimeout)
		defer cancel()

		err := h.srv.Shutdown(shutdownCtx)
		_ = <-errCh

		return err
	}
}
