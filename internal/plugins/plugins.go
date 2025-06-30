package plugins

import (
	"sort"
	"sync"
)

// BillingUsage holds credit usage information returned by a CLI.
type BillingUsage struct {
	TotalGranted   float64 `json:"total_granted"`
	TotalUsed      float64 `json:"total_used"`
	TotalAvailable float64 `json:"total_available"`
}

// Plugin defines the interface all CLI plugins must implement.
type Plugin interface {
	Name() string
	Detect() bool
	Invoke(args []string) error
	BillingURL() string
	Usage() (*BillingUsage, error)
}

var (
	mu       sync.RWMutex
	registry = make(map[string]Plugin)
)

// Register adds a plugin to the registry.
func Register(p Plugin) {
	mu.Lock()
	defer mu.Unlock()
	registry[p.Name()] = p
}

// Get retrieves a plugin by name.
func Get(name string) (Plugin, bool) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := registry[name]
	return p, ok
}

// All returns all registered plugins sorted by name.
func All() []Plugin {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	res := make([]Plugin, 0, len(names))
	for _, n := range names {
		res = append(res, registry[n])
	}
	return res
}
