package sync

// Or комбинация done-каналов. Сработает, когда любой из входных каналов закроется
func Or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			// Вызываем рекурсивно OR с остатком каналов.
			// В итоге когда сработает либо выше, либо внутри цепочки - все схлопнется
			// (производные горутины завершатся вместе с родительской)
			case <-Or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}
