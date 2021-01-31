package schedule

import (
	"time"

	"github.com/evrone/go-service-template/business-logic/domain"
)

type Scheduler struct {
	useCase    domain.EntityUseCase
	interval   time.Duration
	stopSignal chan struct{}
}

func NewScheduler(uc domain.EntityUseCase, seconds int) *Scheduler {
	return &Scheduler{
		useCase:    uc,
		interval:   time.Duration(seconds) * time.Second,
		stopSignal: make(chan struct{}, 1),
	}
}

func (s *Scheduler) Start() {
	go func() {
	LOOP:
		for range time.Tick(s.interval) {
			select {
			case <-s.stopSignal:
				break LOOP
			default:
			}

			s.useCase.DoTranslate()
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopSignal)
}
