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
	mgr     *app.SessionManager
	logger  *logging.Logger
	baseDir string
	cfg     *app.Config
	mux     *http.ServeMux
}

// New creates a Server with its routes configured.
func New(mgr *app.SessionManager, logger *logging.Logger, baseDir string, cfg *app.Config) *Server {
	s := &Server{mgr: mgr, logger: logger, baseDir: baseDir, cfg: cfg, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/sessions", s.handleSessions)
	s.mux.HandleFunc("/resource", s.handleResource)
	s.mux.HandleFunc("/models", s.handleModels)
	s.mux.HandleFunc("/theme", s.handleTheme)
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

func (s *Server) handleTheme(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.respondJSON(w, map[string]string{"theme": s.cfg.Theme})
	case http.MethodPost:
		var req struct {
			Theme string `json:"theme"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if !app.ValidTheme(req.Theme) {
			http.Error(w, "invalid theme", http.StatusBadRequest)
			return
		}
		s.cfg.Theme = req.Theme
		if err := app.SaveConfig(s.baseDir, *s.cfg); err != nil {
			s.logger.Error("save theme: " + err.Error())
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		s.respondJSON(w, map[string]string{"theme": s.cfg.Theme})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Start listens on the given address.
func (s *Server) Start(addr string) error {
	s.logger.Info("server starting on " + addr)
	return http.ListenAndServe(addr, s.mux)
}
