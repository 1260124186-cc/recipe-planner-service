package domain

import (
	"fmt"
	"strings"
)

// PantryItem records available portions of a normalized ingredient name.
type PantryItem struct {
	Name     string `json:"name"`
	Portions int    `json:"portions"`
}

func NormalizePantryItems(items []PantryItem) ([]PantryItem, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("at least one pantry item is required")
	}
	merged := make(map[string]int, len(items))
	for _, item := range items {
		name := strings.ToLower(strings.TrimSpace(item.Name))
		if name == "" {
			return nil, fmt.Errorf("pantry item name is required")
		}
		if item.Portions <= 0 {
			return nil, fmt.Errorf("pantry item %q portions must be positive", item.Name)
		}
		merged[name] += item.Portions
	}
	normalized := make([]PantryItem, 0, len(merged))
	for name, portions := range merged {
		normalized = append(normalized, PantryItem{Name: name, Portions: portions})
	}
	return normalized, nil
}
