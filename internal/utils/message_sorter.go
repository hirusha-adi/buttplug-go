package utils

import (
	"context"
	"sync"
	"time"

	"github.com/hirusha-adi/buttplug-go/internal/messages"
)

const maxMessageID = 4294967295

type waitResult struct {
	msg messages.Message
	err error
}

// MessageSorter correlates outgoing requests with incoming responses by message Id.
type MessageSorter struct {
	mu      sync.Mutex
	nextID  uint32
	pending map[int]chan waitResult
}

// NewMessageSorter creates a new MessageSorter.
func NewMessageSorter() *MessageSorter {
	return &MessageSorter{
		nextID:  1,
		pending: make(map[int]chan waitResult),
	}
}

// GetNextID returns the next available message ID.
func (s *MessageSorter) GetNextID() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	msgID := int(s.nextID)
	if s.nextID >= maxMessageID {
		s.nextID = 1
	} else {
		s.nextID++
	}
	return msgID
}

// WaitForResponse waits for a response with the matching message ID.
func (s *MessageSorter) WaitForResponse(ctx context.Context, msgID int, timeout time.Duration) (messages.Message, error) {
	ch := make(chan waitResult, 1)

	s.mu.Lock()
	s.pending[msgID] = ch
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.pending, msgID)
		s.mu.Unlock()
	}()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case result := <-ch:
		return result.msg, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Resolve resolves a pending request with a response.
func (s *MessageSorter) Resolve(msgID int, response messages.Message) bool {
	s.mu.Lock()
	ch, ok := s.pending[msgID]
	s.mu.Unlock()
	if !ok {
		return false
	}
	select {
	case ch <- waitResult{msg: response}:
		return true
	default:
		return false
	}
}

// RejectAll rejects all pending requests with an error.
func (s *MessageSorter) RejectAll(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for msgID, ch := range s.pending {
		select {
		case ch <- waitResult{err: err}:
		default:
		}
		delete(s.pending, msgID)
	}
}

// PendingCount returns the number of pending requests.
func (s *MessageSorter) PendingCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.pending)
}

// SetNextIDForTest sets the next message ID for tests.
func SetNextIDForTest(s *MessageSorter, nextID uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID = nextID
}
