package generators

func Map[T, S any](done <-chan struct{}, inputStream <-chan T, operator func(T) S) <-chan S {
	mappedStream := make(chan S, 1)

	go func() {
		defer close(mappedStream)

		for {
			select {
			case <-done:
				return
			case i, ok := <-inputStream:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case mappedStream <- operator(i):
				}
			}
		}
	}()

	return mappedStream
}
