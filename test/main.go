package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/valkey-io/valkey-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mainContext, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"localhost:6379"}})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	key := "testkey"
	value := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	cacheTTL := time.Second * 5

	valueAsJsonStr, err := json.Marshal(&value)
	if err != nil {
		log.Fatalln(fmt.Errorf("json marshal error: %w", err))
	}

	if err := client.Do(mainContext, client.B().Set().Key(key).Value(valkey.BinaryString(valueAsJsonStr)).Nx().Ex(cacheTTL).Build()).Error(); err != nil {
		log.Printf(fmt.Errorf("[1] client error: %w", err).Error())
	}

	d, err := client.Do(mainContext, client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		log.Printf(fmt.Errorf("[2] client error: %w", err).Error())
	}

	fmt.Printf("Result before expire: %s\n", d)

	// wait for cache to expire
	time.Sleep(cacheTTL)

	e, err := client.Do(mainContext, client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		log.Printf(fmt.Errorf("[3] client error: %w", err).Error())
	}

	fmt.Printf("Result after expire: %s\n", e)

	// trying to reset expire time
	if err := client.Do(mainContext, client.B().Set().Key(key).Value(valkey.BinaryString(valueAsJsonStr)).Nx().Ex(cacheTTL).Build()).Error(); err != nil {
		log.Printf(fmt.Errorf("[4] client error: %w", err).Error())
	}

	// wait for cache half of cacheTTL
	time.Sleep(cacheTTL / 2)

	// trying to reset expire time
	if err := client.Do(mainContext, client.B().Set().Key(key).Value(valkey.BinaryString(valueAsJsonStr)).Nx().Ex(cacheTTL).Build()).Error(); err != nil {
		log.Printf(fmt.Errorf("[5] client error: %w", err).Error())
	}

	a, err := client.Do(mainContext, client.B().Ttl().Key(key).Build()).ToInt64()
	if err != nil {
		log.Printf(fmt.Errorf("[6] client error: %w", err).Error())
	}

	fmt.Printf("Current TTL: %d\n", a)
}
