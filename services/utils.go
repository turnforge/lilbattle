package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// newRandomId generates a new unique random ID of specified length (default 8 chars)
// It is upto the caller to check for collissions
func newRandomId(numChars ...int) (string, error) {
	const maxRetries = 10

	// Default to 8 characters if not specified
	length := 8
	if len(numChars) > 0 && numChars[0] > 0 {
		length = numChars[0]
	}

	// Calculate number of bytes needed (2 hex chars per byte)
	numBytes := (length + 1) / 2

	bytes := make([]byte, numBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func newProtoInstance[T proto.Message]() (out T) {
	var zero T
	tType := reflect.TypeOf(zero)

	// If T is a pointer type, create new instance
	if tType.Kind() == reflect.Ptr {
		out = reflect.New(tType.Elem()).Interface().(T)
	} else {
		panic("only pointer types supported")
	}
	return
}
