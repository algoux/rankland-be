package pubsub

import (
	"context"
	"rankland/load"
)

func Subscribe(ctx context.Context, key string, s chan<- string) {
	sub := load.GetRedis().Subscribe(ctx, key)
	for msg := range sub.Channel() {
		if msg.Channel == key {
			s <- msg.Payload
		}
	}
}
