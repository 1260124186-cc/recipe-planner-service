package store

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
)

// MemoryStore is a concurrency-safe local store intended for a single service process.
type MemoryStore struct {
	mu      sync.RWMutex
	recipes map[string]domain.Recipe
	pantry  map[string]int
	plans   map[string]domain.MealPlan
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		recipes: make(map[string]domain.Recipe),
		pantry:  make(map[string]int),
		plans:   make(map[string]domain.MealPlan),
	}
}

func (s *MemoryStore) SaveRecipe(ctx context.Context, recipe domain.Recipe) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.recipes[recipe.ID]; exists {
		return fmt.Errorf("recipe %q already exists", recipe.ID)
	}
	s.recipes[recipe.ID] = cloneRecipe(recipe)
	return nil
}

func (s *MemoryStore) ListRecipes(ctx context.Context) ([]domain.Recipe, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	recipes := make([]domain.Recipe, 0, len(s.recipes))
	for _, recipe := range s.recipes {
		recipes = append(recipes, cloneRecipe(recipe))
	}
	sort.Slice(recipes, func(i, j int) bool { return recipes[i].ID < recipes[j].ID })
	return recipes, nil
}

func (s *MemoryStore) GetRecipe(ctx context.Context, id string) (domain.Recipe, error) {
	if err := ctx.Err(); err != nil {
		return domain.Recipe{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	recipe, exists := s.recipes[id]
	if !exists {
		return domain.Recipe{}, fmt.Errorf("recipe %q was not found", id)
	}
	return cloneRecipe(recipe), nil
}

func (s *MemoryStore) Restock(ctx context.Context, items []domain.PantryItem) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		s.pantry[item.Name] += item.Portions
	}
	return nil
}

func (s *MemoryStore) PantrySnapshot(ctx context.Context) (map[string]int, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot := make(map[string]int, len(s.pantry))
	for name, portions := range s.pantry {
		snapshot[name] = portions
	}
	return snapshot, nil
}

func (s *MemoryStore) SavePlan(ctx context.Context, plan domain.MealPlan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.plans[plan.ID]; exists {
		return fmt.Errorf("plan %q already exists", plan.ID)
	}
	s.plans[plan.ID] = clonePlan(plan)
	return nil
}

func (s *MemoryStore) GetPlan(ctx context.Context, id string) (domain.MealPlan, error) {
	if err := ctx.Err(); err != nil {
		return domain.MealPlan{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	plan, exists := s.plans[id]
	if !exists {
		return domain.MealPlan{}, fmt.Errorf("plan %q was not found", id)
	}
	return clonePlan(plan), nil
}

// CompleteMeal consumes ingredients and marks an entry cooked under one lock.
func (s *MemoryStore) CompleteMeal(ctx context.Context, planID, date string, needs []domain.IngredientNeed) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	plan, exists := s.plans[planID]
	if !exists {
		return fmt.Errorf("plan %q was not found", planID)
	}
	entryIndex := -1
	for i := range plan.Entries {
		if plan.Entries[i].Date == date {
			entryIndex = i
			break
		}
	}
	if entryIndex < 0 {
		return fmt.Errorf("plan %q has no meal on %s", planID, date)
	}
	if plan.Entries[entryIndex].Status == domain.MealCooked {
		return fmt.Errorf("meal on %s is already cooked", date)
	}
	for _, need := range needs {
		if s.pantry[need.Name] < need.Portions {
			return fmt.Errorf("not enough %s to cook the meal", need.Name)
		}
	}
	for _, need := range needs {
		s.pantry[need.Name] -= need.Portions
	}
	plan.Entries[entryIndex].Status = domain.MealCooked
	s.plans[planID] = plan
	return nil
}

func cloneRecipe(recipe domain.Recipe) domain.Recipe {
	recipe.Tags = append([]string(nil), recipe.Tags...)
	recipe.Steps = append([]string(nil), recipe.Steps...)
	recipe.Ingredients = append([]domain.IngredientNeed(nil), recipe.Ingredients...)
	return recipe
}

func clonePlan(plan domain.MealPlan) domain.MealPlan {
	plan.Entries = append([]domain.MealEntry(nil), plan.Entries...)
	return plan
}
