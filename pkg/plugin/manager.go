package plugin

import (
	"errors"
	"fmt"
	"plugin"
	"sync"
)

// Plugin represents a plugin that can be loaded and used by the application
type Plugin interface {
	// Initialize initializes the plugin with configuration
	Initialize(config map[string]interface{}) error

	// Name returns the plugin's name
	Name() string

	// Version returns the plugin's version
	Version() string

	// Type returns the plugin's type (e.g., "framework", "template")
	Type() string
}

// Manager manages plugins
type Manager struct {
	plugins map[string]Plugin
	mutex   sync.RWMutex
}

// NewManager creates a new plugin manager
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
		mutex:   sync.RWMutex{},
	}
}

// LoadPlugin loads a plugin from the given path
func (m *Manager) LoadPlugin(path string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Load the plugin
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look up the Plugin symbol
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin does not export 'Plugin' symbol: %w", err)
	}

	// Assert that it implements the Plugin interface
	plug, ok := sym.(Plugin)
	if !ok {
		return errors.New("plugin does not implement Plugin interface")
	}

	// Initialize the plugin
	err = plug.Initialize(nil)
	if err != nil {
		return fmt.Errorf("failed to initialize plugin: %w", err)
	}

	// Add it to our map
	m.plugins[plug.Name()] = plug

	return nil
}

// GetPlugin returns a plugin by name
func (m *Manager) GetPlugin(name string) (Plugin, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	plug, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", name)
	}

	return plug, nil
}

// ListPlugins returns a list of all loaded plugins
func (m *Manager) ListPlugins() []Plugin {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	plugs := make([]Plugin, 0, len(m.plugins))
	for _, plug := range m.plugins {
		plugs = append(plugs, plug)
	}

	return plugs
}