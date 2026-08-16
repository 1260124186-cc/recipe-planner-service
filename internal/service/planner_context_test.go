package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
	"github.com/1260124186-cc/recipe-planner-service/internal/service"
	"github.com/1260124186-cc/recipe-planner-service/internal/store"
)

type cancelingRecipeRepository struct {
	cancel context.CancelFunc
}

func (r cancelingRecipeRepository) SaveRecipe(context.Context, domain.Recipe) error { return nil }

func (r cancelingRecipeRepository) ListRecipes(context.Context) ([]domain.Recipe, error) {
	r.cancel()
	return []domain.Recipe{{
		ID:          "soup",
		Name:        "Soup",
		Ingredients: []domain.IngredientNeed{{Name: "beans", Portions: 1}},
	}}, nil
}

func (r cancelingRecipeRepository) GetRecipe(context.Context, string) (domain.Recipe, error) {
	return domain.Recipe{}, errors.New("not used")
}

type contextCheckingPlanRepository struct {
	saved bool
}

func (r *contextCheckingPlanRepository) SavePlan(ctx context.Context, _ domain.MealPlan) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.saved = true
	return nil
}

func (r *contextCheckingPlanRepository) GetPlan(context.Context, string) (domain.MealPlan, error) {
	return domain.MealPlan{}, errors.New("not used")
}

func (r *contextCheckingPlanRepository) CompleteMeal(context.Context, string, string, []domain.IngredientNeed) error {
	return errors.New("not used")
}

func TestGeneratePlanStopsWhenRequestIsCanceledAfterRecipesAreRead(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	repository := &contextCheckingPlanRepository{}
	planner := service.NewPlannerService(cancelingRecipeRepository{cancel: cancel}, repository)

	_, err := planner.GeneratePlan(ctx, "cancelled-plan", "2026-08-17", 2, "")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("GeneratePlan error = %v, want context cancellation", err)
	}
	if repository.saved {
		t.Fatal("GeneratePlan saved a plan after its request was canceled")
	}
}

func TestMemoryStoreRejectsCanceledPlanSave(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	memory := store.NewMemoryStore()
	plan := domain.MealPlan{ID: "cancelled", Entries: []domain.MealEntry{{
		Date: "2026-08-17", RecipeID: "soup", Status: domain.MealPlanned,
	}}}

	err := memory.SavePlan(ctx, plan)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("SavePlan error = %v, want context cancellation", err)
	}
	if _, err := memory.GetPlan(context.Background(), plan.ID); err == nil {
		t.Fatal("canceled SavePlan unexpectedly persisted a plan")
	}
}
