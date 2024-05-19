// HTTP Reciever.
package tp

import (
	"net"
	"sync"

	"github.com/vynovikov/highLoadParser/internal/controllers"
	"github.com/vynovikov/highLoadParser/internal/logger"
	timeops "github.com/vynovikov/highLoadParser/pkg/timeOps"
)

type TpServer struct {
	l net.Listener
}

type TpReceiver interface {
	Run()
	//HandleRequestFull(net.Conn, string, *sync.WaitGroup)
	//Stop(*sync.WaitGroup)
}

type tpReceiverStruct struct {
	ctrl controllers.ParserController
	srv  *TpServer
	wg   sync.WaitGroup
}

func NewTpReceiver(c controllers.ParserController) *tpReceiverStruct {

	li, err := net.Listen("tcp", ":3000")
	if err != nil {
		logger.L.Error(err)
	}
	logger.L.Info("listening localhost:3000")

	s := &TpServer{l: li}

	return &tpReceiverStruct{
		ctrl: c,
		srv:  s,
	}
}

func (r *tpReceiverStruct) Run() {
	for {
		conn, err := r.srv.l.Accept()

		logger.L.Infof("in tp.Run new request %v, err: %v\n", conn, err)
		if err != nil && conn == nil {

			r.wg.Wait()
			//r.A.ChanInClose()

			return

		}

		r.wg.Add(1)
		ts := timeops.NewTS()
		// consider serial execution -> mutex
		r.ctrl.HandleRequestFull(conn, ts, &r.wg)
	}

}
