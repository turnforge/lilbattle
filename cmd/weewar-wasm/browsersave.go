//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
)

// BrowserSaveHandler calls JavaScript functions to save via browser APIs
type BrowserSaveHandler struct {
	apiEndpoint string
}

// NewBrowserSaveHandler creates a new browser-based save handler
func NewBrowserSaveHandler(apiEndpoint string) *BrowserSaveHandler {
	return &BrowserSaveHandler{
		apiEndpoint: apiEndpoint,
	}
}

func (h *BrowserSaveHandler) Save(sessionData []byte) error {
	// Check if we're running in a browser environment
	if js.Global().IsUndefined() {
		return fmt.Errorf("browser save handler requires browser environment")
	}
	
	// Call JavaScript function to handle the save
	saveHandler := js.Global().Get("gameSaveHandler")
	if saveHandler.IsUndefined() {
		return fmt.Errorf("gameSaveHandler JavaScript function not found")
	}
	
	// Convert to string for JS
	sessionDataStr := string(sessionData)
	
	// Call the JS function and wait for result
	result := saveHandler.Invoke(sessionDataStr)
	
	// Check if the save was successful
	success := result.Get("success").Bool()
	if !success {
		errorMsg := result.Get("error").String()
		return fmt.Errorf("browser save failed: %s", errorMsg)
	}
	
	return nil
}