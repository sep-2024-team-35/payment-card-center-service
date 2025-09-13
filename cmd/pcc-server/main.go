package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sep-2024-team-35/payment-card-center-service/internal/config"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/handler"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/repository"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/service"
)

func main() {
	// Učitavanje konfiguracije iz config.yaml
	config.Load("config.yaml")

	// Inicijalizacija repozitorijuma, servisa i handlera
	repo := repository.NewBankRepository()
	svc := service.NewPCCService(repo)
	h := handler.NewTransactionHandler(svc)

	// Postavljanje HTTP mux-a i rute
	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", h.Execute)

	// Start HTTPS server sa TLS certifikatima iz konfiguracije
	addr := fmt.Sprintf(":%s", config.Global.Server.Port)
	log.Printf("PCC listening on %s …", addr)
	log.Fatal(http.ListenAndServeTLS(
		addr,
		config.Global.TLS.CertFile,
		config.Global.TLS.KeyFile,
		mux,
	))
}
