package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
	"github.com/1260124186-cc/recipe-planner-service/internal/store"
)

func TestMissingPlanReadDoesNotBlockLaterWrites(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	memory := store.NewMemoryStore()
	if _, err := memory.GetPlan(ctx, "missing"); err == nil {
		t.Fatal("expected missing plan error")
	}

	done := make(chan struct{})
	go func() {
		_ = memory.SavePlan(ctx, testPlan("after-missing-read"))
		close(done)
	}()
	waitForStoreOperation(t, done)
}

func TestMissingPlanCompletionDoesNotBlockLaterReads(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	memory := store.NewMemoryStore()
	if err := memory.CompleteMeal(ctx, "missing", "2026-08-17", nil); err == nil {
		t.Fatal("expected missing plan error")
	}

	done := make(chan struct{})
	go func() {
		_, _ = memory.GetPlan(ctx, "another-missing")
		close(done)
	}()
	waitForStoreOperation(t, done)
}

func testPlan(id string) domain.MealPlan {
	return domain.MealPlan{ID: id, Entries: []domain.MealEntry{{
		Date: "2026-08-17", RecipeID: "soup", Status: domain.MealPlanned,
	}}}
}

func waitForStoreOperation(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(150 * time.Millisecond):
		t.Fatal("a failed plan operation left the memory store blocked")
	}
}
