package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {

	gateway_server := NewAdkgGatewayQueryServer(":8880")
	go func() {
		defer log.Info().Msg("gateway server exit")
		log.Info().Msg("starting the gateway server....")
		if err := gateway_server.Start(); err != nil {
			log.Error().Err(err).Msg("fail to start the gateway server")
		}
	}()

	// wait....
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("stop signal received")

	// stop gateway
	if err := gateway_server.Stop(); err != nil {
		log.Fatal().Err(err).Msg("Failed to stop gateway query server")
	}

}
