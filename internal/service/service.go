package service

import (
	"github.com/vynovikov/study/highLoadParser/internal/logger"
	"github.com/vynovikov/study/highLoadParser/internal/repository"
)

type ParcerService interface {
	Serve(ParserServiceDTO)
}

type parcerServiceStruct struct {
	repo repository.ParserRepository
}

func NewParserService(r repository.ParserRepository) *parcerServiceStruct {
	return &parcerServiceStruct{
		repo: r,
	}
}

func (s *parcerServiceStruct) Serve(sDTO ParserServiceDTO) {
	logger.L.Infoln("in service.Serve got some data")

	for _, v := range sDTO.U {

		s.repo.Register(v)

	}
}
