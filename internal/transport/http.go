package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/1260124186-cc/recipe-planner-service/internal/domain"
	"github.com/1260124186-cc/recipe-planner-service/internal/service"
)

type Server struct {
	catalog  *service.CatalogService
	planner  *service.PlannerService
	shopping *service.ShoppingService
}

func NewServer(catalog *service.CatalogService, planner *service.PlannerService, shopping *service.ShoppingService) *Server {
	return &Server{catalog: catalog, planner: planner, shopping: shopping}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /recipes", s.createRecipe)
	mux.HandleFunc("GET /recipes", s.listRecipes)
	mux.HandleFunc("POST /pantry/restock", s.restockPantry)
	mux.HandleFunc("POST /plans/generate", s.generatePlan)
	mux.HandleFunc("GET /plans/{id}", s.getPlan)
	mux.HandleFunc("POST /plans/{id}/meals/{date}/cook", s.cookMeal)
	mux.HandleFunc("GET /shopping", s.shoppingList)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) createRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe domain.Recipe
	if err := decodeJSON(r, &recipe); err != nil {
		writeError(w, err)
		return
	}
	if strings.TrimSpace(recipe.ID) == "" {
		recipe.ID = service.NewID()
	}
	if err := s.catalog.CreateRecipe(r.Context(), recipe); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, recipe.Normalized())
}

func (s *Server) listRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := s.catalog.ListRecipes(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, recipes)
}

func (s *Server) restockPantry(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Items []domain.PantryItem `json:"items"`
	}
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}
	if err := s.catalog.RestockPantry(r.Context(), request.Items); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "restocked"})
}

func (s *Server) generatePlan(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID        string `json:"id"`
		StartDate string `json:"start_date"`
		Days      int    `json:"days"`
		Tag       string `json:"tag"`
	}
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}
	if strings.TrimSpace(request.ID) == "" {
		request.ID = service.NewID()
	}
	plan, err := s.planner.GeneratePlan(r.Context(), request.ID, request.StartDate, request.Days, request.Tag)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, plan)
}

func (s *Server) getPlan(w http.ResponseWriter, r *http.Request) {
	plan, err := s.planner.GetPlan(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, plan)
}

func (s *Server) cookMeal(w http.ResponseWriter, r *http.Request) {
	if err := s.planner.CookMeal(r.Context(), r.PathValue("id"), r.PathValue("date")); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cooked"})
}

func (s *Server) shoppingList(w http.ResponseWriter, r *http.Request) {
	planID := strings.TrimSpace(r.URL.Query().Get("plan_id"))
	if planID == "" {
		writeError(w, fmt.Errorf("plan_id is required"))
		return
	}
	items, err := s.shopping.BuildList(r.Context(), planID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("invalid JSON request: %w", err)
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		status = http.StatusRequestTimeout
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
