package pubsub

import (
	"context"
	"rankland/load"
)

func Publish(ctx context.Context, key string, value interface{}) error {
	err := load.GetRedis().Publish(ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}
