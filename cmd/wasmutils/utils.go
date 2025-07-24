package wasmutils

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

// WASMFunction represents a function that takes js.Value args and returns (data, error)
type WASMFunction func(args []js.Value) (any, error)

// WASMResponse represents a standardized JavaScript-friendly response
type WASMResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// createWrapper creates a generic wrapper for WASM functions with validation and error handling
func CreateWrapper(minArgs, maxArgs int, fn WASMFunction) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		// Validate argument count
		if len(args) < minArgs {
			return CreateErrorResponse(fmt.Sprintf("Expected at least %d arguments, got %d", minArgs, len(args)))
		}
		if maxArgs >= 0 && len(args) > maxArgs {
			return CreateErrorResponse(fmt.Sprintf("Expected at most %d arguments, got %d", maxArgs, len(args)))
		}

		// Call the function and handle response
		result, err := fn(args)
		if err != nil {
			return CreateErrorResponse(err.Error())
		}

		return CreateSuccessResponse(result)
	})
}

// =============================================================================
// Response Helpers
// =============================================================================

func CreateSuccessResponse(data any) js.Value {
	response := WASMResponse{
		Success: true,
		Message: "Operation completed successfully",
		Data:    data,
	}
	return marshalToJS(response)
}

func CreateErrorResponse(error string) js.Value {
	response := WASMResponse{
		Success: false,
		Error:   error,
	}
	return marshalToJS(response)
}

func marshalToJS(obj any) js.Value {
	bytes, _ := json.Marshal(obj)
	return js.Global().Get("JSON").Call("parse", string(bytes))
}
