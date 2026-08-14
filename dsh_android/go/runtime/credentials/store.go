package credentials

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Store manages credential references and values
type Store struct {
	mu      sync.RWMutex
	path    string
	data    map[string]Credential
	providers map[string]Provider
}

// Credential represents a stored credential
type Credential struct {
	Ref      string `json:"ref"`
	Value    string `json:"value,omitempty"`
	Source   string `json:"source"`
	Writable bool   `json:"writable"`
}

// Provider resolves credential references to values
type Provider interface {
	Name() string
	Resolve(ref string) (string, bool)
	Set(ref, value string) error
}

// NewStore creates a new credential store
func NewStore(dataDir string) (*Store, error) {
	path := filepath.Join(dataDir, "credentials.json")
	store := &Store{
		path:      path,
		data:      make(map[string]Credential),
		providers: make(map[string]Provider),
	}
	
	// Load existing credentials
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &store.data)
	}
	
	// Register default providers
	store.RegisterProvider(&EnvProvider{})
	store.RegisterProvider(&FileProvider{Dir: dataDir})
	
	return store, nil
}

// Get returns a credential by reference
func (s *Store) Get(ref string) (Credential, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.data[ref]
	return c, ok
}

// Set stores a credential value
func (s *Store) Set(ref, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	c := Credential{
		Ref:      ref,
		Value:    value,
		Source:   "local",
		Writable: true,
	}
	s.data[ref] = c
	return s.save()
}

// Resolve resolves a credential reference to its value
func (s *Store) Resolve(ref string) (string, bool) {
	s.mu.RLock()
	c, ok := s.data[ref]
	s.mu.RUnlock()
	
	if !ok {
		// Try providers
		for _, provider := range s.providers {
			if val, found := provider.Resolve(ref); found {
				return val, true
			}
		}
		return "", false
	}
	
	if c.Value != "" {
		return c.Value, true
	}
	
	// Try providers for stored refs
	for _, provider := range s.providers {
		if val, found := provider.Resolve(c.Ref); found {
			return val, true
		}
	}
	
	return "", false
}

// RegisterProvider adds a credential provider
func (s *Store) RegisterProvider(p Provider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providers[p.Name()] = p
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

// EnvProvider reads from environment variables
type EnvProvider struct{}

func (p *EnvProvider) Name() string { return "env" }

func (p *EnvProvider) Resolve(ref string) (string, bool) {
	val := os.Getenv(ref)
	return val, val != ""
}

func (p *EnvProvider) Set(ref, value string) error {
	os.Setenv(ref, value)
	return nil
}

// FileProvider reads from files
type FileProvider struct {
	Dir string
}

func (p *FileProvider) Name() string { return "file" }

func (p *FileProvider) Resolve(ref string) (string, bool) {
	path := filepath.Join(p.Dir, ref)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func (p *FileProvider) Set(ref, value string) error {
	path := filepath.Join(p.Dir, ref)
	return os.WriteFile(path, []byte(value), 0600)
}
