package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1260124186-cc/recipe-planner-service/internal/service"
	"github.com/1260124186-cc/recipe-planner-service/internal/store"
	"github.com/1260124186-cc/recipe-planner-service/internal/transport"
)

func TestHTTPWorkflow(t *testing.T) {
	t.Parallel()
	memory := store.NewMemoryStore()
	server := transport.NewServer(
		service.NewCatalogService(memory, memory),
		service.NewPlannerService(memory, memory),
		service.NewShoppingService(memory, memory, memory),
	)
	handler := server.Routes()

	postJSON(t, handler, "/recipes", `{"id":"soup","name":"Soup","tags":["warm"],"ingredients":[{"name":"beans","portions":2}]}`, http.StatusCreated)
	postJSON(t, handler, "/pantry/restock", `{"items":[{"name":"beans","portions":2}]}`, http.StatusOK)
	response := postJSON(t, handler, "/plans/generate", `{"id":"plan-a","start_date":"2026-08-17","days":1,"tag":"warm"}`, http.StatusCreated)
	var plan struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&plan); err != nil {
		t.Fatalf("decode plan: %v", err)
	}
	if plan.ID != "plan-a" {
		t.Fatalf("plan id = %q, want plan-a", plan.ID)
	}
	postJSON(t, handler, "/plans/plan-a/meals/2026-08-17/cook", "", http.StatusOK)

	request := httptest.NewRequest(http.MethodGet, "/plans/plan-a", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("get plan status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"cooked"`)) {
		t.Fatalf("cooked status missing from %s", recorder.Body.String())
	}
}

func postJSON(t *testing.T, handler http.Handler, path, body string, wantStatus int) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != wantStatus {
		t.Fatalf("POST %s status = %d, want %d; body=%s", path, recorder.Code, wantStatus, recorder.Body.String())
	}
	return recorder
}
