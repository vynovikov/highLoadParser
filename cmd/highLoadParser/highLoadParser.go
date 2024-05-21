// Central point
package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vynovikov/highLoadParser/internal/controllers"
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/infrastructure"
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/internal/routers/tp"
	"github.com/vynovikov/highLoadParser/internal/service"
	"github.com/vynovikov/highLoadParser/internal/transmitters"
)

var (
	wgMain sync.WaitGroup
)

func main() {

	dh := dataHandler.NewMemoryDataHandler()
	repo := repository.NewParserRepository(dh)
	trans := transmitters.NewTransmitter()
	inf := infrastructure.NewInfraStructure(repo, trans)
	srv := service.NewParserService(inf)
	ctr := controllers.NewController(srv)
	router := tp.NewTpReceiver(ctr)

	go router.Run()
	signalListen()
}

// signalListen listens for Interrupt signal, when receiving one invokes stop function
func signalListen() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	<-sigChan

}
