package weewar

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// =============================================================================
// Memory SaveHandler (for testing)
// =============================================================================

// MemorySaveHandler stores sessions in memory - useful for testing
type MemorySaveHandler struct {
	sessions map[string][]byte
	mutex    sync.RWMutex
}

// NewMemorySaveHandler creates a new in-memory save handler
func NewMemorySaveHandler() *MemorySaveHandler {
	return &MemorySaveHandler{
		sessions: make(map[string][]byte),
	}
}

func (h *MemorySaveHandler) Save(sessionData []byte) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Extract session ID from the data to use as key
	// For now, we'll use a simple approach - in real implementation,
	// we'd parse the JSON to get the sessionId
	sessionID := fmt.Sprintf("session_%d", len(h.sessions))
	h.sessions[sessionID] = sessionData

	return nil
}

// Load, List, Delete not needed - UI handles these operations directly

// =============================================================================
// File SaveHandler (for CLI mode)
// =============================================================================

// FileSaveHandler stores sessions as JSON files on disk
type FileSaveHandler struct {
	saveDirectory string
}

// NewFileSaveHandler creates a new file-based save handler
func NewFileSaveHandler(saveDirectory string) (*FileSaveHandler, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(saveDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create save directory: %w", err)
	}

	return &FileSaveHandler{
		saveDirectory: saveDirectory,
	}, nil
}

func (h *FileSaveHandler) Save(sessionData []byte) error {
	// For now, generate a filename based on timestamp
	// In real implementation, we'd extract sessionId from the JSON
	filename := fmt.Sprintf("session_%d.json", len(sessionData))
	filepath := filepath.Join(h.saveDirectory, filename)

	if err := os.WriteFile(filepath, sessionData, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// Load, List, Delete not needed - UI handles these operations directly

// Browser SaveHandler moved to wasm package
// SaveHandler factory removed - game creators set handlers directly
