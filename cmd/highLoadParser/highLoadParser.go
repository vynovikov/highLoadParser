// Central point
package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vynovikov/highLoadParser/internal/controllers"
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/internal/routers/tp"
	"github.com/vynovikov/highLoadParser/internal/service"
)

var (
	wgMain sync.WaitGroup
)

func main() {

	dh := dataHandler.NewMemoryDataHandler()
	repo := repository.NewParserRepository(dh)
	srv := service.NewParserService(repo)
	ctr := controllers.NewController(srv)
	router := tp.NewTpReceiver(ctr)

	go router.Run()
	signalListen()
}

// signalListen listens for Interrupt signal, when receiving one invokes stop function
func signalListen() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	logger.L.Println("signal listening")
	sig := <-sigChan

	logger.L.Println("signal detected", sig)

}
