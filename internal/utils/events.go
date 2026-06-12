package utils

// EventHandler is a simple event handler supporting multiple callbacks.
type EventHandler[T any] struct {
	callbacks []func(T)
}

// NewEventHandler creates a new EventHandler.
func NewEventHandler[T any]() *EventHandler[T] {
	return &EventHandler[T]{}
}

// Add registers a callback.
func (h *EventHandler[T]) Add(callback func(T)) {
	h.callbacks = append(h.callbacks, callback)
}

// Remove unregisters a callback.
func (h *EventHandler[T]) Remove(callback func(T)) {
	for i, cb := range h.callbacks {
		if callbackBodyEqual(cb, callback) {
			h.callbacks = append(h.callbacks[:i], h.callbacks[i+1:]...)
			return
		}
	}
}

// Emit invokes all registered callbacks.
func (h *EventHandler[T]) Emit(value T) {
	for _, callback := range h.callbacks {
		callback(value)
	}
}

// Clear removes all callbacks.
func (h *EventHandler[T]) Clear() {
	h.callbacks = nil
}

// HasCallbacks reports whether any callbacks are registered.
func (h *EventHandler[T]) HasCallbacks() bool {
	return len(h.callbacks) > 0
}

func callbackBodyEqual[T any](a, b func(T)) bool {
	return &a == &b
}
