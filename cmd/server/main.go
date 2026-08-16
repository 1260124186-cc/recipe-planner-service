package main

import (
	"log"
	"net/http"
	"os"

	"github.com/1260124186-cc/recipe-planner-service/internal/service"
	"github.com/1260124186-cc/recipe-planner-service/internal/store"
	"github.com/1260124186-cc/recipe-planner-service/internal/transport"
)

func main() {
	address := os.Getenv("RECIPE_PLANNER_ADDR")
	if address == "" {
		address = ":8080"
	}

	memoryStore := store.NewMemoryStore()
	catalog := service.NewCatalogService(memoryStore, memoryStore)
	planner := service.NewPlannerService(memoryStore, memoryStore)
	shopping := service.NewShoppingService(memoryStore, memoryStore, memoryStore)
	server := transport.NewServer(catalog, planner, shopping)

	log.Printf("recipe planner listening on %s", address)
	if err := http.ListenAndServe(address, server.Routes()); err != nil {
		log.Fatal(err)
	}
}
