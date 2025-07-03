package app

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"cli-wrapper/internal/history"
	"cli-wrapper/internal/logging"
	"cli-wrapper/internal/telemetry"
	"github.com/shirou/gopsutil/v3/process"
)

// newUUID returns a random UUID string.
func newUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// Session represents a running CLI invocation.
type Session struct {
	ID        string
	Model     string
	Prompt    string
	Cmd       *exec.Cmd
	cancel    context.CancelFunc
	done      chan error
	procMu    sync.Mutex
	proc      *process.Process
	bufMu     sync.Mutex
	throttled bool
	output    chan string
	metricsMu sync.Mutex
	metrics   []history.Metric
}

// SessionManager schedules CLI sessions with a concurrency limit.
type SessionManager struct {
	logger *logging.Logger
	queue  chan sessionRequest
	active map[string]*Session
	mu     sync.Mutex
	sem    chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *Config
	hist   *history.Store
}

type sessionRequest struct {
	id     string
	model  string
	prompt string
	args   []string
	res    chan error
}

// NewSessionManager creates a manager with the given concurrency limit.
func NewSessionManager(logger *logging.Logger, concurrency int, cfg *Config, hist *history.Store) *SessionManager {
	ctx, cancel := context.WithCancel(context.Background())
	m := &SessionManager{
		logger: logger,
		queue:  make(chan sessionRequest, 100),
		active: make(map[string]*Session),
		sem:    make(chan struct{}, concurrency),
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
		hist:   hist,
	}
	go m.dispatch()
	return m
}

func (m *SessionManager) dispatch() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case req := <-m.queue:
			m.sem <- struct{}{}
			go m.run(req)
		}
	}
}

func (m *SessionManager) streamOutput(sess *Session, r io.ReadCloser, buf *bytes.Buffer, ch chan<- string) {
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		sess.bufMu.Lock()
		buf.WriteString(line + "\n")
		sess.bufMu.Unlock()
		select {
		case ch <- line:
		case <-m.ctx.Done():
			return
		}
	}
}

// run executes a session and tracks it until completion.
func (m *SessionManager) run(req sessionRequest) {
	ctx, cancel := context.WithCancel(m.ctx)
	cmd := exec.CommandContext(ctx, req.model, req.args...)
	cmd.Dir = m.cfg.WorkingDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		req.res <- fmt.Errorf("stdout pipe: %w", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		req.res <- fmt.Errorf("stderr pipe: %w", err)
		return
	}

	var buf bytes.Buffer
	sess := &Session{ID: req.id, Model: req.model, Prompt: req.prompt, Cmd: cmd, cancel: cancel, done: make(chan error, 1), output: make(chan string, 32)}
	m.logger.Info(fmt.Sprintf("session %s running model %s", req.id, req.model))
	m.mu.Lock()
	m.active[req.id] = sess
	m.mu.Unlock()

	if err := cmd.Start(); err != nil {
		m.logger.Error(fmt.Sprintf("start %s: %v", req.id, err))
		cancel()
		req.res <- fmt.Errorf("start: %w", err)
		m.remove(req.id)
		<-m.sem
		return
	}
	go m.streamOutput(sess, stdout, &buf, sess.output)
	go m.streamOutput(sess, stderr, &buf, sess.output)
	proc, err := process.NewProcess(int32(cmd.Process.Pid))
	if err == nil {
		sess.procMu.Lock()
		sess.proc = proc
		sess.procMu.Unlock()
		go m.monitor(sess)
	} else {
		m.logger.Error(fmt.Sprintf("session %s monitor: %v", req.id, err))
	}
	go func() {
		err := cmd.Wait()
		sess.done <- err
		sess.procMu.Lock()
		if sess.proc != nil {
			telemetry.ReleaseProcess(sess.proc)
			sess.proc = nil
		}
		sess.procMu.Unlock()
		m.logger.Info(fmt.Sprintf("session %s finished", req.id))
		if m.hist != nil {
			sess.bufMu.Lock()
			resp := buf.String()
			sess.bufMu.Unlock()
			rec := history.Record{
				ID:       req.id,
				Prompt:   req.prompt,
				Response: resp,
				Model:    req.model,
				Success:  err == nil,
			}
			if err := m.hist.Add(rec); err != nil {
				m.logger.Error("history add " + err.Error())
			}
			sess.metricsMu.Lock()
			metrics := append([]history.Metric(nil), sess.metrics...)
			sess.metricsMu.Unlock()
			for _, mtr := range metrics {
				if err := m.hist.AddMetric(mtr); err != nil {
					m.logger.Error("metric add " + err.Error())
				}
			}
		}
		m.remove(req.id)
		close(sess.output)
		<-m.sem
	}()
	req.res <- nil
}

func (m *SessionManager) monitor(s *Session) {
	ticker := time.NewTicker(time.Duration(m.cfg.PollInterval) * time.Second)
	defer func() {
		ticker.Stop()
		s.procMu.Lock()
		s.procMu.Lock()
		if s.proc != nil {
			telemetry.ReleaseProcess(s.proc)
			s.proc = nil
		}
		s.procMu.Unlock()
		s.procMu.Unlock()
	}()
	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			s.procMu.Lock()
			proc := s.proc
			s.procMu.Unlock()
			if proc == nil {
				continue
			}
			usage, err := telemetry.ReadProcess(proc)
			if err != nil {
				m.logger.Error(fmt.Sprintf("metrics %s: %v", s.ID, err))
				continue
			}
			s.metricsMu.Lock()
			s.metrics = append(s.metrics, history.Metric{
				ID:        s.ID,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				CPU:       usage.CPU,
				Memory:    usage.Memory,
			})
			s.metricsMu.Unlock()
			if usage.CPU > m.cfg.CPUThreshold || usage.Memory > m.cfg.MemoryThreshold {
				m.logger.Info(fmt.Sprintf("session %s high usage cpu %.1f mem %.1f", s.ID, usage.CPU, usage.Memory))
				if !s.throttled {
					if err := throttleProcess(s.Cmd); err != nil {
						m.logger.Error("throttle " + s.ID + ": " + err.Error())
					} else {
						s.throttled = true
						m.logger.Info("session " + s.ID + " throttled")
					}
				} else {
					_ = m.Terminate(s.ID)
					m.logger.Info("session " + s.ID + " cancelled")
					return
				}
			}
		}
	}
}

func (m *SessionManager) remove(id string) {
	m.mu.Lock()
	delete(m.active, id)
	m.mu.Unlock()
}

// OutputChannel returns a channel of output lines for the given session.
func (m *SessionManager) OutputChannel(id string) (<-chan string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.active[id]
	if !ok {
		return nil, fmt.Errorf("unknown session %s", id)
	}
	return sess.output, nil
}

// AddSession queues a CLI invocation and returns its ID when scheduled.
func (m *SessionManager) AddSession(model, prompt string, args []string) (string, error) {
	var err error
	model, err = SanitizeModel(model)
	if err != nil {
		return "", err
	}
	prompt = SanitizePrompt(prompt)
	for i, a := range args {
		args[i] = SanitizePrompt(a)
	}
	id, err := newUUID()
	if err != nil {
		return "", err
	}
	res := make(chan error, 1)
	req := sessionRequest{id: id, model: model, prompt: prompt, args: args, res: res}
	select {
	case m.queue <- req:
	case <-time.After(5 * time.Second):
		return "", fmt.Errorf("queue timeout")
	}
	if err := <-res; err != nil {
		return "", err
	}
	return id, nil
}

// Terminate stops a running session by ID.
func (m *SessionManager) Terminate(id string) error {
	m.mu.Lock()
	sess, ok := m.active[id]
	m.mu.Unlock()
	if !ok {
		return fmt.Errorf("unknown session %s", id)
	}
	if err := terminateProcess(sess.Cmd); err != nil {
		if !errors.Is(err, os.ErrProcessDone) {
			return err
		}
	}
	sess.cancel()
	<-sess.done
	sess.procMu.Lock()
	if sess.proc != nil {
		telemetry.ReleaseProcess(sess.proc)
		sess.proc = nil
	}
	sess.procMu.Unlock()
	return nil
}

func terminateProcess(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return fmt.Errorf("process not started")
	}
	if runtime.GOOS == "windows" {
		return cmd.Process.Kill()
	}
	_ = cmd.Process.Signal(syscall.SIGTERM)
	time.Sleep(100 * time.Millisecond)
	return cmd.Process.Kill()
}

// Sessions returns the IDs of running sessions.
func (m *SessionManager) Sessions() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := make([]string, 0, len(m.active))
	for id := range m.active {
		ids = append(ids, id)
	}
	return ids
}

// Close stops all workers and running sessions.
func (m *SessionManager) Close() {
	m.cancel()
	m.mu.Lock()
	for _, s := range m.active {
		terminateProcess(s.Cmd)
		s.cancel()
		<-s.done
		if s.proc != nil {
			telemetry.ReleaseProcess(s.proc)
			s.proc = nil
		}
	}
	m.active = make(map[string]*Session)
	m.mu.Unlock()
}
