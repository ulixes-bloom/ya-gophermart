package workerpool

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

type jobCh[T any] chan T
type jobResCh[T any] chan T
type jobErrCh[T any] chan error
type jobHandler[T any] func(context.Context, T) (T, error)

type pool[T any] struct {
	numOfWorkers int            // количество worker'ов в pool'е
	jobCh        jobCh[T]       // канал с job'ами
	jobHandler   jobHandler[T]  // обработчик job'ы
	wg           sync.WaitGroup // синхронизатор работы worker'ов в pool'e
	results      jobResCh[T]    // канал с результатами работы всех job
	errors       jobErrCh[T]    // канал с ошибками работы job
	once         sync.Once      // гарантирует, что каналы закроются только один раз
}

func New[T any](ctx context.Context, numOfWorkers, jobChSize int, jobHandler jobHandler[T]) *pool[T] {
	jobCh := make(chan T, jobChSize)
	results := make(chan T, jobChSize)
	errors := make(chan error, jobChSize)

	p := pool[T]{
		jobCh:        jobCh,
		numOfWorkers: numOfWorkers,
		jobHandler:   jobHandler,
		results:      results,
		errors:       errors,
	}

	// Запуск worker'ов
	for range p.numOfWorkers {
		go p.startWorker(ctx)
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
	// Выполняется строго один раз
	p.once.Do(func() {
		close(p.jobCh)
		p.wg.Wait()
		close(p.results)
		close(p.errors)
	})
}

// Запуск worker'a
func (p *pool[T]) startWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Завершение работы при отмене контекста
			log.Debug().Msg("Worker pool stopping due to context cancellation")
			return
		case job, ok := <-p.jobCh:
			if !ok {
				// Если канал job'ов закрыт, завершаем worker'a
				return
			}

			res, err := p.jobHandler(ctx, job)
			if err != nil {
				p.errors <- err
			} else {
				p.results <- res
			}

			p.wg.Done()
		}
	}
}

// Получение результатов
func (p *pool[T]) Results() <-chan T {
	return p.results
}

// Получение ошибок
func (p *pool[T]) Errors() <-chan error {
	return p.errors
}
