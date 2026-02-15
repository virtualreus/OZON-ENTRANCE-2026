package in_memory

import (
	"ozon_entrance/internal/domain/entities"
	"sync"
)

type InMemory struct {
	Mutex sync.RWMutex

	ByShort    map[string]entities.Link
	ByOriginal map[string]string
}

func NewInMemory() *InMemory {
	return &InMemory{
		ByShort:    make(map[string]entities.Link),
		ByOriginal: make(map[string]string),
	}
}
