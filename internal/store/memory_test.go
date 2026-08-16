package store_test

import (
	"context"
	"testing"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
	"github.com/1260124186-cc/recipe-planner-service/internal/store"
)

func TestRestockDoesNotApplyValidPrefixWhenBatchContainsInvalidItem(t *testing.T) {
	t.Parallel()
	memory := store.NewMemoryStore()
	ctx := context.Background()

	err := memory.Restock(ctx, []domain.PantryItem{
		{Name: "tomato", Portions: 2},
		{Name: "pasta", Portions: 0},
	})
	if err == nil {
		t.Fatal("expected invalid batch to fail")
	}
	snapshot, err := memory.PantrySnapshot(ctx)
	if err != nil {
		t.Fatalf("pantry snapshot: %v", err)
	}
	if len(snapshot) != 0 {
		t.Fatalf("failed batch partially changed pantry: %#v", snapshot)
	}
}
