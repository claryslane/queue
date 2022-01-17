package main

import (
	"context"
	"time"

	"github.com/Claryslane/queue/internal/channel"
)

func SendAll(userCh *chan channel.Message, chanName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer close(*userCh)

	for message := range *userCh {
		persons := channels[chanName].Out
		for _, person := range persons {
			go func() {
				*person <- message
			}()
		}
	}

	select {
	case <-ctx.Done():
		cancel()
		return
	}
}
