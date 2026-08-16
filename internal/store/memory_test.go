package store

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
)

type firstCheckContext struct {
	context.Context
	checked chan<- struct{}
	once    sync.Once
}

func (c *firstCheckContext) Err() error {
	c.once.Do(func() { close(c.checked) })
	return c.Context.Err()
}

func TestSavePlanRejectsCancellationWhileWaitingForLock(t *testing.T) {
	memory := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	memory.mu.Lock()
	checked := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		result <- memory.SavePlan(&firstCheckContext{Context: ctx, checked: checked}, domain.MealPlan{
			ID: "cancelled-while-waiting",
			Entries: []domain.MealEntry{{
				Date: "2026-08-17", RecipeID: "soup", Status: domain.MealPlanned,
			}},
		})
	}()

	select {
	case <-checked:
		cancel()
		memory.mu.Unlock()
	case <-time.After(time.Second):
		memory.mu.Unlock()
		t.Fatal("SavePlan did not check the context before waiting for the lock")
	}

	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("SavePlan error = %v, want context cancellation", err)
		}
	case <-time.After(time.Second):
		t.Fatal("SavePlan did not return after the lock was released")
	}

	if _, err := memory.GetPlan(context.Background(), "cancelled-while-waiting"); err == nil {
		t.Fatal("SavePlan persisted a plan after cancellation while waiting for the lock")
	}
}
