package domain

import (
	"fmt"
	"sort"
	"strings"
)

// IngredientNeed describes the number of portions of an ingredient a recipe needs.
type IngredientNeed struct {
	Name     string `json:"name"`
	Portions int    `json:"portions"`
}

// Recipe is the reusable cooking template selected by a meal plan.
type Recipe struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Tags        []string         `json:"tags"`
	Steps       []string         `json:"steps"`
	Ingredients []IngredientNeed `json:"ingredients"`
}

func (r Recipe) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return fmt.Errorf("recipe id is required")
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("recipe name is required")
	}
	if len(r.Ingredients) == 0 {
		return fmt.Errorf("recipe must include at least one ingredient")
	}

	seen := make(map[string]struct{}, len(r.Ingredients))
	for _, ingredient := range r.Ingredients {
		key := strings.ToLower(strings.TrimSpace(ingredient.Name))
		if key == "" {
			return fmt.Errorf("ingredient name is required")
		}
		if ingredient.Portions <= 0 {
			return fmt.Errorf("ingredient %q portions must be positive", ingredient.Name)
		}
		if _, exists := seen[key]; exists {
			return fmt.Errorf("ingredient %q is listed more than once", ingredient.Name)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func (r Recipe) MatchesTag(tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return true
	}
	for _, candidate := range r.Tags {
		if strings.TrimSpace(candidate) == tag {
			return true
		}
	}
	return false
}

func (r Recipe) Normalized() Recipe {
	normalized := r
	normalized.Name = strings.TrimSpace(r.Name)
	normalized.Tags = append([]string(nil), r.Tags...)
	normalized.Steps = append([]string(nil), r.Steps...)
	normalized.Ingredients = append([]IngredientNeed(nil), r.Ingredients...)
	for i := range normalized.Tags {
		normalized.Tags[i] = strings.ToLower(strings.TrimSpace(normalized.Tags[i]))
	}
	for i := range normalized.Steps {
		normalized.Steps[i] = strings.TrimSpace(normalized.Steps[i])
	}
	sort.Strings(normalized.Tags)
	for i := range normalized.Ingredients {
		normalized.Ingredients[i].Name = strings.ToLower(strings.TrimSpace(normalized.Ingredients[i].Name))
	}
	return normalized
}
