package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type (
	AdkgGatewayServer struct {
		s      *http.Server
		logger zerolog.Logger
	}
)

func NewAdkgGatewayQueryServer(addr string) *AdkgGatewayServer {
	ns := &AdkgGatewayServer{
		logger: log.With().Str("module", "adkg_gateway_query_server").Logger(),
	}
	s := &http.Server{
		Addr:    addr,
		Handler: ns.newHandler(),
	}
	ns.s = s
	return ns
}

func (ns *AdkgGatewayServer) newHandler() http.Handler {
	router := mux.NewRouter()
	router.Handle("/api/v1/ping", http.HandlerFunc(ns.pingHandler)).Methods(http.MethodGet)
	// router.Handle("/bc_height", http.HandlerFunc(ns.GetBlockChainHeightHandler)).Methods(http.MethodGet)
	router.Handle("/api/v1/get-app-config/", http.HandlerFunc(ns.GetAppConfigHandler)).Methods(http.MethodGet)
	return router
}

func (ns *AdkgGatewayServer) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type Creds struct {
	Provider string `json:"verifier"`
	ClientID string `json:"client_id"`
	Domain   string `json:"domain"`
}

type appConfigResp struct {
	Creds  []Creds `json:"cred"`
	Global bool    `json:"global"`
}

func (ns *AdkgGatewayServer) GetAppConfigHandler(w http.ResponseWriter, req *http.Request) {
	app_id := req.URL.Query().Get("id")

	if app_id == "" {
		http.Error(w, "app id is empty", http.StatusBadRequest)
		return
	}

	_creds := []Creds{
		{
			Provider: "google",
			ClientID: "938953074000-d96t0tkvjfkf4ag9i04f8man85nsu6kb.apps.googleusercontent.com",
			Domain:   "",
		},
	}
	body := appConfigResp{Global: false, Creds: _creds}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
	return
}

func (ns *AdkgGatewayServer) Start() error {
	if ns.s == nil {
		return errors.New("invalid narada query http server instance")
	}

	if err := ns.s.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return fmt.Errorf("fail to start narada query http server: %w", err)
		}
	}
	return nil
}

func (ns *AdkgGatewayServer) Stop() error {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := ns.s.Shutdown(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shutdown the narada query server gracefully")
	}
	return err
}
