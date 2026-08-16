package service

import (
	"context"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
)

type CatalogService struct {
	recipes RecipeRepository
	pantry  PantryRepository
}

func NewCatalogService(recipes RecipeRepository, pantry PantryRepository) *CatalogService {
	return &CatalogService{recipes: recipes, pantry: pantry}
}

func (s *CatalogService) CreateRecipe(ctx context.Context, recipe domain.Recipe) error {
	recipe = recipe.Normalized()
	if err := recipe.Validate(); err != nil {
		return err
	}
	return s.recipes.SaveRecipe(ctx, recipe)
}

func (s *CatalogService) ListRecipes(ctx context.Context) ([]domain.Recipe, error) {
	return s.recipes.ListRecipes(ctx)
}

func (s *CatalogService) RestockPantry(ctx context.Context, items []domain.PantryItem) error {
	return s.pantry.Restock(ctx, items)
}
