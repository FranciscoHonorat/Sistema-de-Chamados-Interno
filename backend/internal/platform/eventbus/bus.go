// Package eventbus is the in-process message bus the modules use to talk to
// each other without importing each other's internals. Messages carry the event
// type and a JSON payload, a transport-independent contract, so a module can be
// extracted into its own service later by swapping the bus for a broker.
package eventbus

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Message struct {
	ID      string
	Type    string
	Payload []byte
}

type Handler func(ctx context.Context, msg Message) error

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]namedHandler
}

type namedHandler struct {
	name    string
	handler Handler
}

func New() *Bus {
	return &Bus{handlers: make(map[string][]namedHandler)}
}

// Subscribe registers handler for every message of eventType. The name only
// shows up in errors, to tell which subscriber failed.
func (b *Bus) Subscribe(eventType, name string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], namedHandler{name: name, handler: handler})
}

// Publish delivers msg synchronously to every subscriber of its type and
// returns the joined errors of the ones that failed. The publisher (the outbox
// relay) keeps the message pending on error and retries it, so delivery is
// at-least-once and subscribers must be idempotent.
func (b *Bus) Publish(ctx context.Context, msg Message) error {
	b.mu.RLock()
	subscribers := append([]namedHandler(nil), b.handlers[msg.Type]...)
	b.mu.RUnlock()

	var errs []error
	for _, s := range subscribers {
		if err := b.deliver(ctx, s, msg); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", s.name, err))
		}
	}
	return errors.Join(errs...)
}

func (b *Bus) deliver(ctx context.Context, s namedHandler, msg Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic handling %s: %v", msg.Type, r)
		}
	}()
	return s.handler(ctx, msg)
}
