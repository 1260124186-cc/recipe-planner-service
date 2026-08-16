package domain

import (
	"fmt"
	"strings"
	"time"
)

type MealStatus string

const (
	MealPlanned MealStatus = "planned"
	MealCooked  MealStatus = "cooked"
)

type MealEntry struct {
	Date     string     `json:"date"`
	RecipeID string     `json:"recipe_id"`
	Status   MealStatus `json:"status"`
}

type MealPlan struct {
	ID      string      `json:"id"`
	Entries []MealEntry `json:"entries"`
}

func (p MealPlan) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("plan id is required")
	}
	if len(p.Entries) == 0 {
		return fmt.Errorf("meal plan must include at least one entry")
	}
	seenDates := make(map[string]struct{}, len(p.Entries))
	for _, entry := range p.Entries {
		if _, err := time.Parse("2006-01-02", entry.Date); err != nil {
			return fmt.Errorf("invalid meal date %q", entry.Date)
		}
		if strings.TrimSpace(entry.RecipeID) == "" {
			return fmt.Errorf("recipe id is required for %s", entry.Date)
		}
		if entry.Status != MealPlanned && entry.Status != MealCooked {
			return fmt.Errorf("invalid meal status for %s", entry.Date)
		}
		if _, exists := seenDates[entry.Date]; exists {
			return fmt.Errorf("meal date %s appears more than once", entry.Date)
		}
		seenDates[entry.Date] = struct{}{}
	}
	return nil
}
