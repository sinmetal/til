package main

import (
	"fmt"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
)

func main() {
	mem := memcache.New("10.14.80.3:11211")
	mem.MaxIdleConns = 1000

	for {
		var wg sync.WaitGroup
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				id := uuid.New().String()
				err := mem.Add(&memcache.Item{
					Key:   id,
					Value: []byte(id),
				})
				if err != nil {
					fmt.Println(err.Error())
				}
			}(i)
		}
	}
}
