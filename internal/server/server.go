package server

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"os"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/history"
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
	hist    *history.Store
}

// New creates a Server with its routes configured.
func New(mgr *app.SessionManager, logger *logging.Logger, baseDir string, cfg *app.Config, hist *history.Store) *Server {
	s := &Server{mgr: mgr, logger: logger, baseDir: baseDir, cfg: cfg, hist: hist, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/sessions", s.handleSessions)
	s.mux.Handle("/stream", websocket.Handler(s.handleStream))
	s.mux.HandleFunc("/resource", s.handleResource)
	s.mux.HandleFunc("/models", s.handleModels)
	s.mux.HandleFunc("/billing", s.handleBilling)
	s.mux.HandleFunc("/theme", s.handleTheme)
	s.mux.HandleFunc("/config", s.handleConfig)
	s.mux.HandleFunc("/history", s.handleHistory)
	s.mux.HandleFunc("/history/search", s.handleHistorySearch)
	s.mux.HandleFunc("/history/export", s.handleHistoryExport)
}

func (s *Server) respondJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, s.mgr.Sessions())
}

func (s *Server) handleStream(ws *websocket.Conn) {
	id := ws.Request().URL.Query().Get("id")
	if id == "" {
		ws.Close()
		return
	}
	ch, err := s.mgr.OutputChannel(id)
	if err != nil {
		ws.Close()
		return
	}
	for line := range ch {
		if err := websocket.Message.Send(ws, line); err != nil {
			break
		}
	}
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

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resp := map[string]any{
			"concurrency": s.cfg.Concurrency,
			"workingDir":  s.cfg.WorkingDir,
			"logLevel":    s.cfg.LogLevel,
			"logPath":     s.cfg.LogPath,
		}
		s.respondJSON(w, resp)
	case http.MethodPost:
		var req struct {
			Concurrency int    `json:"concurrency"`
			WorkingDir  string `json:"workingDir"`
			LogLevel    string `json:"logLevel"`
			LogPath     string `json:"logPath"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if req.Concurrency < 1 || req.Concurrency > 5 {
			http.Error(w, "invalid concurrency", http.StatusBadRequest)
			return
		}
		if req.WorkingDir == "" {
			http.Error(w, "working directory required", http.StatusBadRequest)
			return
		}
		if info, err := os.Stat(req.WorkingDir); err != nil || !info.IsDir() {
			http.Error(w, "invalid working directory", http.StatusBadRequest)
			return
		}
		switch req.LogLevel {
		case "error", "info", "debug":
		default:
			http.Error(w, "invalid log level", http.StatusBadRequest)
			return
		}
		if req.LogPath == "" {
			http.Error(w, "log path required", http.StatusBadRequest)
			return
		}
		s.cfg.Concurrency = req.Concurrency
		s.cfg.WorkingDir = req.WorkingDir
		s.cfg.LogLevel = req.LogLevel
		s.cfg.LogPath = req.LogPath
		if err := app.SaveConfig(s.baseDir, *s.cfg); err != nil {
			s.logger.Error("save config: " + err.Error())
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		resp := map[string]any{"concurrency": s.cfg.Concurrency, "workingDir": s.cfg.WorkingDir, "logLevel": s.cfg.LogLevel, "logPath": s.cfg.LogPath}
		s.respondJSON(w, resp)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		recs, err := s.hist.All()
		if err != nil {
			s.logger.Error("history all: " + err.Error())
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		s.respondJSON(w, recs)
	case http.MethodPost:
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if err := s.hist.Import(data); err != nil {
			s.logger.Error("import history: " + err.Error())
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		s.respondJSON(w, map[string]string{"status": "ok"})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleHistorySearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	recs, err := s.hist.Search(q, 20)
	if err != nil {
		s.logger.Error("search history: " + err.Error())
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	s.respondJSON(w, recs)
}

func (s *Server) handleHistoryExport(w http.ResponseWriter, r *http.Request) {
	data, err := s.hist.Export()
	if err != nil {
		s.logger.Error("export history: " + err.Error())
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (s *Server) handleBilling(w http.ResponseWriter, r *http.Request) {
	tool, err := app.DetectCLITool()
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	url, err := app.BillingURL(tool)
	if err != nil {
		s.logger.Error("billing url: " + err.Error())
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	resp := map[string]any{
		"tool": tool,
		"url":  url,
	}
	if usage, err := app.FetchUsage(tool); err == nil {
		resp["usage"] = usage
	}
	s.respondJSON(w, resp)
}

// Start listens on the given address.
func (s *Server) Start(addr string) error {
	s.logger.Info("server starting on " + addr)
	return http.ListenAndServe(addr, s.mux)
}
