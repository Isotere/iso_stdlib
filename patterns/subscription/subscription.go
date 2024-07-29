package subscription

import (
	"context"
	"log"
	"time"
)

type Subscription[T any] interface {
	Updates() <-chan T
}

type Fetcher[T any] interface {
	Fetch() (T, error)
}

func NewSubscription[T any](ctx context.Context, fetcher Fetcher[T], freq int) Subscription[T] {
	s := &subscription[T]{
		fetcher: fetcher,
		updates: make(chan T),
	}

	// Запуск задачи, предназначенной для получения наших данных
	go s.serve(ctx, freq)

	return s
}

type subscription[T any] struct {
	fetcher Fetcher[T]
	updates chan T
}

func (s *subscription[T]) Updates() <-chan T {
	return s.updates
}

func (s *subscription[T]) serve(ctx context.Context, checkFrequency int) {
	clock := time.NewTicker(time.Second / time.Duration(checkFrequency))

	type fetchResult struct {
		fetched T
		err     error
	}

	fetchDone := make(chan fetchResult, 1)

	for {
		select {
		// Таймер, который запускает фетчер
		case <-clock.C:
			go func() {
				fetched, err := s.fetcher.Fetch()
				fetchDone <- fetchResult{fetched, err}
			}()
		// Случай, когда результат фетчера готов к использованию
		case result := <-fetchDone:
			fetched := result.fetched
			if result.err != nil {
				log.Printf("Fetch error: %v \n Waiting the next iteration...", result.err.Error())
				break
			}
			s.updates <- fetched
		// Случай, когда нам нужно закрыть сервер
		case <-ctx.Done():
			return
		}
	}
}
