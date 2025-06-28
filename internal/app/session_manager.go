package app

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"cli-wrapper/internal/logging"
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
	ID     string
	Cmd    *exec.Cmd
	cancel context.CancelFunc
	done   chan error
}

// SessionManager schedules CLI sessions with a concurrency limit.
type SessionManager struct {
	baseDir string
	logger  *logging.Logger
	queue   chan sessionRequest
	active  map[string]*Session
	mu      sync.Mutex
	sem     chan struct{}
	ctx     context.Context
	cancel  context.CancelFunc
}

type sessionRequest struct {
	id   string
	tool string
	args []string
	res  chan error
}

// NewSessionManager creates a manager with the given concurrency limit.
func NewSessionManager(baseDir string, logger *logging.Logger, concurrency int) *SessionManager {
	ctx, cancel := context.WithCancel(context.Background())
	m := &SessionManager{
		baseDir: baseDir,
		logger:  logger,
		queue:   make(chan sessionRequest, 100),
		active:  make(map[string]*Session),
		sem:     make(chan struct{}, concurrency),
		ctx:     ctx,
		cancel:  cancel,
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

// run executes a session and tracks it until completion.
func (m *SessionManager) run(req sessionRequest) {
	ctx, cancel := context.WithCancel(m.ctx)
	cmd := exec.CommandContext(ctx, req.tool, req.args...)
	sess := &Session{ID: req.id, Cmd: cmd, cancel: cancel, done: make(chan error, 1)}
	m.mu.Lock()
	m.active[req.id] = sess
	m.mu.Unlock()

	if err := cmd.Start(); err != nil {
		m.logger.Error(fmt.Sprintf("start %s: %v", req.id, err))
		req.res <- fmt.Errorf("start: %w", err)
		m.remove(req.id)
		<-m.sem
		return
	}
	go func() {
		err := cmd.Wait()
		sess.done <- err
		m.logger.Info(fmt.Sprintf("session %s finished", req.id))
		m.remove(req.id)
		<-m.sem
	}()
	req.res <- nil
}

func (m *SessionManager) remove(id string) {
	m.mu.Lock()
	delete(m.active, id)
	m.mu.Unlock()
}

// AddSession queues a CLI invocation and returns its ID when scheduled.
func (m *SessionManager) AddSession(tool string, args []string) (string, error) {
	id, err := newUUID()
	if err != nil {
		return "", err
	}
	res := make(chan error, 1)
	req := sessionRequest{id: id, tool: tool, args: args, res: res}
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
	}
	m.active = make(map[string]*Session)
	m.mu.Unlock()
}
