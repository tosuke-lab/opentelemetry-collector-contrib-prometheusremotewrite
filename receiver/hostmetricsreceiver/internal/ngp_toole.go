package internal

import "time"

func CurrentTime() (uint64, error) {
	currentTime := float64(time.Now().UnixNano()) / float64(time.Second)
	t := currentTime
	return uint64(t), nil
}
