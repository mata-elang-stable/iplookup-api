package cache

import (
	"context"
	"github.com/fadhilyori/iplookup-go/internal/logger"
	"github.com/valkey-io/valkey-go"
	"time"
)

var log = logger.GetLogger()

type ValKeyInstance struct {
	client valkey.Client
	ttl    time.Duration
}

func MustNewValkey(addresses []string, ttl time.Duration) *ValKeyInstance {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: addresses,
	})
	if err != nil {
		log.Fatalf("failed to create valkey client: %v", err)
	}

	if ttl == 0 {
		ttl = time.Hour
	}

	return &ValKeyInstance{
		client: client,
		ttl:    ttl,
	}
}

func (v *ValKeyInstance) Cache(ctx context.Context, key string, value string) error {
	return v.client.Do(ctx, v.client.B().Set().Key(key).Value(value).Nx().Ex(v.ttl).Build()).Error()
}

func (v *ValKeyInstance) Get(ctx context.Context, key string) (string, error) {
	return v.client.DoCache(ctx, v.client.B().Get().Key(key).Cache(), time.Second*10).ToString()
}
