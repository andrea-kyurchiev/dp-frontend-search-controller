package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
)

type Cacher interface {
	Close()
	Get(key string) (interface{}, bool)
	Set(key string, data interface{})
	AddUpdateFunc(key string, updateFunc func() (interface{}, error))
	StartUpdates(ctx context.Context, channel chan error)
}

type Cache struct {
	data           sync.Map
	updateInterval *time.Duration
	close          chan struct{}
	updateFuncs    map[string]func() (interface{}, error)
}

func NewCache(updateInterval *time.Duration) (Cacher, error) {
	if updateInterval != nil {
		if *updateInterval <= 0 {
			return nil, errors.New("cache update interval duration is less than or equal to 0")
		}
	}

	return &Cache{
		data:           sync.Map{},
		updateInterval: updateInterval,
		close:          make(chan struct{}),
		updateFuncs:    make(map[string]func() (interface{}, error)),
	}, nil
}

func (dc *Cache) Get(key string) (interface{}, bool) {
	return dc.data.Load(key)
}

func (dc *Cache) Set(key string, data interface{}) {
	dc.data.Store(key, data)
}

func (dc *Cache) Close() {
	if dc.updateInterval != nil {
		dc.close <- struct{}{}
		for key, _ := range dc.updateFuncs {
			dc.data.Store(key, "")
		}
		dc.updateFuncs = make(map[string]func() (interface{}, error))
	}
}

func (dc *Cache) AddUpdateFunc(key string, updateFunc func() (interface{}, error)) {
	dc.updateFuncs[key] = updateFunc
}

func (dc *Cache) UpdateContent(ctx context.Context) error {
	for key, updateFunc := range dc.updateFuncs {
		updatedContent, err := updateFunc()
		if err != nil {
			return fmt.Errorf("failed to update search cache for %s. error: %v", key, err)
		}
		dc.Set(key, updatedContent)
	}
	return nil
}

func (dc *Cache) StartUpdates(ctx context.Context, errorChannel chan error) {
	if len(dc.updateFuncs) == 0 {
		return
	}

	err := dc.UpdateContent(ctx)
	if err != nil {
		errorChannel <- err
		dc.Close()
		return
	}

	if dc.updateInterval != nil {
		ticker := time.NewTicker(*dc.updateInterval)

		for {
			select {
			case <-ticker.C:
				err := dc.UpdateContent(ctx)
				if err != nil {
					log.Error(ctx, err.Error(), err)
				}

			case <-dc.close:
				return
			case <-ctx.Done():
				return
			}
		}
	}
}