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
	log.Println("🚀 [BOOT] Starting Payment Card Center API...")
	
	config.Load("config.yaml")
	log.Printf("⚙️  [CONFIG] Loaded successfully. Server will run on port %s", config.Global.Server.Port)

	repo := repository.NewBankRepository()
	log.Println("💾 [INIT] BankRepository initialized")

	svc := service.NewPCCService(repo)
	log.Println("🛠️  [INIT] PCCService initialized")

	h := handler.NewTransactionHandler(svc)
	log.Println("📦 [INIT] TransactionHandler initialized")

	router := routes.SetupRoutes(h)
	log.Println("🛣️  [ROUTER] Routes configured successfully")

	addr := fmt.Sprintf(":%s", config.Global.Server.Port)
	log.Printf("🌐 [SERVER] Listening on %s …", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("❌ [FATAL] Could not start server: %v", err)
	}
}
