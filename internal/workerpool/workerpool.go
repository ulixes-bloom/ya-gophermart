package workerpool

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type jobCh[T any] chan T
type jobResCh[T any] chan T
type jobHandler[T any] func(T) (T, error)

type pool[T any] struct {
	numOfWorkers int            // количество worker'ов в pool'е
	jobCh        jobCh[T]       // канал с job'ами
	jobHandler   jobHandler[T]  // обработчик job'ы
	wg           sync.WaitGroup // синхронизатор работы worker'ов в pool'e
	Results      jobResCh[T]    // канал с результатами работы всех job
}

func New[T any](numOfWorkers, jobChSize int, jobHandler jobHandler[T]) *pool[T] {
	jobCh := make(chan T, jobChSize)
	results := make(chan T, jobChSize)
	p := pool[T]{
		jobCh:        jobCh,
		numOfWorkers: numOfWorkers,
		jobHandler:   jobHandler,
		Results:      results,
	}

	for range p.numOfWorkers {
		go p.startWorker()
	}

	return &p
}

// Добавление новой job'ы в worker pool
func (p *pool[T]) Submit(job T) {
	p.wg.Add(1)
	p.jobCh <- job
}

// Остановка и ожидание окончания работы worker pool'а
func (p *pool[T]) StopAndWait() {
	close(p.jobCh)
	p.wg.Wait()
	close(p.Results)
}

// Запуск worker'a
func (p *pool[T]) startWorker() {
	for j := range p.jobCh {
		res, err := p.jobHandler(j)

		if err != nil {
			log.Warn().Msg(err.Error())
		} else {
			p.Results <- res
		}

		p.wg.Done()
	}
}
