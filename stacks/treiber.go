package stacks

import (
	"sync/atomic"
)

// Основное отличие стека Трайбера от однопоточного заключается в том,
// что несколько потоков имеют доступ к данным в стеке одновременно, а значит, могут удалять и добавлять элементы.
//
// Для его реализации используется операция CAS.
// (Сравнение с обменом (англ. compare and set, compare and swap, CAS) — атомарная инструкция,
// сравнивающая значение в памяти с первым аргументом и, в случае успеха, записывающая второй аргумент в память.)
//
// Условия:
// 1. Добавлять новый элемент только убедившись, что на момент окончания операции, указатель на голову стека остался тот же.
//    Другими словами, элемент, выбранный нами в качестве next, на момент окончания операции все еще актуален.
// 2. При удалении элемента, перед его возвратом, нужно быть уверенным, что мы действительно удаляем текущую голову стека
//    и в качестве новой головы предъявляем H.next

type treiberNode[T any] struct {
	value T
	next  atomic.Pointer[treiberNode[T]]
}

type TreiberStack[T any] struct {
	size atomic.Int64
	head atomic.Pointer[treiberNode[T]]
}

func NewTreiberStack[T any]() *TreiberStack[T] {
	return &TreiberStack[T]{}
}

func (s *TreiberStack[T]) Push(value T) {
	// Создадим новый элемент, который хотим добавить в начало стека.
	node := &treiberNode[T]{value: value}

	for {
		// Запомним, куда указывает голова стека
		head := s.head.Load()
		// Указатель на следующее значение для него — head.
		node.next.Store(head)

		// Попробуем переместить H на новый элемент, при помощи CAS. Если это удалось — добавление прошло успешно.
		// Если нет, то кто-то другой изменил стек, пока мы пытались добавить элемент. Придется начинать сначала в бесконечном цикле.
		if s.head.CompareAndSwap(head, node) {
			s.size.Add(1)
			return
		}
	}
}

func (s *TreiberStack[T]) Pop() (item T, ok bool) {
	for {
		// Запомним, на что указывает голова стека
		// Значение, которое хранит в себе head, — то, что необходимо будет вернуть.
		head := s.head.Load()
		if head == nil {
			return
		}

		next := head.next.Load()
		// Попробуем переместить голову стеком CAS-ом. Если удалось — вернем head.value.
		// Если нет, то это означает, что с момента начала операции стек был изменен. Поэтому попробуем проделать операцию заново.
		// (в бесконечном цикле)
		if s.head.CompareAndSwap(head, next) {
			s.size.Add(-1)
			return head.value, true
		}
	}
}

func (s *TreiberStack[T]) Size() int64 {
	return s.size.Load()
}

func (s *TreiberStack[T]) Empty() bool {
	return s.head.Load() == nil
}
