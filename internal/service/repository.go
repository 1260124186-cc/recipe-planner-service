package service

import (
	"context"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
)

type RecipeRepository interface {
	SaveRecipe(context.Context, domain.Recipe) error
	ListRecipes(context.Context) ([]domain.Recipe, error)
	GetRecipe(context.Context, string) (domain.Recipe, error)
}

type PantryRepository interface {
	Restock(context.Context, []domain.PantryItem) error
	PantrySnapshot(context.Context) (map[string]int, error)
}

type PlanRepository interface {
	SavePlan(context.Context, domain.MealPlan) error
	GetPlan(context.Context, string) (domain.MealPlan, error)
	CompleteMeal(context.Context, string, string, []domain.IngredientNeed) error
}
