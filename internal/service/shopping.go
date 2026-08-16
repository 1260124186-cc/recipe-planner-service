package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
)

type ShoppingItem struct {
	Name     string `json:"name"`
	Portions int    `json:"portions"`
}

type ShoppingService struct {
	recipes RecipeRepository
	pantry  PantryRepository
	plans   PlanRepository
}

func NewShoppingService(recipes RecipeRepository, pantry PantryRepository, plans PlanRepository) *ShoppingService {
	return &ShoppingService{recipes: recipes, pantry: pantry, plans: plans}
}

func (s *ShoppingService) BuildList(ctx context.Context, planID string) ([]ShoppingItem, error) {
	plan, err := s.plans.GetPlan(ctx, planID)
	if err != nil {
		return nil, err
	}
	required := make(map[string]int)
	for _, entry := range plan.Entries {
		if entry.Status == domain.MealCooked {
			continue
		}
		recipe, err := s.recipes.GetRecipe(ctx, entry.RecipeID)
		if err != nil {
			return nil, err
		}
		for _, ingredient := range recipe.Ingredients {
			required[ingredient.Name] += ingredient.Portions
		}
	}
	pantry, err := s.pantry.PantrySnapshot(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ShoppingItem, 0, len(required))
	for name, needed := range required {
		if shortage := needed - pantry[name]; shortage > 0 {
			items = append(items, ShoppingItem{Name: name, Portions: shortage})
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func (s *ShoppingService) ValidatePlanExists(ctx context.Context, planID string) error {
	if _, err := s.plans.GetPlan(ctx, planID); err != nil {
		return fmt.Errorf("shopping list unavailable: %w", err)
	}
	return nil
}
