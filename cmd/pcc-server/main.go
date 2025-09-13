package main

import (
	"fmt"
	"log"
	"net/http"
	"payment-card-center-service"
)

func main() {
	confi.Load("config.yaml")

	repo := repository.NewBankRepository()
	svc := service.NewPCCService(repo)
	h := handler.NewTransactionHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", h.Execute)

	addr := fmt.Sprintf(":%d", config.Global.Server.Port)
	log.Printf("PCC listening on %s …", addr)
	log.Fatal(http.ListenAndServeTLS(addr,
		config.Global.TLS.CertFile,
		config.Global.TLS.KeyFile,
		mux))
}
