package server

import (
	"encoding/json"
	"net/http"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/logging"
	"cli-wrapper/internal/telemetry"
)

// Server exposes HTTP endpoints for the frontend.
type Server struct {
	mgr    *app.SessionManager
	logger *logging.Logger
	mux    *http.ServeMux
}

// New creates a Server with its routes configured.
func New(mgr *app.SessionManager, logger *logging.Logger) *Server {
	s := &Server{mgr: mgr, logger: logger, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/sessions", s.handleSessions)
	s.mux.HandleFunc("/resource", s.handleResource)
	s.mux.HandleFunc("/models", s.handleModels)
}

func (s *Server) respondJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, s.mgr.Sessions())
}

func (s *Server) handleResource(w http.ResponseWriter, r *http.Request) {
	u, err := telemetry.Read()
	if err != nil {
		s.logger.Error("read metrics: " + err.Error())
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	s.respondJSON(w, u)
}

func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	tools, err := app.DetectCLITools()
	if err != nil {
		s.logger.Error("detect models: " + err.Error())
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	s.respondJSON(w, tools)
}

// Start listens on the given address.
func (s *Server) Start(addr string) error {
	s.logger.Info("server starting on " + addr)
	return http.ListenAndServe(addr, s.mux)
}
