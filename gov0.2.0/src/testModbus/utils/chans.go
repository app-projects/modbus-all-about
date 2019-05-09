package utils

func IsChanClosed(ch chan interface{}) bool {
	select {
	case <-ch:
		return false
	}
	return true
}

func CloseChan(ch chan interface{}) {
	if !IsChanClosed(ch) {
		close(ch)
	}
}
