// @title Payment Card Center API
// @version 1.0
// @description API for routing interbank payment transactions
// @termsOfService http://example.com/terms/

// @contact.name Luka Usljebrka
// @contact.email lukauslje13@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
package main

import (
	"fmt"
	"log"

	"github.com/sep-2024-team-35/payment-card-center-service/internal/config"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/handler"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/repository"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/routes"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/service"
)

func main() {
	// Load config
	config.Load("config.yaml")

	// Init dependencies
	repo := repository.NewBankRepository()
	svc := service.NewPCCService(repo)
	h := handler.NewTransactionHandler(svc)

	// Setup router
	router := routes.SetupRoutes(h)

	// Start HTTPS server
	addr := fmt.Sprintf(":%s", config.Global.Server.Port)
	log.Printf("PCC listening on %s …", addr)
	log.Fatal(router.RunTLS(
		addr,
		config.Global.TLS.CertFile,
		config.Global.TLS.KeyFile,
	))
}
