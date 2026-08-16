package service_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
	"github.com/1260124186-cc/recipe-planner-service/internal/service"
	"github.com/1260124186-cc/recipe-planner-service/internal/store"
)

func TestPlanCookingAndShoppingWorkflow(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	memory := store.NewMemoryStore()
	catalog := service.NewCatalogService(memory, memory)
	planner := service.NewPlannerService(memory, memory)
	shopping := service.NewShoppingService(memory, memory, memory)

	recipe := domain.Recipe{
		ID:   "tomato-pasta",
		Name: "Tomato Pasta",
		Tags: []string{"Quick"},
		Ingredients: []domain.IngredientNeed{
			{Name: "Tomato", Portions: 2},
			{Name: "Pasta", Portions: 1},
		},
	}
	if err := catalog.CreateRecipe(ctx, recipe); err != nil {
		t.Fatalf("create recipe: %v", err)
	}
	if err := catalog.RestockPantry(ctx, []domain.PantryItem{{Name: "tomato", Portions: 2}, {Name: "pasta", Portions: 1}}); err != nil {
		t.Fatalf("restock pantry: %v", err)
	}

	plan, err := planner.GeneratePlan(ctx, "week-1", "2026-08-17", 2, "quick")
	if err != nil {
		t.Fatalf("generate plan: %v", err)
	}
	if len(plan.Entries) != 2 || plan.Entries[0].RecipeID != recipe.ID {
		t.Fatalf("unexpected generated plan: %+v", plan)
	}

	items, err := shopping.BuildList(ctx, plan.ID)
	if err != nil {
		t.Fatalf("build shopping list: %v", err)
	}
	wantBeforeCooking := []service.ShoppingItem{{Name: "pasta", Portions: 1}, {Name: "tomato", Portions: 2}}
	if !reflect.DeepEqual(items, wantBeforeCooking) {
		t.Fatalf("shopping list before cooking = %#v, want %#v", items, wantBeforeCooking)
	}

	if err := planner.CookMeal(ctx, plan.ID, "2026-08-17"); err != nil {
		t.Fatalf("cook meal: %v", err)
	}
	updated, err := planner.GetPlan(ctx, plan.ID)
	if err != nil {
		t.Fatalf("get plan: %v", err)
	}
	if updated.Entries[0].Status != domain.MealCooked {
		t.Fatalf("first meal status = %q, want cooked", updated.Entries[0].Status)
	}
}

func TestRestockRejectsInvalidBatchWithoutChangingPantry(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	memory := store.NewMemoryStore()
	catalog := service.NewCatalogService(memory, memory)

	err := catalog.RestockPantry(ctx, []domain.PantryItem{
		{Name: "tomato", Portions: 2},
		{Name: "pasta", Portions: 0},
	})
	if err == nil {
		t.Fatal("expected invalid restock batch to fail")
	}
	snapshot, err := memory.PantrySnapshot(ctx)
	if err != nil {
		t.Fatalf("pantry snapshot: %v", err)
	}
	if len(snapshot) != 0 {
		t.Fatalf("invalid restock changed pantry: %#v", snapshot)
	}
}
