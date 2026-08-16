package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
)

type PlannerService struct {
	recipes RecipeRepository
	plans   PlanRepository
}

func NewPlannerService(recipes RecipeRepository, plans PlanRepository) *PlannerService {
	return &PlannerService{recipes: recipes, plans: plans}
}

func (s *PlannerService) GeneratePlan(ctx context.Context, planID, startDate string, days int, tag string) (domain.MealPlan, error) {
	if days < 1 || days > 14 {
		return domain.MealPlan{}, fmt.Errorf("days must be between 1 and 14")
	}
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return domain.MealPlan{}, fmt.Errorf("start_date must use YYYY-MM-DD")
	}
	recipes, err := s.recipes.ListRecipes(ctx)
	if err != nil {
		return domain.MealPlan{}, err
	}
	matches := make([]domain.Recipe, 0, len(recipes))
	for _, recipe := range recipes {
		if recipe.MatchesTag(tag) {
			matches = append(matches, recipe)
		}
	}
	if len(matches) == 0 {
		return domain.MealPlan{}, fmt.Errorf("no recipe matches tag %q", strings.TrimSpace(tag))
	}

	plan := domain.MealPlan{ID: strings.TrimSpace(planID), Entries: make([]domain.MealEntry, 0, days)}
	for offset := 0; offset < days; offset++ {
		if err := ctx.Err(); err != nil {
			return domain.MealPlan{}, err
		}
		recipe := matches[offset%len(matches)]
		plan.Entries = append(plan.Entries, domain.MealEntry{
			Date:     start.AddDate(0, 0, offset).Format("2006-01-02"),
			RecipeID: recipe.ID,
			Status:   domain.MealPlanned,
		})
	}
	if err := plan.Validate(); err != nil {
		return domain.MealPlan{}, err
	}
	if err := s.plans.SavePlan(ctx, plan); err != nil {
		return domain.MealPlan{}, err
	}
	return plan, nil
}

func (s *PlannerService) GetPlan(ctx context.Context, planID string) (domain.MealPlan, error) {
	return s.plans.GetPlan(ctx, planID)
}

func (s *PlannerService) CookMeal(ctx context.Context, planID, date string) error {
	plan, err := s.plans.GetPlan(ctx, planID)
	if err != nil {
		return err
	}
	for _, entry := range plan.Entries {
		if entry.Date != date {
			continue
		}
		recipe, err := s.recipes.GetRecipe(ctx, entry.RecipeID)
		if err != nil {
			return err
		}
		return s.plans.CompleteMeal(ctx, planID, date, recipe.Ingredients)
	}
	return fmt.Errorf("plan %q has no meal on %s", planID, date)
}
