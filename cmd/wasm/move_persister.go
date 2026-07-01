//go:build js && wasm
// +build js,wasm

package main

import (
	"context"
	"fmt"
	"syscall/js"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	pj "google.golang.org/protobuf/encoding/protojson"
)

// jsCallbackPersister bridges Go's singleton.MovePersister interface to a
// JavaScript async function registered via lilbattle.registerMovePersister.
//
// The Save call marshals state + group to protojson (readable in devtools;
// size penalty is trivial for a single group per user action), invokes the
// JS callback with (gameId, stateJSON, groupJSON), and awaits the returned
// Promise. Any Promise rejection surfaces as a Go error so
// BaseGamesService.ProcessMoves can propagate it and the FE can log or
// react.
type jsCallbackPersister struct {
	callback js.Value // JS function: (gameId, stateJSON, groupJSON) -> Promise<void>
}

// Save marshals the arguments and awaits the JS callback's Promise. A
// non-resolvable callback (missing, wrong type, throws synchronously)
// produces a descriptive error; the FE surfaces persistence failures
// via the game-log panel.
func (p *jsCallbackPersister) Save(ctx context.Context, gameID string, state *v1.GameState, group *v1.GameMoveGroup) error {
	if !p.callback.Truthy() {
		return fmt.Errorf("no move persister registered from JS")
	}

	stateJSON, err := pj.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	groupJSON, err := pj.Marshal(group)
	if err != nil {
		return fmt.Errorf("marshal group: %w", err)
	}

	promise := p.callback.Invoke(gameID, string(stateJSON), string(groupJSON))
	return awaitPromise(promise)
}

// awaitPromise blocks until the JS Promise settles and returns a Go error
// if it rejected. Uses Then/Catch to bridge the JS event loop back to the
// Go goroutine via a done channel — this is the standard syscall/js pattern
// for turning a Promise into synchronous Go control flow.
func awaitPromise(promise js.Value) error {
	if !promise.Truthy() {
		return nil
	}
	if then := promise.Get("then"); !then.Truthy() {
		return nil
	}

	done := make(chan error, 1)

	onResolve := js.FuncOf(func(this js.Value, args []js.Value) any {
		done <- nil
		return nil
	})
	defer onResolve.Release()

	onReject := js.FuncOf(func(this js.Value, args []js.Value) any {
		msg := "js persister rejected"
		if len(args) > 0 {
			if s := args[0].Get("message"); s.Truthy() {
				msg = s.String()
			} else {
				msg = args[0].String()
			}
		}
		done <- fmt.Errorf("%s", msg)
		return nil
	})
	defer onReject.Release()

	promise.Call("then", onResolve).Call("catch", onReject)
	return <-done
}
