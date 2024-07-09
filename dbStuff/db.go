package db

import "sync"

type db struct {
	path string
	mu   *sync.Mutex
}
