package generators

func Filter[T any](done <-chan struct{}, inputStream <-chan T, operator func(T) bool) <-chan T {
	filteredStream := make(chan T)

	go func() {
		defer close(filteredStream)

		for {
			select {
			case <-done:
				return
			case i, ok := <-inputStream:
				if !ok {
					return
				}

				if !operator(i) {
					break
				}

				select {
				case <-done:
					return
				case filteredStream <- i:
				}
			}
		}
	}()

	return filteredStream
}
