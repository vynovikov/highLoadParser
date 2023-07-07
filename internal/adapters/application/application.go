// Application layer. Central place through witch all adapters interact
package application

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/vynovikov/highLoadParser/internal/adapters/driven/rpc"
	"github.com/vynovikov/highLoadParser/internal/adapters/driven/store"
	"github.com/vynovikov/highLoadParser/internal/repo"

	"github.com/google/go-cmp/cmp"
)

type AppService struct {
	stopping        bool
	chanInClosed    bool
	transmitterLock sync.Mutex
	appRWLock       sync.RWMutex

	W repo.WaitGroups
	C repo.Channels
}

type PieceKey struct {
	TS   string
	Part int
}

// Spy logger interface for testing
type Logger interface {
	LogStuff(repo.AppUnit)
}

// All adapters combined
type App struct {
	T rpc.Transmitter
	A AppService
	S store.Store
	L Logger // Logger interface is testDouble spy for testing Handle method
}

func NewAppFull(s store.Store, t rpc.Transmitter) (*App, chan struct{}) {
	done := make(chan struct{})

	App := &App{
		T: t,
		A: NewAppService(done),
		S: s,
	}
	return App, done
}

func NewAppEmpty() *App {

	App := &App{
		A: NewAppService(make(chan struct{})),
	}
	return App
}

func NewAppService(done chan struct{}) AppService {
	return AppService{
		W: repo.WaitGroups{
			M: make(map[repo.AppStoreKeyGeneral]*sync.WaitGroup),
		},
		C: repo.Channels{
			ChanIn:  make(chan repo.AppFeederUnit, 10),
			ChanOut: make(chan repo.AppDistributorUnit, 10),
			ChanLog: make(chan string, 10),
			Done:    done,
		},
	}
}

func (a *App) MountLogger(l Logger) {
	a.L = l
}

type Application interface {
	Start()
	AddToFeeder(repo.ReceiverUnit)
	Stop()
	ChanInClose()
	SetStopping()
	Stopping() bool
}

func (a *App) Start() {

	a.A.W.Workers.Add(1)
	go a.Work(1)

	a.A.W.Sender.Add(1)
	go a.Send()

	go a.Log()

}

func (a *App) toChanOut(adu repo.AppDistributorUnit) {

	if a.L != nil {
		a.L.LogStuff(adu)
	}
	a.A.C.ChanOut <- adu

}

func (a *App) toChanLog(s string) {
	a.A.C.ChanLog <- s
}

// Work is the central function of whole application.
// Handles data from receiver, sends results to transmitter
func (a *App) Work(i int) {

	for afu := range a.A.C.ChanIn {
		if len(afu.R.B.B) == 0 {
			continue
		}
		// Reading feederUnit bytes chunk, finding to boundary appearance and slicing it into dataPieces
		dataPieces := repo.Slicer(afu)

		for _, v := range dataPieces {
			a.Handle(v, afu.R.H.Bou)
		}
	}
	close(a.A.C.ChanLog)
	close(a.A.C.ChanOut)

}

// AddToFeeder updates receiver data and sends it to chanIn
func (a *App) AddToFeeder(in repo.ReceiverUnit) {
	if a.L != nil {
		a.L.LogStuff(in)
	}

	A := repo.NewAppFeederUnit(in)

	askg := repo.NewAppStoreKeyGeneralFromFeeder(A)

	if _, ok := a.A.W.M[askg]; !ok {
		a.A.W.M[askg] = &sync.WaitGroup{}
	}

	a.A.C.ChanIn <- A
}

// Send is running as gourutine. Initiates transmission for any data got from chanOut
func (a *App) Send() {

	for adu := range a.A.C.ChanOut {
		a.T.Transmit(adu)
	}
	a.A.W.Sender.Done()
	close(a.A.C.Done)

}

func (a *App) Log() {
	for l := range a.A.C.ChanLog {
		a.T.Log(l)
	}
}

func NewPieceKeyFromAPU(apu repo.AppPieceUnit) PieceKey {
	return PieceKey{
		Part: apu.APH.Part,
		TS:   apu.APH.TS,
	}
}

// Handles dataPieces depending on its parameters and state of store.
// Tested in application_test.go
func (a *App) Handle(d repo.DataPiece, bou repo.Boundary) {
	prepErrs := make([]error, 0)
	//logger.L.Infof("application.Handle got dataPiece header %v, body %q\n", d.GetHeader(), d.GetBody(0))
	a.toChanLog(fmt.Sprintf("in highLoadparser application.Handle was invoked for dataPiece with header %v, body %q", d.GetHeader(), d.GetBody(0)))

	if d.B() == repo.False && (d.E() == repo.False || d.E() == repo.Last) { // siglePart data

		adu := repo.NewAppDistributorUnitKafka(d, bou)
		a.toChanOut(adu)
		return
	}

	// multiPart data

	adu := repo.AppDistributorUnit{}
	adub, header, bErr := CalcBody(d, bou)
	//logger.L.Infof("application.Handle got dataPiece header %v, body %q\n", d.GetHeader(), d.GetBody(0))

	presence, err := a.S.Presence(d)
	if err != nil {
		prepErrs = append(prepErrs, err)
		return
	}
	sc, scErr := repo.NewStoreChange(d, presence, bou)
	if scErr != nil && bErr != nil && scErr.Error() != bErr.Error() && !strings.Contains(scErr.Error(), "changed") {
		prepErrs = append(prepErrs, scErr)
		return
	}
	a.S.Act(d, sc)
	if bErr != nil {
		prepErrs = append(prepErrs, bErr)
		if strings.Contains(bErr.Error(), "is not full") {
			return
		}
		if strings.Contains(bErr.Error(), "is ending part") {
			if scErr != nil && scErr.Error() != bErr.Error() && !strings.Contains(scErr.Error(), "changed") {
				prepErrs = append(prepErrs, scErr)
				return
			}

			if scErr != nil && strings.Contains(scErr.Error(), "changed") {
				adu = repo.NewAppDistributorUnitKafkaPrepared(d, header[bytes.Index(header, []byte("Content-Disposition")):], bErr, adub)

			} else {
				aduh := repo.NewAppDistributorHeaderKafkaFromSC(d, sc)
				adu = repo.NewAppDistributorUnit(aduh, adub)
			}

			if !cmp.Equal(adu, repo.AppDistributorUnit{}) {
				a.toChanOut(adu)
				if d.E() == repo.Last {
					a.S.Reset(repo.NewAppStoreKeyGeneralFromDataPiece(d))
				}
				return
			}
		}
	}

	if d.E() == repo.Probably {
		if d.IsSub() {
			return
		}
		adu := repo.NewAppDistributorUnitKafka(d, bou)
		a.toChanOut(adu)
		return
	}
	if len(adub.B) > 0 && len(header) > 0 {
		aduh := repo.NewAppDistributorHeaderKafkaFromSC(d, sc)
		adu := repo.NewAppDistributorUnit(aduh, repo.AppDistributorBody{})
		if d.E() != repo.Last || !repo.IsLastBoundaryEnding(header, bou) {
			adu = repo.NewAppDistributorUnitKafkaPrepared(d, header, nil, adub)
		} else {
			a.S.Reset(repo.NewAppStoreKeyGeneralFromDataPiece(d))
		}
		a.toChanOut(adu)
		return
	}

	if len(adub.B) > 0 && len(header) == 0 {
		if scErr == nil || (scErr != nil && strings.Contains(scErr.Error(), "no header found")) {
			prepErrs = append(prepErrs, scErr)

			if len(sc.From[repo.NewAppStoreKeyDetailed(d)]) == 2 {

				old := sc.From[repo.NewAppStoreKeyDetailed(d)][true].D.H
				adub.B = append(old, adub.B...)
			}

			aduh := repo.NewAppDistributorHeaderKafkaFromSC(d, sc)
			adu := repo.NewAppDistributorUnit(aduh, repo.AppDistributorBody{})

			if d.E() == repo.Last || d.E() == repo.False {
				a.S.Reset(repo.NewAppStoreKeyGeneralFromDataPiece(d))
				if d.E() != repo.Last || !repo.IsLastBoundaryEnding(d.GetBody(0), bou) {
					adu.B = adub
				}
			} else {
				adu.B = adub
			}

			a.toChanOut(adu)

		}
	}
}

// CalcBody creates body of unit to be transfered. Tested in application_test.go
func CalcBody(d repo.DataPiece, bou repo.Boundary) (repo.AppDistributorBody, []byte, error) {
	var err error
	b := d.GetBody(0)
	adub, header := repo.AppDistributorBody{}, make([]byte, 0)
	if d.IsSub() {
		return adub, b, nil
	}
	header, err = d.H(bou)

	if err != nil {
		if !strings.Contains(err.Error(), "is ending part") &&
			!strings.Contains(err.Error(), "no header found") {
			return repo.AppDistributorBody{}, d.GetBody(repo.MaxHeaderLimit), err
		}
	}

	adub = repo.AppDistributorBody{B: b}

	if len(header) > 0 && len(header) < len(b) {
		adub = repo.AppDistributorBody{B: d.GetBody(0)[len(header):]}
		if err != nil && strings.Contains(err.Error(), "is ending part") {
			return adub, header, err
		}
		return adub, header, nil
	}

	return adub, header, nil
}
func (a *App) ChanInClose() {
	close(a.A.C.ChanIn)
}
func (a *App) SetStopping() {
	a.A.stopping = true
}
func (a *App) Stopping() bool {
	return a.A.stopping
}

func (a *App) Stop() {}
