package service_test

import (
	"context"
	"testing"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
	"github.com/1260124186-cc/recipe-planner-service/internal/service"
	"github.com/1260124186-cc/recipe-planner-service/internal/store"
)

func TestGeneratePlanOnlyUsesRecipesMatchingRequestedTag(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	memory := store.NewMemoryStore()
	catalog := service.NewCatalogService(memory, memory)
	planner := service.NewPlannerService(memory, memory)

	for _, recipe := range []domain.Recipe{
		{ID: "quick-pasta", Name: "Quick Pasta", Tags: []string{"quick"}, Ingredients: []domain.IngredientNeed{{Name: "pasta", Portions: 1}}},
		{ID: "warm-soup", Name: "Warm Soup", Tags: []string{"warm"}, Ingredients: []domain.IngredientNeed{{Name: "beans", Portions: 1}}},
	} {
		if err := catalog.CreateRecipe(ctx, recipe); err != nil {
			t.Fatalf("create recipe %s: %v", recipe.ID, err)
		}
	}

	plan, err := planner.GeneratePlan(ctx, "warm-week", "2026-08-17", 2, "warm")
	if err != nil {
		t.Fatalf("generate tagged plan: %v", err)
	}
	for _, entry := range plan.Entries {
		if entry.RecipeID != "warm-soup" {
			t.Fatalf("tagged plan selected %q, want only warm-soup", entry.RecipeID)
		}
	}
}

func TestRecipeTagMatchingIgnoresCallerCase(t *testing.T) {
	t.Parallel()
	recipe := domain.Recipe{Tags: []string{"Quick"}}
	if !recipe.MatchesTag("QUICK") {
		t.Fatal("tag matching should not depend on capitalization")
	}
}
