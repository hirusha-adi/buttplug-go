package buttplug_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hirusha-adi/buttplug-go/internal/messages"
	"github.com/hirusha-adi/buttplug-go/internal/utils"
)

func TestMessageSorterGetNextIDIncrements(t *testing.T) {
	sorter := utils.NewMessageSorter()
	if sorter.GetNextID() != 1 || sorter.GetNextID() != 2 || sorter.GetNextID() != 3 {
		t.Fatal("unexpected ids")
	}
}

func TestMessageSorterWaitAndResolve(t *testing.T) {
	sorter := utils.NewMessageSorter()
	msgID := sorter.GetNextID()
	response := &messages.Ok{BaseMessage: messages.BaseMessage{ID: msgID}}

	resultCh := make(chan messages.Message, 1)
	errCh := make(chan error, 1)
	go func() {
		msg, err := sorter.WaitForResponse(context.Background(), msgID, time.Second)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- msg
	}()

	time.Sleep(10 * time.Millisecond)
	if !sorter.Resolve(msgID, response) {
		t.Fatal("expected resolve to succeed")
	}

	select {
	case result := <-resultCh:
		if ok, _ := result.(*messages.Ok); ok == nil || result.GetID() != msgID {
			t.Fatalf("unexpected result: %+v", result)
		}
	case err := <-errCh:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestMessageSorterResolveUnknownID(t *testing.T) {
	sorter := utils.NewMessageSorter()
	if sorter.Resolve(999, &messages.Ok{BaseMessage: messages.BaseMessage{ID: 999}}) {
		t.Fatal("expected false")
	}
}

func TestMessageSorterTimeout(t *testing.T) {
	sorter := utils.NewMessageSorter()
	msgID := sorter.GetNextID()
	_, err := sorter.WaitForResponse(context.Background(), msgID, 10*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}

func TestMessageSorterRejectAll(t *testing.T) {
	sorter := utils.NewMessageSorter()
	msgID1 := sorter.GetNextID()
	msgID2 := sorter.GetNextID()

	errCh1 := make(chan error, 1)
	errCh2 := make(chan error, 1)
	go func() {
		_, err := sorter.WaitForResponse(context.Background(), msgID1, time.Second)
		errCh1 <- err
	}()
	go func() {
		_, err := sorter.WaitForResponse(context.Background(), msgID2, time.Second)
		errCh2 <- err
	}()

	time.Sleep(10 * time.Millisecond)
	testErr := errors.New("Test error")
	sorter.RejectAll(testErr)

	if err := <-errCh1; err == nil || err.Error() != testErr.Error() {
		t.Fatalf("unexpected err1: %v", err)
	}
	if err := <-errCh2; err == nil || err.Error() != testErr.Error() {
		t.Fatalf("unexpected err2: %v", err)
	}
}

func TestMessageSorterPendingCount(t *testing.T) {
	sorter := utils.NewMessageSorter()
	if sorter.PendingCount() != 0 {
		t.Fatal("expected 0 pending")
	}

	msgID := sorter.GetNextID()
	done := make(chan struct{})
	go func() {
		_, _ = sorter.WaitForResponse(context.Background(), msgID, time.Second)
		close(done)
	}()

	time.Sleep(10 * time.Millisecond)
	if sorter.PendingCount() != 1 {
		t.Fatalf("expected 1 pending, got %d", sorter.PendingCount())
	}

	sorter.Resolve(msgID, &messages.Ok{BaseMessage: messages.BaseMessage{ID: msgID}})
	<-done

	if sorter.PendingCount() != 0 {
		t.Fatal("expected 0 pending after resolve")
	}
}

func TestMessageSorterIDWrapsAtMax(t *testing.T) {
	sorter := utils.NewMessageSorter()
	utils.SetNextIDForTest(sorter, 4294967295)
	id1 := sorter.GetNextID()
	id2 := sorter.GetNextID()
	if id1 != 4294967295 || id2 != 1 {
		t.Fatalf("expected wrap, got %d and %d", id1, id2)
	}
}
