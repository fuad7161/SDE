package main

import (
	"log"
	httpHand "main/internal/adapters/http"
	"main/internal/adapters/postgres"
	"main/internal/usecases"
	"net/http"
)

func main() {

	// 1. Infrastructure (adapters)
	dbRepo := &postgres.Repo{}

	// 2. Core (use cases depend on output ports)
	userUC := usecases.NewUserService(dbRepo)

	// 3. Transport (adapters depend on input ports)
	handler := httpHand.NewHandler(userUC)

	// 4. Run
	mux := http.NewServeMux()
	mux.HandleFunc("/users", handler.GetUser)
	log.Println("Server running on: 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
